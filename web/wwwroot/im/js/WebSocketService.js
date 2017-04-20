var WebSocketService = function( webSocket) {
	var webSocketService = this;
	
	var webSocket = webSocket;
	//var model = model;

	var TypeReq  = "1";
	var TypePush = "3";
	
	this.hasConnection = false;
	
	this.welcomeHandler = function(data) {
        webSocketService.hasConnection = true;
        console.log("welcomeHandler:");
        console.log(data);

    };

    this.failedHandler = function(data) {
        webSocketService.hasConnection = true;
        console.log("failedHandler:");
        console.log(data);
        alert("加入房间失败!")

    };
	
	this.updateHandler = function(data) {
		var newtp = false;
		//console.log( "updateHandler:" );
		 console.log( data );
 
	}
	
	this.messageHandler = function(data) {
		console.log( "messageHandler:" );
		console.log( data );
		
	}
	
	this.closedHandler = function(data) {

	}
	
	this.redirectHandler = function(data) {
		if (data.url) {
			if (authWindow) {
				authWindow.document.location = data.url;
			} else {
				document.location = data.url;
			}			
		}
	}
	
	this.noneHandler = function(data) {
		 
	}
	
	this.processMessage = function(data) {
		//console.log("processMessage:");
		console.log(data);
		var fn = webSocketService[data.type + 'Handler'];
		if (fn) {
			fn(data);
		}
	}
	
	this.connectionClosed = function() {
		webSocketService.hasConnection = false;
		 
	};
	
	this.sendUpdate = function(tadpole) {
		
		//console.log("sendUpdate:")
		//console.log(tadpole);
		var sendObj = {
			type: 'update',
			x: tadpole.x.toFixed(1),
			y: tadpole.y.toFixed(1),
			id:tadpole.id,
			angle: tadpole.angle.toFixed(3),
			momentum: tadpole.momentum.toFixed(3)
		};

		if(tadpole.name) {
			sendObj['name'] = tadpole.name;
		}
        str = this.wrapReqMessage( 'Update',tadpole.id,0,sendObj)
		webSocket.send(str);
	}

	this.wrapReqMessage = function( cmd,sid,reqid,msg ){
		str = msg
		if( typeof(msg)=="undefined" ){
			return false
		}
		if( typeof(msg)=="null" ){
			return false
		}
		if( typeof(msg)=="object" ){
			str =  JSON.stringify(msg)
		}

		return  TypeReq+"||"+cmd+"||"+sid+"||"+reqid+"||"+str

	}
	this.sendMessage = function( msg  ) {
		console.log("sendMessage:"+msg);
		var regexp = /name: ?(.+)/i;
		if(regexp.test(msg)) {
			model.userTadpole.name = msg.match(regexp)[1];
			$.cookie('todpole_name', model.userTadpole.name, {expires:14});
			return;
		}
		
		var sendObj = {
			type: 'message',
			message: msg,
			id:model.userTadpole.id
		};
        str = this.wrapReqMessage( 'Message',model.userTadpole.id,0,sendObj)
		webSocket.send(str);
	}
    this.joinChannel = function( channel_id  ) {
        console.log("joinChannel:"+channel_id);

        str = this.wrapReqMessage( 'JoinChannel',model.userTadpole.id,0,channel_id)
        webSocket.send(str);
    }
	
	this.authorize = function(token,verifier) {
		var sendObj = {
			type: 'authorize',
			token: token,
			verifier: verifier
		};
        str = this.wrapReqMessage( 'Authorize',"",0,sendObj)
		webSocket.send(str);
	}
}