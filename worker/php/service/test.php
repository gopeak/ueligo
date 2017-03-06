<?php
namespace php\service ;

	use php\model as model;
	/**
	 * 业务逻辑相关的接口
	 *
	 * @author sevenb@mcross.cn
	 *
	 */
	class test  extends  BaseService  {
		 
	 

		public function __construct(){

			  
			$this->uid = 1;
			 
			parent::__construct();
				
		}
		 
	  
	  
	 
  
        /**
         *
         * 登录接口，1成功，2用户不存在，3密码错误
         * 发送 { "cmd":"test.login","params":{"user":"admin_xbd","password":"258369"}}
         * 监听返回 {"cmd":"socket.login","data":{"code":1,"token":"1041"},"ret":200,"time":1431656856,"trace":{},"req_id":0}
         * @param  $params  示例{"user":"admin_xbd","password":"258369"}
         */
        public function login( $params='{"user":"admin_xbd","password":"258369"}' )
        {

            $final['code'] = 0;  
            $final['token'] ='' ; 
            $name =  addslashes( $_REQUEST['params']['user'] ); 
            $pwd  =  addslashes( $_REQUEST['params']['password'] );
              
             
             //    
             
            return $final;
                
        }
        
		 

		 
   }





