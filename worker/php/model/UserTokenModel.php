<?php

namespace php\model;

/**
 *  
 * 用户token模块
 *
 * @author seven@haowan11.com
 *
 */
class UserTokenModel extends BaseCacheDbModel
{
	/**
	 * 数据库表名
	 */
	public  $table = 'dist_token';
	
	const   DATA_KEY = 'dist_token/'; 
	
	const   VALID_TOKEN_RET_OK = 1;
	
	const   VALID_TOKEN_RET_NOT_EXIST = 2;
	
	const   VALID_TOKEN_RET_EXPIRE = 3;
	
	public   $uid = '';
	
	
	public function __construct( $uid ='',$PERSISTENT=false )
	{
		parent::__construct( $uid,$PERSISTENT );
	
		$this->uid = $uid;
			
	}
	
	/**
	 * 生成token
	 * @param string $uid
	 * @param string $pass
	 * @param string $tokens
	 */
	static function makeUserToken( $uid, $pass ) {
	
		$token_cfg   = parent::getConfigVar( 'data' );
		$public_key  = $token_cfg['token']['public_key'];
		$secret_key  = $token_cfg['token']['secret_key'];
		$expire_time = $token_cfg['token']['expire_time'];
	
		return md5(md5( $uid ).$public_key.$secret_key.$pass .time()).md5(time().$expire_time.$uid);
	
	}
	
	/**
	 * 刷新token
	 * @param string $uid
	 * @param string $pass
	 * @param string $tokens
	 */
	static function makeUserRefreshToken( $uid, $pass ) {
	
		$token_cfg   = parent::getConfigVar( 'data' );
		$public_key  = $token_cfg['token']['public_key'];
		$secret_key  = $token_cfg['token']['secret_key'];
		$expire_time = $token_cfg['token']['expire_time'];
	
		return md5( $uid.$public_key.md5($pass).$secret_key.time() ).md5($uid.$expire_time);
	
	}
	
	/**
	 * 校验token是否有效
	 * @param string $uid
	 * @param string $token
	 * @return array string
	 */
	public function validToken( $uid, $token ) {
		
		$row = $this->getUserToken( $uid );
		
		if( !isset($row['tk_token']) ||  $row['tk_token']!=$token ){
			
			return array( self::VALID_TOKEN_RET_NOT_EXIST,'token值错误!');
		}
		
		$data_config    =  parent::getConfigVar( 'data' );
		
		if( ( time()- intval($row['tk_token_time']))> intval($data_config['token']['expire'])  ){
			return array( self::VALID_TOKEN_RET_EXPIRE,'token值过期了!');
		}
		
		return array( self::VALID_TOKEN_RET_OK,'ok');
		
	}
	
	/**
	 * 生成和刷新token
	 * @param array $user
	 * @return array
	 */
	public function makeToken( $user  ){
		
		$token_row = $this->getUserToken( $user['uid'] );
		
		$token =  self::makeUserToken( $user['uid'], $user['password'] );
		$refresh_token =  self::makeUserRefreshToken( $user['uid'], $user['password'] );
		$userTokenInfo['tk_user_id'] =  $user['uid'];
		$userTokenInfo['tk_token'] =  $token;
		$userTokenInfo['tk_token_time'] = time();
		$userTokenInfo['tk_refresh_token'] =  $refresh_token;
		$userTokenInfo['tk_refresh_token_time'] = time();
		
		if( !isset($token_row['token'])  ) {
				
			$ret = $this->insertUserToken( $userTokenInfo );
				
		}else{
				
			$ret = $this->updateUserToken( $user['uid'] ,$userTokenInfo );
		}
		
		return array( $ret ,$token,$refresh_token) ;
		
	}
	
	
	/**
	 *  获取用户token的记录信息 
	 *  @param $
	 * @return array
	 */
	public function getUserToken( $uid )
	{
		//使用缓存机制
		$fileds	=	'* ';
		$where	=	" Where `tk_user_id`='$uid'  limit 1 ";
		$key	=	self::DATA_KEY.$uid;
		$table  =   $this->table ;
		$final	=	parent::getRowByKey( $table, $fileds, $where, $key );
		return $final;
	
	}
 
 	/**
	 *  获取用户token的记录信息 
	 *  @param $
	 * @return array
	 */
	public function getRowByToken( $token )
	{
		//使用缓存机制
		$fileds	=	'* ';
		$where	=	" Where `tk_token`='$token'  limit 1 ";
		$key	=	'';
		$table  =   $this->table ;
		$final	=	parent::getRowByKey( $table, $fileds, $where, $key );
		return $final;
	
	}
    
    /**
     * 
     * 插入一条用户token记录
     * @param bool
     */
    public function insertUserToken( $insertInfo )
    {
    	
        $key = self::DATA_KEY.$insertInfo['uid'];
        
        $re = parent::insertInfoByKey( $this->table, $insertInfo, $key );
        
        return $re;
        
    }
    
    /**
     * 
     * @param string $uid
     * @param string $updateinfo
     * @return boolean
     */
    public function updateUserToken( $uid ,$updateinfo)
    {
    	if(empty($updateinfo))
    	{
    		return false;
    	}
    	if(!is_array($updateinfo))
    	{
    		return false;
    	}
    	$key  = self::DATA_KEY.$uid;
    	$where= "  where `tk_user_id`='$uid'";
    	$flag =$this->updateInfoByKey($this->table,$where,$updateinfo,$key);
    	 
    	return 	$flag;
    }
    
    
	/**    
	 * 删除用户token记录
	 * Enter description here ...
	 */
    public function delUserToken( $uid )
    {
    	
    	$key   = "";// self::DATA_KEY.$id;
    	$table = $this->table;
        $where = " Where tk_user_id = '$uid'";

        $flag =  parent::deleteBykey( $table, $where, $key );
        return $flag;
    }
    
 
    
}?>