var Assets = (function () {

	var backgrounds = [];
	backgrounds[0] = 0x0089ab;
	backgrounds[1] = 0x1099bb;
	backgrounds[2] = 0x20A9Cb;
	backgrounds[3] = 0x30B9DB;
	backgrounds[4] = 0x40C9EB;
	backgrounds[5] = 0x30DBB9;

	return {
		backgrounds: backgrounds
	}
})();

var Network = (function () {

	/**
	 * the websocket connection
	 */
	var conn;

	var BSON = new bson().BSON;

	var messageHandler;

	var tileSize = 20;

	function onMessage(evt) {
		try {
			var reader = new FileReader();
			reader.onload = function () {
				var msg = new Uint8Array(this.result);
				messageHandler(BSON.deserialize(msg));
			};
			reader.readAsArrayBuffer(evt.data);
		} catch (err) {
			console.log('Failed to deserialize: ', err);
		}
	}

	function onOpen(evt) {
		console.log("connection opened");
		sendMsg({
			"topic": "time_request",
			"client": window.performance.now()
		});
	}

	function onClose(evt) {
		console.log("connection closed");
	}

	function sendMsg(msg) {
		conn.send(BSON.serialize(msg, false, true, false));
	}

	return {
		connect: function (msgHandler) {
			messageHandler = msgHandler
			conn = new WebSocket("ws://" + document.location.host + "/ws");
			conn.binaryType = "blob";
			conn.onclose = onClose;
			conn.onmessage = onMessage;
			conn.onopen = onOpen;
		},
		send: sendMsg
	};

})();

