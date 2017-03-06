<?php
namespace php\service {

	use php\model as model;
	/**
	 * 业务逻辑相关的接口
	 *
	 * @author sevenb@mcross.cn
	 *
	 */
	class socket  extends  BaseService  {
		 
	 
        public $config ;
        
		public function __construct(){

			  
			$this->uid = 1;
            $sdk =   \Application::getInstance(); 
			$this->config =  $sdk ->getServerConfig();
			parent::__construct();
				
		}
		 
	  
	public function enabled( $params='{}' )
    {
        $sdk =   \Application::getInstance(); 
        $final = $sdk->enabled(); 
         //      
        //var_dump( $final );      
        return $final;                    
                
    }    
    
    public function createChannel( $params=array() )
    {
       
    	$params = $params['params'];
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $ret = $sdk->createChannel( $params['name'] );   
        
        return $ret;
            
    }  
    
    public function joinChannel( $params='{"sid":"xxx", "name":"channel_11"}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $ret = $sdk->joinChannel( $params['sid'], $params['name'] ); 
         //    
    	return $ret;
            
            
    } 
    
    public function getChannels( $params='{}' )
    {
        
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $ret = $sdk->getChannels(  ); 
        //var_dump( $ret ); 
        return $ret;
            
    } 
    
    public function getUserJoinChannels( $params='{"sid":"xxx"}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $ret = $sdk->getUserJoinChannels( $params['sid']  ); 
           
    	return $ret;
            
    } 
    
    public function leaveChannel( $params='{"sid":"xxx", "name":"channel_11"}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $final = $sdk->leaveChannel( $params['sid'], $params['name'] ); 
         //    
        //var_dump( $final ); 
        return $final;
            
    } 
    
    public function removeChannel( $params='{ "name":"channel_11"}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $final = $sdk->removeChannel( $params['name'] ); 
         //    
        //var_dump( $final ); 
        return $final;
            
    } 
    
    public function broadcast( $params='{}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        
        print_r( $params );
        
        $params = $params['params'];
        
        //print_r( $params );
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $final = $sdk->broadcast( $params['msg'], $params['name'] ); 
         //    
        //var_dump( $final ); 
        return $final;
              
    } 
    
    public function push( $params='{"sid":"xxx", "msg":"", "from_sid":""}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        //print_r( $params );
        $req_id   = isset( $params['req_id'] ) ? $params['req_id']:0;
        $from_sid = isset( $params['from_sid'] ) ? $params['from_sid']:'';
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $final = $sdk->push( $params['sid'], $params['msg'], $from_sid ,$req_id ); 
            
        //var_dump( $final ); 
        return $final;
            
    } 
    
    
    public function pushBySids( $params='{"sids":"xxx", "msg":"msg ...."}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        //print_r( $params );
        $sdk =   \ChannelService::getInstance( $this->config ); 
        $final = $sdk->pushBySids( $params['sids'], $params['msg'] ); 
         //    
        //var_dump( $final ); 
        return $final;
            
    } 
    
    public function kickBySid( $params='{"sid":"xxx"}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        //print_r( $params );
         
        $sdk =   \SessionService::getInstance( $this->config ); 
        $final = $sdk->kickBySid( $params['sid']  ); 
         //    
        //var_dump( $final ); 
        return $final;
            
    } 
    
    public function getAllSessions( $params='{}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        //print_r( $params );
         
        $sdk =   \SessionService::getInstance( $this->config );  
        $final = $sdk->get_all_session(  );  
         //    
        //var_dump( $final ); 
        return $final;        
              
    }   
    
    public function getSession( $params='{"sid":"xxx"}' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        //print_r( $params );
         
        $sdk =   \SessionService::getInstance( $this->config ); 
        $final = $sdk->getBySid( $params['sid'] ); 
         //    
        //var_dump( $final ); 
        return $final;
            
    } 
    
    public function updateSession( $params='{"sid":"xxx","user":{} }' )
    {
        if( is_string( $params ) ) $params = json_decode( $params , true );
        $params = $params['params'];
        //print_r( $params );
         
        $sdk =   \SessionService::getInstance( $this->config ); 
        $final = $sdk->updateUserBySid( $params['sid'] , $params['user'] );  
        
        //var_dump( $final ); 
        return $final;
                 
    }     
    
	       
	
	/** 
	 *  关闭通知,用于客户端显式关闭连接    
	 * @example 发送 { "cmd":"socket.close" }
			监听返回 {"cmd":"socket.close","data":{"code":1,"msg":"closed"},"ret":200,"time":1431658917,"trace":{},"req_id":0}
	 */
	public function close(   )
	{
 
		$final['ret'] = 1; 
		$final['msg'] = 'closed';  
		return $final;
			
	}   
	
    //              
	public function user_login( $params =array()  ){
			
		//print_r( $params );
       // $userTokenModel = new model\UserTokenModel();


      //  if(!$userTokenModel->getRowByToken($params['token'])){
        //    $final['sid'] = 'error_sid';
      //  }else{
           
           

       // }
  

        return $params['client_idf'];

		
	}
	
	 
 	
		 

		 
	}







}