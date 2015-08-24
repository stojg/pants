/**
 * Created by slindqvist on 24/08/15.
 */
function mouseDown() {
	var msg = {
		"topic": "input",
		"action": "clickDown",
		"id": this.id
	};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseUp() {
	var msg = {
		"topic": "input",
		"action": "clickUp",
		"id": this.id
	};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseOver() {
	var msg = {
		"topic": "input",
		"action": "mouseOver",
		"id": this.id
	};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
}

function mouseOut() {
	var msg = {
		"topic": "input",
		"action": "mouseOut",
		"id": this.id
	};
	var serialisedMsg = BSON.serialize(msg, false, true, false);
	conn.send(serialisedMsg);
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

