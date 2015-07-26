$(document).ready(function () {
	var baseUrl = location.href.split('?')[0],
		params,
		select = document.getElementById('version'),
		themeSelect = document.getElementById('theme'),
		editor;

	params = App.getUrlParams();

	document.title = 'pixi.js - ' + params.title;


		console.log('Loading local pixi');
		App.loadPixi('/js/vendor/pixi.min.js',onPixiLoaded);

	function onPixiLoaded()
	{
		console.log('pixi loaded');
		loadExample('examples/' + params.s + '/' + params.f);
	}

	function loadExample(url)
	{
		// load the example code and executes it
		App.loadScript(url, 'example-script');
	}

});