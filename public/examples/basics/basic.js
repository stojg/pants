var BSON = new bson().BSON;
var Long = new bson().Long;

var renderer = PIXI.autoDetectRenderer(window.innerWidth, window.innerHeight, {backgroundColor: 0x1099bb});
document.body.appendChild(renderer.view);

// create the root of the scene graph
var stage = new PIXI.Container();

var interactive = true;
var container = new PIXI.Container(0x66FF99, interactive);
stage.addChild(container);

/*
 * All the bunnies are added to the container with the addChild method
 * when you do this, all the bunnies become children of the container, and when a container moves,
 * so do all its children.
 * This gives you a lot of flexibility and makes it easier to position elements on the screen
 */
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
	var position = {x: 0, y: 0},
		diffX,
		diffY;

	diffX = to.x - from.x;
	if (Math.abs(diffX) < 0.1) {
		position.x = from.x;
	} else {
		position.x = from.x + coef * diffX;
	}
	diffY = to.y - from.y;
	if (Math.abs(diffY) < 0.1) {
		position.y = from.y;
	} else {
		position.y = from.y + coef * diffY;
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
		var curr = timeDiff + sprite.snapshots[0].timestamp;

		// wait until we have enough snapshots
		if (sprite.snapshots.length < 2) {
			continue;
		}

		var next = timeDiff + sprite.snapshots[1].timestamp;

		var coefficient = (now - curr) / (next - curr);
		var t = linearInterpolation(sprite.snapshots[0], sprite.snapshots[1], coefficient);

		var snapshot = sprite.snapshots[0];

		sprite.height = 20;
		sprite.width = 20;
		sprite.rotation = snapshot.rotation;
		sprite.x = t.x;
		sprite.y = t.y;

		// we passed the time for the next timestamp
		if (coefficient > 1) {
			sprite.snapshots.shift();
		}

		//sprite.x += (sprite.speed * (elapsed / 1000));
	}
}

function handleMessage(msg) {
	var topic = msg.topic;
	var data = msg.data;

	// http://gamedev.stackexchange.com/a/93662
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
			sprites[snapshot.id].anchor.x = 0.5;
			sprites[snapshot.id].anchor.y = 0.5;
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
				handleMessage(BSON.deserialize(sprites))
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
	}
} else {
	console.log('Your browser does not support WebSockets');
}

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


