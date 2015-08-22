var BSON = new bson().BSON;
var Long = new bson().Long;

var renderer = PIXI.autoDetectRenderer(window.innerWidth, window.innerHeight, {backgroundColor: 0x1099bb});
document.body.appendChild(renderer.view);

// create the root of the scene graph
var stage = new PIXI.Container();

var interactive = true;
var container = new PIXI.Container(0x66FF99, interactive);
stage.addChild(container);

var debugContainer = new PIXI.Container();
stage.addChild(debugContainer);


var graphics = {
	stage: stage,
	graphGraphics: container,
	domContainer: renderer.view,
	renderer: renderer
};

var debugs = [];

container.x = 0;
container.y = 0;

// start animating
animate(window.performance.now());

/**
 * @Param currentTime - high resolution time
 */
function animate() {
	requestAnimationFrame(animate);
	renderer.render(stage);
	update();
}

var sprites = [];

var timeDiff = 0;
var latency = 0;

// always be a couple of snapshots behind
var maxSnapshotBuffer = 5;

var linearInterpolation = function (from, to, coef) {
	var position = {x: 0, y: 0, orientation: 0};

	var diffX = to.x - from.x;
	if (Math.abs(diffX) < 0.1) {
		position.x = from.x;
	} else {
		position.x = from.x + coef * diffX;
	}
	var diffY = to.y - from.y;
	if (Math.abs(diffY) < 0.1) {
		position.y = from.y;
	} else {
		position.y = from.y + coef * diffY;
	}
	var diffOrientation = to.orientation - from.orientation;
	if (Math.abs(diffOrientation) < 0.1) {
		position.orientation = to.orientation;
	} else {
		position.orientation = from.orientation + coef * diffOrientation;
	}
	return position;
};

// @todo(stig): This should probably be measured
var interpolationDelay = 100;

function update() {
	debugContainer.removeChildren();

	for (var key in sprites) {
		if (!sprites.hasOwnProperty(key)) {
			continue;
		}
		var sprite = sprites[key];

		var now = window.performance.now() - latency;

		// wait until we have enough snapshots
		if (sprite.snapshots.length < 2) {
			continue;
		}

		var currentTimestamp = timeDiff + sprite.snapshots[0].timestamp;
		var nextTimestamp = timeDiff + sprite.snapshots[1].timestamp;

		var coefficient = (now - currentTimestamp) / (nextTimestamp - currentTimestamp);
		var position = linearInterpolation(sprite.snapshots[0], sprite.snapshots[1], coefficient);

		//console.log(sprite.type);
		if(sprite.type === 'graphics') {
			//sprite.clear();
			sprite.clear();
			sprite.lineStyle(1, 0xffffff, 1);
			sprite.beginFill(0xffffff, 1);
			sprite.moveTo(position.x, position.y);
			//console.log(sprite.snapshots[0].data.toX);
			sprite.lineTo(Number(sprite.snapshots[0].data.toX), Number(sprite.snapshots[0].data.toY));
			//console.log(Number(sprite.snapshots[0].data.toX*10), position.x);
			//console.log(sprite.snapshots[0].data.toX, sprite.snapshots[0].data.toY);
			//sprite.lineTo(0, 0);
			//sprite.moveTo(t.x, t.y);
			//console.log(t.data);
		} else {
			sprite.x = position.x;
			sprite.y = position.y;
			sprite.rotation = position.orientation;
			sprite.height = 20;
			sprite.width = 20;
		}
		// we passed the time for the next timestamp
		if (coefficient > 1) {
			sprite.snapshots.shift();
		}
	}
}

