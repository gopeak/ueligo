var WebSocketService = function( webSocket) {
	var webSocketService = this;
	
	var webSocket = webSocket;
	var TypeReq  = "1";
	var TypePush = "3";
	
	this.hasConnection = false;
	
	this.welcomeHandler = function(data) {
        webSocketService.hasConnection = true;
        console.log("welcomeHandler:",data);

        webSocketService.subscripeGroup()
    };


    this.failedHandler = function(data) {
        webSocketService.hasConnection = true;
        console.log("failedHandler:");
        console.log(data);
        console.log("加入失败!")

    };

    this.subscripeGroupfailedHandler = function(data) {
        console.log(data);
        console.log("加入群组失败!")

    };

    this.subscripeGroupHandler = function(data) {
        console.log(data);
        console.log("加入群组成功!")

    };

    this.failedHandler = function(data) {
        webSocketService.hasConnection = true;
        console.log("failedHandler:");
        console.log(data);
        alert("加入失败!")

    };
	
	this.updateHandler = function(data) {
		var newtp = false;
		//console.log( "updateHandler:" );
		 console.log( data );
 
	}
	
	this.messageHandler = function(data) {
		console.log( "messageHandler:" );
        from_info = data.msg.from_info
        var from_sid = data.from_sid

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
            ,type: data.msg_catlog
			,mine:false
            ,content: data.msg.content
        }


        layui.use('layim', function(layim){
            layim.getMessage(obj);
        });
		
	}

    this.messageGroupHandler = function(data) {
        console.log( "groupMessageHandler:" );

        from_info = data.msg.from_info

        group_id = ""
        for(var i=0; i<GlobalGroups.length; i++)
        {
            if  (GlobalGroups[i].channel_id==data.group_channel_id){
                group_id = GlobalGroups[i].id;
            }
        }

        var obj = {
            username:from_info.username
            ,avatar: from_info.avatar
            ,id: group_id
            ,fromid:from_info.id
			,mine:false
            ,type: "group"
            ,content: data.msg.content
        }
        console.log( "messageGroupHandler obj:" );
        console.log( obj );

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


	this.authorize = function(token,sid) {
		var sendObj = {
			type: 'authorize',
			token: token,
			sid: sid
		};
        str = this.wrapReqMessage( 'Authorize',sid,0,sendObj)
		webSocket.send(str);
	}

    this.subscripeGroup = function( ) {
        var sendObj = {
            type: 'SubscripeGroup',
            token: GlobalToken,
            sid: GlobalSid
        };
        str = webSocketService.wrapReqMessage( 'SubscripeGroup',GlobalSid,0,sendObj)
        webSocket.send(str);
    }

}