
var App = function( aCanvas) {
	var app = this;

	var
			canvas,
			context,
			webSocket,
			webSocketService,
			messageQuota = 5
	;

	app.update = function() {

	};

	app.draw = function() {

	};

	app.onSocketOpen = function(e) {
		var sendObj = {
				type: 'auth'
			};
		// 认证请求
        str = webSocketService.wrapReqMessage( 'Auth',"",0,"area-global")
		webSocket.send( str );

		//console.log('Socket opened!', e);

		uri = parseUri(document.location)
		if ( uri.queryKey.oauth_token ) {
			app.authorize(uri.queryKey.oauth_token, uri.queryKey.oauth_verifier)
		}

	};

	app.onSocketClose = function(e) {

		webSocketService.connectionClosed();
	};

	app.onSocketMessage = function(e) {
		try {
			data_arr = e.data.split('||')
            _type = data_arr[0]
            _cmd = data_arr[1]
            _sid = data_arr[2]
            _reqid = data_arr[3]
            _data = data_arr[4]

			var obj = JSON.parse(_data);

			webSocketService.processMessage(obj);
		} catch(e) {}
	};

	app.sendMessage = function( msg ) {

	    webSocketService.sendMessage( msg  );

	}

	app.authorize = function(token,verifier) {
		//webSocketService.authorize(token,verifier);
	}

	app.init = function(aCanvas ) {
		canvas = aCanvas;
		context = canvas.getContext('2d');

		webSocket 				= new WebSocket( 'ws://'+document.domain+':9898/ws' );
		webSocket.onopen 		= app.onSocketOpen;
		webSocket.onclose		= app.onSocketClose;
		webSocket.onmessage 	= app.onSocketMessage;

		webSocketService		= new WebSocketService( webSocket );

	}

}
