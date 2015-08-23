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
	container: container,
	domContainer: renderer.view,
	renderer: renderer
};

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

var serverTimeDiff = 0;
var serverLatency = 0;

// always be a couple of snapshots behind
var maxSnapshotBuffer = 5;

function update() {
	debugContainer.removeChildren();

	for (var key in sprites) {
		if (!sprites.hasOwnProperty(key)) {
			continue;
		}
		var sprite = sprites[key];

		var now = window.performance.now() - serverLatency;

		// wait until we have enough snapshots
		if (sprite.snapshots.length < 2) {
			continue;
		}

		var currentTimestamp = serverTimeDiff + sprite.snapshots[0].timestamp;
		var nextTimestamp = serverTimeDiff + sprite.snapshots[1].timestamp;

		var coefficient = (now - currentTimestamp) / (nextTimestamp - currentTimestamp);
		var position = linearInterpolation(sprite.snapshots[0], sprite.snapshots[1], coefficient);

		//console.log(sprite.type);
		if (sprite.type === 'graphics') {
			sprite.clear();
			sprite.lineStyle(1, 0xffffff, 1);
			sprite.beginFill(0xffffff, 1);
			sprite.moveTo(position.x, position.y);
			sprite.lineTo(Number(sprite.snapshots[0].data.toX), Number(sprite.snapshots[0].data.toY));
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

if (window["WebSocket"]) {
	conn = new WebSocket("ws://" + document.location.host + "/ws");
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
		//mouse(graphics);
	}
} else {
	console.log('Your browser does not support WebSockets');
}

function handleMessage(msg) {
	var topic = msg.topic;
	var data = msg.data;

	var now = window.performance.now();
	if (topic === 'time_request') {
		serverLatency = (now - msg.client) / 2;
		serverTimeDiff = (now - msg.server) + serverLatency;
		console.log("serverLatency: " + serverLatency + " serverTimeDiff " + serverTimeDiff);
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
			if (data[key].type === "sprite") {
				sprite = createSprite(data[key]);
				container.addChild(sprite);
				sprites[sprite.id] = sprite;
				sprites[sprite.id].snapshots = [];
			} else {
				console.log("unknown type " + data[key].type);
			}

		}

		if (data[key].dead) {
			container.removeChild(sprites[data[key].id]);
			// @todo should be a dead animation?
			delete(sprites[data[key].id]);
			continue;
		}

		data[key].timestamp = msg.timestamp;
		sprites[data[key].id].snapshots.push(data[key]);

		if (typeof sprites[data[key].id] !== 'undefined') {
			while (sprites[data[key].id].snapshots.length > maxSnapshotBuffer) {
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
		var sprite = new PIXI.Sprite(texture);
		//sprites[snapshot.id] = bunny;
		sprite.id = entity.id;
		sprite.type = "sprite";
		sprite.x = entity.x;
		sprite.y = entity.y;
		sprite.anchor.x = 0.5;
		sprite.anchor.y = 0.5;
		sprite.height = 20;
		sprite.width = 20;

		//sprite.interactive = true;
		//sprite
		//	// set the mouse down and touch start callback...
		//	.on('mousedown', mouseDown)
		//	.on('touchstart', mouseDown)
		//	// set the mouse up and touch end callback...
		//	.on('mouseup', mouseUp)
		//	.on('touchend', mouseUp)
		//	.on('mouseupoutside', mouseUp)
		//	.on('touchendoutside', mouseUp)
		//	// set the mouse over callback...
		//	.on('mouseover', mouseOver)
		//	// set the mouse out callback...
		//	.on('mouseout', mouseOut);
		return sprite;
	}
}



//var mouse = function (graphics) {
//
//	addWheelListener(graphics.domContainer, function (e) {
//		zoom(e.clientX, e.clientY, e.deltaY < 0);
//	});
//
//	function zoom(x, y, isZoomIn) {
//		var direction = isZoomIn ? 1 : -1;
//		var factor = (1 + direction * 0.01);
//		graphics.container.scale.x *= factor;
//		graphics.container.scale.y *= factor;
//		//console.log(graphics.stage.getMousePosition());
//		// Technically code below is not required, but helps to zoom on mouse
//		// cursor, instead center of container coordinates
//		//var beforeTransform = getGraphCoordinates(x, y);
//		graphics.container.updateTransform();
//		//var afterTransform = getGraphCoordinates(x, y);
//		//container.position.x += (afterTransform.x - beforeTransform.x) * container.scale.x;
//		//container.position.y += (afterTransform.y - beforeTransform.y) * container.scale.y;
//		graphics.container.updateTransform();
//	}
//};

function mouseDown() {
	var msg = {"topic": "input", "action": "clickDown", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseUp() {
	var msg = {"topic": "input", "action": "clickUp", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseOver() {
	var msg = {"topic": "input", "action": "mouseOver", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseOut() {
	var msg = {"topic": "input", "action": "mouseOut", "id": this.id};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

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

