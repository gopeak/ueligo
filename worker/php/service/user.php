<?php

namespace php\service {
	
	
	use php\model as model;
	//require_once APP_PATH.'/service/BaseService.php' ;
	
	/**
	 * 用户相关的接口
	 *
	 * @author sevenb@mcross.cn
	 *
	 */
	class user  extends   BaseService  {


	 
		
		
		public function getUserByToken( $token )
		{
			 
			$userTokenModel = new model\UserTokenModel();
			$tokenRow = $userTokenModel->getRowByToken( $token );
			
			//var_dump($uid);
			$userModel = new  model\UserModel( $tokenRow['uid'] ,true );
			$userModel->uid = $uid; 
			$user = $userModel->getUser();
			if( isset($user['password']) ) unset($user['password']);
				
			return  $user;

		}

		/**
		 * 获取单个用户信息
		 * @param string $uid
		 * @return array
		 */
		public function getUser( $arg )
		{
			$uid = 	$arg->params;		 
			if( empty($uid ) ) { 
				throw new GameException(array('user_id_err','没有uid'));
			}
			//var_dump($uid);
			$userModel = new  model\UserModel( $uid ,true );
			$userModel->uid = $uid;

			$user = $userModel->getUser();
			if( isset($user['password']) ) unset($user['password']);
				
			return  $user;

		}
		
		public function room_list(  $arg  )
		{
			global $cfgArr;
            
            
			$final   = $cfgArr['area'];
			  
			return $final;
		}
		
		public function send(  $arg  )
		{
			global $cfgArr;
            var_dump( $arg );
			$params   = $arg['params'];
			$room_id  = $params['channel'];
			//unset( $arg,$params ); 
			$final['code']  = 1; 
			
			$msg = array();
			$msg['cmd'] = 'room.push'; 
			$msg['req_id'] = 0;
            
            $sdk =   \ChannelService::getInstance( $cfgArr ); 
             
            $data['from'] = $params['from'];
            $data['target'] = $params['target'];
            $data['content'] = $params['content'];
            $ret = $sdk->broadcast( $data, $params['channel'] ); 
             
			//v( $ret  );
			$final['msg']   = 'send msg ok';
			
			return $final;
		}
	 
		public function disconnect(  $arg  )
		{
			$params   = $arg->params;
			$username = $params->user; 	
			$room_id  = $params->channel;
			unset( $arg,$params );
			$final['msg'] = '';  
			$final['code']  = 1; 
			if( !function_exists('apc_fetch') ) {
				$final['msg']   = 'server error ,apc extension not installed!';
				$final['userlist']   = $userlist;
				$final['code']  = 2;
				return $final;
			}
			$key = 'rooms_'.$room_id;
			$cache_data = apc_fetch( $key ); 
			if( $cache_data === false  ) {
				$userlist = [];
			}else{
				$userlist = json_decode( $cache_data, true);	
				if( $key= array_search( $username , $userlist   )!==false ) {
					unset($userlist[$key]);
				}
			} 		
			$ret = apc_store( $key , json_encode( $userlist ) ); 
			 
			$final['msg']   = 'logout ok';
			
			return $final;
		}

		/**
		 * 登录接口，1成功，2用户不存在，3密码错误
		 * 登录成功后
		 * @param   $username  string
		 * @param   $password string
		 * @param   $lat string
		 * @param   $long string
		 * @return  number
		 */
		public function login(  $arg  )
		{
			var_dump( $arg ); 
            
			$params   = $arg['params'];
			$username = $params['user'];
			$password = $params['pwd'];	
			$room_id  = $params['channel'];
		 
			
			$final['msg'] = '';  
			$final['code']  = 1; 
			$userlist = array();
			
			if( !function_exists('apc_fetch') ) {
				$final['msg']   = 'server error ,apc extension not installed!';
				$final['users']   = $userlist;
				$final['code']  = 2;
				return $final;
			}
			$key = 'rooms_'.$room_id;
			$cache_data = apc_fetch( $key );
			
			//apc_delete( $key); 
			if( $cache_data === false  ) {
				$userlist[] = $username;;
			}else{
				$userlist = json_decode( $cache_data, true);	
				if( array_search( $username , $userlist   )===false ) {
					$userlist[] = $username;
				}
			} 		
			$ret = apc_store( $key , json_encode( $userlist ) );  
            var_dump(  $key, $ret );
            $final['sid'] = 'sid_'.strval( mt_rand(1000,9999) ) ;
			$final['users']   = $userlist;
			$final['msg']   = 'login ok';
			
			return $final;
			 
		}
		/**
		 * 角色初始化
		 */
		public function act_init(  $arg  )
		{  
			$params   = $arg->params; 
			$act_id = $params->act_id;
			unset( $arg,$params );
			
			$final = array("id"=>"1000080","area"=>"area-1"); 
			$final['code']  = 1;
			$final['msg']   = 'ok';
						
			return $final;
		}

 

		/**
		 * 注销接口
		 * @return  number
		 */
		public function logout(   )
		{

			//清除会话
			
			$final['code']  = 1;
			$final['msg']  = 'ok';
			return $final;
		}

		  


 
		/**
		 *
		 * @throws GameException
		 * @return number
		 */
		public function throwExceptionTest(  ) {

			throw new GameException(array('key'=>'nologin','content'=>'登录已经失效了'));

		}
 

	}







}