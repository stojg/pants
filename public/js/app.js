var App = {

	loadPixi: function (url, callback) {
		this.loadScript(url, 'lib-script', callback);
	},

	loadScript: function (url, id, cb) {
		var script = document.getElementById(id) || document.createElement('script');
		if (script.parent) {
			script.remove();
		}
		script.setAttribute('src', url);
		if (cb) {
			function loadHandler() {
				script.removeEventListener('load', loadHandler);
				cb();
			}

			script.addEventListener('load', loadHandler);
		}
		document.body.appendChild(script);
	}
};