var Basic = (function (assets, network) {

	var serverTimeDiff = 0;

	var serverLatency = 0;

	var backgrounds = assets.backgrounds;

	var sprites = [];

	var renderer = PIXI.autoDetectRenderer(window.innerWidth, window.innerHeight, {backgroundColor: 0x00698B});

	window.onresize = function (event) {
		var w = window.innerWidth;
		var h = window.innerHeight;
		//this part resizes the canvas but keeps ratio the same
		renderer.view.style.width = w + "px";
		renderer.view.style.height = h + "px";
		//this part adjusts the ratio:
		renderer.resize(w, h);
	};

	// create the root of the scene graph
	var stage = new PIXI.Container();


	var backgroundStage = new PIXI.Container();
	backgroundStage.interactive = true;
	backgroundStage.on('mousedown', function (event) {
		var x = Math.floor(event.data.originalEvent.offsetX / tileSize);
		var y = Math.floor(event.data.originalEvent.offsetY / tileSize);
		network.send({
			topic: 'input',
			type: 'click',
			data: [x, y]
		});
		console.log('input');
	});

	stage.addChild(backgroundStage);

	var spriteContainer = new PIXI.ParticleContainer(10000, {
		scale: true,
		position: true,
		rotation: true,
		uvs: false,
		alpha: false,
		x: 0,
		y: 0
	});

	stage.addChild(spriteContainer);

	/**
	 * @Param currentTime - high resolution time
	 */
	function animate() {
		requestAnimationFrame(animate);
		renderer.render(stage);
		update();
	}

	function update() {
		for (var key in sprites) {
			if (!sprites.hasOwnProperty(key)) {
				continue;
			}
			var sprite = sprites[key];

			var latency = serverLatency;

			if (latency < 100) {
				latency = 100;
			}

			var now = window.performance.now() - serverLatency;

			// wait until we have enough snapshots
			if (sprite.snapshots.length < 2) {
				continue;
			}

			var currentTimestamp = serverTimeDiff + sprite.snapshots[0].timestamp;
			var nextTimestamp = serverTimeDiff + sprite.snapshots[1].timestamp;

			var coefficient = (now - currentTimestamp) / (nextTimestamp - currentTimestamp);
			var position = linearInterpolation(sprite.snapshots[0], sprite.snapshots[1], coefficient);

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

	function createGraphic(entity) {
		var graphic = new PIXI.Graphics();
		graphic.type = "graphics";
		graphic.clear();
		graphic.lineStyle(1, 0xffffff, 1);
		graphic.beginFill(0xffffff, 1);
		graphic.moveTo(entity.x, entity.y);
		graphic.lineTo(Number(entity.data.toX), Number(entity.data.toY));
		return graphic;
	}

	function createSprite(entity) {
		var data = entity.data;
		var texture = PIXI.Texture.fromImage(data.Image, true, PIXI.SCALE_MODES.LINEAR);
		var sprite = new PIXI.Sprite(texture);
		sprite.id = entity.id;
		sprite.type = "sprite";
		sprite.x = entity.x;
		sprite.y = entity.y;
		sprite.anchor.x = 0.5;
		sprite.anchor.y = 0.5;
		sprite.height = 20;
		sprite.width = 20;
		return sprite;
	}

	function mapMessage(msg) {
		var tileSize = 20;
		var graphics = new PIXI.Graphics();
		graphics.lineStyle(1, 0x0069AB, 1);
		console.log("Map received");
		var layerContainer = new PIXI.Container();
		for (var key in msg.data) {
			var tile = msg.data[key];
			if (typeof backgrounds[tile.tiletype] !== 'undefined') {
				graphics.beginFill(backgrounds[tile.tiletype], 1);
				graphics.drawRect(tile.x * tileSize, tile.y * tileSize, tileSize, tileSize);
			}
		}
		layerContainer.addChild(graphics);
		backgroundStage.removeChildren();
		backgroundStage.addChild(layerContainer);
	}

	function timeMessage(msg) {
		var now = window.performance.now();
		serverLatency = (now - msg.client) / 2;
		serverTimeDiff = (now - msg.server) + serverLatency;
		console.log("serverLatency: " + serverLatency + " serverTimeDiff " + serverTimeDiff);
	}

	function entityUpdates(spriteContainer, msg) {
		var data = msg.data;
		for (var key in data) {
			if (!data.hasOwnProperty(key)) {
				continue;
			}
			entityUpdate(data[key], msg.timestamp, spriteContainer);
		}
	}

	function entityUpdate(spriteData, msgTimeStamp, container) {
		// create new graphics
		if (spriteData.type === 'graphics') {
			createGraphic(spriteData);
			return;
		}
		// create new sprite
		if (typeof sprites[spriteData.id] === "undefined") {
			var sprite;
			if (spriteData.type === "sprite") {
				sprite = createSprite(spriteData);
				container.addChild(sprite);
				sprites[sprite.id] = sprite;
				sprites[sprite.id].snapshots = [];
			} else {
				console.log("unknown sprite type " + spriteData.type);
				return;
			}
		}

		if (spriteData.dead) {
			container.removeChild(sprites[spriteData.id]);
			// @todo should be a dead animation?
			delete(sprites[spriteData.id]);
			return;
		}

		// copy the message timestamp to the spritedata
		spriteData.timestamp = msgTimeStamp;
		sprites[spriteData.id].snapshots.push(spriteData);

		if (typeof sprites[spriteData.id] !== 'undefined') {
			// keep maxium 5 snapshots
			while (sprites[spriteData.id].snapshots.length > 5) {
				sprites[spriteData.id].snapshots.shift();
			}
		}
	}

	function handleMessage(msg) {
		switch (msg.topic) {
			case 'time_request':
				timeMessage(msg);
				break;
			case 'map':
				mapMessage(msg);
				break;
			case 'update':
				entityUpdates(spriteContainer, msg);
				break;
		}
	}

	return {
		init: function () {

			document.body.appendChild(renderer.view);

			// start animating
			animate(window.performance.now());

			if (window["WebSocket"]) {
				network.connect(handleMessage);
			} else {
				console.log('Your browser does not support WebSockets');
			}
		}
	}
})(Assets, Network);

Basic.init();