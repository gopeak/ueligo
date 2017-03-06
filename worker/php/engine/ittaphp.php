<?php
namespace php\engine {
use php\service as srv;
use php\model as model;

	/**
	 * 开发框架核心文件
	 *
	 */
	class ittaphp
	{
		private $_ctrl;
		private $_action;
		private $_cmd;


		function __construct()
		{
				
			 
		} 

		/**
		 * 开发框架 路由分发，动态调用方法以及构建返回
		 * @throws \Exception
		 *
		 */
		public function route( $args )
		{
		 
			$debugs = array();
			$final = array();
			
			// 处理API接口的请求
		  
			$this->_cmd = $args['cmd'];
			$final['cmd'] = $this->_cmd ;
			try {
					
				if( !strpos( $this->_cmd, '.' ) )
				{
					throw new \Exception('接口调用错误!无效的参数调用!', 500);
				}
				list($service, $method) = explode( '.', $this->_cmd );
				 
				 
				$service_class = 'php\\service\\' . $service; 
				$service_obj = new $service_class;
				
				if( !method_exists($service_obj, $method) )
				{
					throw new \Exception( $method.'方法不存在;', 500);
					die();
				} 
				//f( APP_PATH.'/tmp/request.log' ,date('Y-m-d:H:i:s').': '.var_export( $args,true ),FILE_APPEND );
				

				//开始执行业务逻辑流程
				$result=  $service_obj->$method( $args );
				unset( $args ); 
				//f( APP_PATH.'/tmp/reponse.log' ,date('Y-m-d:H:i:s').': '.var_export( $result,true )."\n\n",FILE_APPEND );
			 
	  
				$final['data'] = $result ;
				$final['code'] = 200;
				$final['time'] = time();

				return $final;
				 
				//捕获游戏异常
			}catch (mmophp\engine\GameException $e){
				$final = array(
						'code'=>$e->getCode(),
						'time'=>time(),
						'data'=>$e->languages,
				);
				return $final;
				//捕获数据库异常
			}catch (PDO\Exception $e){
				  
				$final =   array(
						'code'=>$e->getCode(),
						'time'=>time(),
						'data'=>array('key'=>'database_error' ,'value'=>$e->getMessage()),
				 		);
				return $final;
			}
			//捕获全局异常
			catch (\Exception $e)
			{ 
			
				$final =  array(
						'code'=>$e->getCode(),
						'time'=>time(),
						'data'=>array('key'=>'Server_error' ,'value'=>$e->getMessage()),
						);
				return $final; 
			} 
			return $final;
				
		}




		/**
		 * 手动关闭数据库连接对象
		 */
		public function closePdoConnect(  ) {

			//关闭主数据库PDO连接对象
			if( isset($GLOBALS['masterdbpools']) && !empty($GLOBALS['masterdbpools']) )
			{
				foreach( $GLOBALS['masterdbpools'] as $k =>$vv ) {
					// 判断如果有close方法，则调用
					$vv->link = NULL;
				}
				$GLOBALS['masterdbpools'] = NULL;
			}
				
			//关闭从数据库PDO连接对象
			if( isset($GLOBALS['slavedbpools']) && !empty($GLOBALS['slavedbpools']) )
			{
				foreach( $GLOBALS['slavedbpools'] as $k =>$vv ) {
					// 判断如果有close方法，则调用
					$vv->link = NULL;
				}
				$GLOBALS['slavedbpools'] = NULL;

			}

		}

		static  public function getConfigVar( $file ){
				
			//v( $file );
			$_file = APP_PATH.'config/'.$file .'.cfg.php';
			if( file_exists(  $_file ) )
			{
				include_once  $_file ;
				
			}else{
				 return array();
			}

				
			return $_config;
				
		}



	}


}