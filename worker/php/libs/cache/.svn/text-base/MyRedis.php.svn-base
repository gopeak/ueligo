<?php
/**
 * 
 * Redis 缓存抽象层
 */
class MyRedis  //implements Icache
{
    public  $config;
    public  $link;
    public  $use;
    public  $connected = false;

	/**
	 *  
	 * 
	 */
	public function __construct( $config, $use = false )
	{
		$this->config = $config;
		$this->use = $use;
	}
	
	public function connect()
	{
		 
		if(!$this->use) return false;
		if ( !is_object($this->link ) )
		{ 
			if (!extension_loaded("redis")){
			    throw  new \Exception('Redis extension is not loaded!', 500);
			}  

			$redis = new Redis();
			foreach ($GLOBALS['redis']['server'] as $info) {
				 $redis->connect($info[0], $info[1]); 
			}
		 
			$redis->setOption(Redis::OPT_SERIALIZER, Redis::SERIALIZER_PHP); 
			$this->link = $redis;
			$this->connected = true;
		}
		return true;
	}
	
	
	public function get($key)
	{ 
		if(!$this->connect()) return false;
		$flag = $this->link->get($key); 
		//var_dump( $flag );
		return  $flag;
	}

	
	/**
	 * 将数据从缓存中取出
	 * @param string $key   array('key0', 'key1', 'key5')
	 * @global $this->link 为Redis类的对象
	 * @return bool
	 */
	public function mget( $keys )
	{
		if(!$this->connect()) return false;
		$flag = $this->link->mGet($keys);
		return  $flag;
	}	
	
	
	/**
	 * 存储多个键值，
	 * @param string $key  array('key0' => 'value0', 'key1' => 'value1')
	 * @global $this->link 为Redis类的对象
	 * @return bool
	 */
	public function mset( $keys ,$life=0 )
	{
		if(!$this->connect()) return false;
		$flag = $this->link->mSet($keys);
		return  $flag;
	}
	
	


	/**
	 * 将数据存入缓存中
	 * @param string $key   key
	 * @param mix    $value 要存入的数据
	 * @param int    $life  存活时间 
	 * @global $this->link 为Redis类的对象
	 * @return bool
	 */
	public function set($key,$value,$life=0)
	{
		if(!$this->connect()) return false;
		$flag = $this->link->set($key,$value,$life); 
	    if(!$flag)
        { 
            $resultMessage = '';
             
            error_log($key.':'.date('Y-m-d H:i:s').$resultMessage."\n",3,TMP_PATH.'/'.date('Y-m-d').'_cache_error.log');
			if(empty($key))
			{
				error_log($key.':'.date('Y-m-d H:i:s').json_encode(debug_backtrace())."\n",3,TMP_PATH.'/'.date('Y-m-d').'_trace_cache_error.log');
			}
        }
		return $flag;
			
	}
	
 
 
	/**
	 * 将数据更新到缓存中，如果存在缓存
	 * @param string $key   key
	 * @param mix    $value 要存入的数据
	 * @param int    $life  存活时间 
	 * @global $this->link 为Redis类的对象
	 * @return bool
	 */
	public function replace($key,$value,$life=0)
	{
		if(!$this->connect()) return false;
		 
		 $flag = $this->link->getSet ($key,$value,$life);  
	 
		
	    if(!$flag)
        {  
            error_log($key.':'.date('Y-m-d H:i:s')."\n",3,TMP_PATH.'/'.date('Y-m-d').'_cache_replace_error.log');
        }		
		return $flag;
	}

	/**
	 * 将数据更新到缓存中，如果存在缓存
	 * @param string $key   key
	 * @param mix    $value 要存入的数据
	 * @param int    $life  存活时间 
	 * @global $this->link 为Redis类的对象
	 * @return bool
	 */
	public function delete($key)
	{
		if(!$this->connect()) return false;
		$flag = $this->link->delete ($key);
	   
	    if(!$flag)
        { 
            $resultMessage = '';  
            error_log($key.':'.date('Y-m-d H:i:s').$resultMessage."\n",3,TMP_PATH.'/'.date('Y-m-d').'_cache_error.log');
        }		
		return $flag;
	} 
	 
	/**
	 * 清除公共缓存
	 * @param $userid
	 */
	public function clearCache($keys) 
	{
		if(!$this->connect()) return false;
	 	 
		$flag = $this->link->delete ($key);
	   
	    if(!$flag)
        { 
            $resultMessage = '';  
            error_log($key.':'.date('Y-m-d H:i:s').$resultMessage."\n",3,TMP_PATH.'/'.date('Y-m-d').'_clearCache_error.log');
        }
	 	 
	}
	
    /**
     * 清除所有缓存，请慎用
     */
    public function flush()
    {
        if(!$this->connect()) return false;
	    return $this->link->flushAll();
    }
	
}
