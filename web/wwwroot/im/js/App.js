
var App = function( aCanvas) {
	var app = this;

	var
			canvas,
			context,
			webSocket,
			webSocketService,
			messageQuota = 5,
			sid
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


		//console.log('Socket opened!', e);

		 app.authorize( GlobalToken, GlobalSid )


	};

	app.onSocketClose = function(e) {

        alert("ws 已经关闭 ")


		webSocketService.connectionClosed();
	};

	app.onSocketMessage = function(e) {

		console.log( e.data )
		//try {
			data_arr = e.data.split('||')
            _type = data_arr[0]
            _cmd = data_arr[1]
            _sid = data_arr[2]
            _reqid = data_arr[3]
            _data = data_arr[4]

			if( _type=="3"){
				var obj = {
					type:"message",
					from_sid:_sid,
					msg:_data
				}

			}else{
				var obj = JSON.parse(_data);
			}

            console.log( obj  )
			webSocketService.processMessage(obj);
		//} catch(e) {}
	};

	app.sendMessage = function( msg ) {


	    webSocketService.sendMessage( msg  );

	}
    app.pushMessage = function( from_sid,to_sid,msg ) {

        var sendObj = {
            sid: to_sid,
            msg: msg,
        };
        webSocketService.pushMessage( from_sid, sendObj  );

    }

	app.authorize = function(token,sid) {
		webSocketService.authorize(token,sid);
	}

	app.init = function(aCanvas ,sid) {
		canvas = aCanvas;
		context = canvas.getContext('2d');

		webSocket 				= new WebSocket( 'ws://'+document.domain+':9898/ws' );
		webSocket.onopen 		= app.onSocketOpen;
		webSocket.onclose		= app.onSocketClose;
		webSocket.onmessage 	= app.onSocketMessage;

		webSocketService		= new WebSocketService( webSocket );

	}

}
