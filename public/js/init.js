$(document).ready(function () {
	document.title = 'Pants';
	App.loadScript("/js/vendor/bson.js");
	App.loadScript("/js/addWheelListener.js");
	App.loadPixi('/js/vendor/pixi.min.js', function () {
		App.loadScript('examples/basics/basic.js', 'example-script');
	});
});