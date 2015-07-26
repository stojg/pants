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
	},

	getUrlParams: function () {
		var params = window.location.search.substr(1).split('&');

		if (params.length > 1) {
			// convert params to object
			params = params.reduce(function (obj, val) {
				val = val.split('=');

				obj[val[0]] = decodeURIComponent(val[1]);

				return obj;
			}, {});

		}
		else {
			// defaults to the basic example, there might be better way to do this
			// but this will do for now
			params = {s: "basics", f: "basic.js", title: "Basic"};
		}

		return params;
	}


}

