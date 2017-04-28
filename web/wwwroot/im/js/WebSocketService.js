var WebSocketService = function( webSocket) {
	var webSocketService = this;
	
	var webSocket = webSocket;
	//var model = model;

	var TypeReq  = "1";
	var TypePush = "3";
	
	this.hasConnection = false;
	
	this.welcomeHandler = function(data) {
        webSocketService.hasConnection = true;
        console.log("welcomeHandler:",data);


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

		var from_sid = data.from_sid

        from_info ={}

        for(var i=0; i<GlobalContacts.length; i++)
        {
            if  (GlobalContacts[i].sid==from_sid){
                from_info = GlobalContacts[i];
            }
        }
        console.log( "from_info:" );
        console.log( from_info );
		obj = {
			username:from_info.username
			,avatar: from_info.avatar
			,id: from_info.id
			,type: "friend"
			,content: data.msg
		}
        layui.use('layim', function(layim){
            layim.getMessage(obj);
        });

		
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

	this.wrapPushMessage = function( sid,msg ){
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

		return  TypePush+"||PushMessage||"+sid+"||0||"+str

	}

	this.wrapPushGroupMessage = function( sid,msg ){
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

		return  TypePush+"||PushGroupMessage||"+sid+"||0||"+str

	}

	this.sendMessage = function( sid, msg  ) {


		var sendObj = {
			type: 'message',
			message: msg,
			id:sid
		};
        str = this.wrapReqMessage( 'Message',sid,0,sendObj)
		console.log("sendMessage:"+str);
		webSocket.send(str);
	}

	this.pushMessage = function( sid, msg  ) {
		console.log("pushMessage:");
        console.log( sid );
        console.log( msg );
		str = this.wrapPushMessage( sid,msg)
		webSocket.send(str);
	}
	this.pushGroupMessage = function( sid, msg  ) {
		console.log("pushGroupMessage:");
		console.log( sid );
		console.log( msg );
		str = this.wrapPushGroupMessage( sid,msg)
		webSocket.send(str);
	}



    this.joinChannel = function( sid,channel_id  ) {
        console.log("joinChannel:"+channel_id);

        str = this.wrapReqMessage( 'JoinChannel',sid ,0,channel_id)
        webSocket.send(str);
    }
	
	this.authorize = function(token,sid) {
		var sendObj = {
			type: 'authorize',
			token: token,
			sid: sid
		};
        str = this.wrapReqMessage( 'Authorize',"",0,sendObj)
		webSocket.send(str);
	}
}