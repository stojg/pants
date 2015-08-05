var BSON = new bson().BSON;
var Long = new bson().Long;

var renderer = PIXI.autoDetectRenderer(window.innerWidth, window.innerHeight, {backgroundColor: 0x1099bb});
document.body.appendChild(renderer.view);

// create the root of the scene graph
var stage = new PIXI.Container();

var interactive = true;
var container = new PIXI.Container(0x66FF99, interactive);
stage.addChild(container);

var graphics = {
	stage: stage,
	graphGraphics: container,
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

var interpolationDelay = 100;

function update() {
	for (var key in sprites) {
		if (!sprites.hasOwnProperty(key)) {
			continue;
		}
		var sprite = sprites[key];

		var now = window.performance.now() - interpolationDelay;
		var currentTimestamp = timeDiff + sprite.snapshots[0].timestamp;

		// wait until we have enough snapshots
		if (sprite.snapshots.length < 2) {
			continue;
		}

		var nextTimestamp = timeDiff + sprite.snapshots[1].timestamp;

		var coefficient = (now - currentTimestamp) / (nextTimestamp - currentTimestamp);
		var t = linearInterpolation(sprite.snapshots[0], sprite.snapshots[1], coefficient);

		sprite.height = 10;
		sprite.width = 10;
		sprite.rotation = t.orientation;
		sprite.x = t.x;
		sprite.y = t.y;

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
		var snapshot = data[key];
		if (typeof sprites[snapshot.id] === "undefined") {
			var texture = PIXI.Texture.fromImage(snapshot.image, true, PIXI.SCALE_MODES.LINEAR);
			var bunny = new PIXI.Sprite(texture);
			sprites[snapshot.id] = bunny;
			sprites[snapshot.id].id = snapshot.id;
			sprites[snapshot.id].x = snapshot.x;
			sprites[snapshot.id].y = snapshot.y;
			sprites[snapshot.id].anchor.x = 0.5;
			sprites[snapshot.id].anchor.y = 0.5;
			sprites[snapshot.id].height = 10;
			sprites[snapshot.id].width = 10;

			sprites[snapshot.id].interactive = true;
			sprites[snapshot.id]
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
			// @todo: this should probably be in the sprites update()
			container.addChild(bunny);
			sprites[snapshot.id].snapshots = [];
		}

		if(snapshot.dead) {
			container.removeChild(sprites[snapshot.id]);
		}

		snapshot.timestamp = msg.timestamp;

		// biggest size of the queue is 20 history items
		sprites[snapshot.id].snapshots.push(snapshot);
		if (sprites[snapshot.id].snapshots.length > maxSnapshotBuffer) {
			sprites[snapshot.id].snapshots.shift();
		}
	}
}

if (window["WebSocket"]) {
	conn = new WebSocket("ws://localhost:8081/ws");
	conn.binaryType = "blob";
	conn.onclose = function (evt) {
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
		console.log(graphics.stage.getMousePosition());
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