function handleMessage(msg) {
	var topic = msg.topic;
	var data = msg.data;

	var now = window.performance.now();
	if (topic === 'time_request') {
		latency = (now - msg.client) / 2;
		timeDiff = (now - msg.server) + latency;
		console.log("latency: " + latency + " timeDiff " + timeDiff);
		return;
	}

	for (var key in data) {
		if (!data.hasOwnProperty(key)) {
			continue;
		}

		if (data[key].type === 'graphics') {
			createGraphic(data[key]);
			continue;
		}

		// create new sprite
		if (typeof sprites[data[key].id] === "undefined") {
			var sprite;
			if(data[key].type === "sprite") {
				sprite = createSprite(data[key]);
				container.addChild(sprite);
				sprites[sprite.id] = sprite;
				sprites[sprite.id].snapshots = [];
			} else {
				console.log("unknown type " + data[key].type);
			}

		}

		if(data[key].dead) {
			container.removeChild(sprites[data[key].id]);
			// @todo should be a dead animation?
			delete(sprites[data[key].id]);
			continue;
		}

		data[key].timestamp = msg.timestamp;
		sprites[data[key].id].snapshots.push(data[key]);

		if(typeof sprites[data[key].id] !== 'undefined') {
			while(sprites[data[key].id].snapshots.length > maxSnapshotBuffer) {
				sprites[data[key].id].snapshots.shift();
			}
		}
	}

	function createGraphic(entity) {
		var graphic = new PIXI.Graphics();
		graphic.type = "graphics";
		graphic.clear();
		graphic.lineStyle(1, 0xffffff, 1);
		graphic.beginFill(0xffffff, 1);
		graphic.moveTo(entity.x, entity.y);
		graphic.lineTo(Number(entity.data.toX), Number(entity.data.toY));
		debugContainer.addChild(graphic)
		return graphic;
	}

	function createSprite(entity) {
		var data = entity.data;
		var texture = PIXI.Texture.fromImage(data.Image, true, PIXI.SCALE_MODES.LINEAR);
		var bunny = new PIXI.Sprite(texture);
		//sprites[snapshot.id] = bunny;
		bunny.id = entity.id;
		bunny.type = "sprite";
		bunny.x = entity.x;
		bunny.y = entity.y;
		bunny.anchor.x = 0.5;
		bunny.anchor.y = 0.5;
		bunny.height = 20;
		bunny.width = 20;

		bunny.interactive = true;
		bunny
			// set the mouse down and touch start callback...
			.on('mousedown', mouseDown)
			.on('touchstart', mouseDown)
			// set the mouse up and touch end callback...
			.on('mouseup', mouseUp)
			.on('touchend', mouseUp)
			.on('mouseupoutside', mouseUp)
			.on('touchendoutside', mouseUp)
			// set the mouse over callback...
			.on('mouseover', mouseOver)
			// set the mouse out callback...
			.on('mouseout', mouseOut);
		return bunny;
	}
}

if (window["WebSocket"]) {
	conn = new WebSocket("ws://"+document.location.host+"/ws");
	conn.binaryType = "blob";
	conn.onclose = function (evt) {
		console.log("connection closed");
	};
	conn.onmessage = function (evt) {
		try {
			var reader = new FileReader();
			reader.onload = function () {
				var sprites = new Uint8Array(this.result);
				handleMessage(BSON.deserialize(sprites));
			};
			reader.readAsArrayBuffer(evt.data);
		}
		catch (err) {
			console.log('Failed to deserialize: ', err);
		}
	};
	conn.onopen = function (evt) {
		console.log("connection opened");
		var msg = {"topic": "time_request", "client": window.performance.now()};
		var serialisedMsg = BSON.serialize(msg, false, true, false);
		conn.send(serialisedMsg);
		mouse(graphics);
	}
} else {
	console.log('Your browser does not support WebSockets');
}

var mouse = function (graphics) {

	var graphGraphics = graphics.graphGraphics;

	addWheelListener(graphics.domContainer, function (e) {
		zoom(e.clientX, e.clientY, e.deltaY < 0);
	});

	//addDragNDrop();
	function zoom(x, y, isZoomIn) {
		direction = isZoomIn ? 1 : -1;
		var factor = (1 + direction * 0.01);
		graphGraphics.scale.x *= factor;
		graphGraphics.scale.y *= factor;
		//console.log(graphics.stage.getMousePosition());
		// Technically code below is not required, but helps to zoom on mouse
		// cursor, instead center of graphGraphics coordinates
		//var beforeTransform = getGraphCoordinates(x, y);
		graphGraphics.updateTransform();
		//var afterTransform = getGraphCoordinates(x, y);
		//graphGraphics.position.x += (afterTransform.x - beforeTransform.x) * graphGraphics.scale.x;
		//graphGraphics.position.y += (afterTransform.y - beforeTransform.y) * graphGraphics.scale.y;
		graphGraphics.updateTransform();
	}
};

function mouseDown() {
	this.down = true;
	var msg = {"topic": "input", "action": "clickDown", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseUp() {
	this.down = false;
	if (this.isOver) {
	}
	var msg = {"topic": "input", "action": "clickUp", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseOver() {
	this.isOver = true;
	var msg = {"topic": "input", "action": "mouseOver", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseOut() {
	this.isOver = false;
	var msg = {"topic": "input", "action": "mouseOut", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

