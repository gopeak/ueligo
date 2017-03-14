<?php


/**
 * APC 抽象缓存层
 *
 */
class Apc implements Icache
{
    public $config;
    public  $link;
    public $use;
    public  $connected = false;

    /**
     * 将数据从缓存中取出
     * @param string $key  key
     * @return bool
     */
    public function __construct( $config, $use = false )
    {
        $this->config = $config;
        $this->use = $use;
    }
    
    public function connect()
    {
        	
        if(!$this->use) return false;
        if (!extension_loaded("apc")){
            throw  new Exception('Apc extension is not loaded!', 500);
        }
        $this->connected = true;

        return true;
    }
    
    public function get($key)
    {
        if(!$this->connect()) return false;
        $flag = apc_fetch($key);
        	
        return  $flag;
    }


    /**
     * 将数据从缓存中取出
     * @param string $key  key
     * @return bool
     */
    public function mget( $keys )
    {
    	if(!$this->connect()) return false;
        $flag = apc_fetch($keys);
        return  $flag;
    }

    
    /**
     * 将数据存入缓存中
     * @param string $key   key
     * @param mix    $value 要存入的数据
     * @param int    $life  存活时间
     * @return bool
     */
    public function add($key,$value,$life=0)
    {
    	if(!$this->connect()) return false;
    	if(!is_string($key)){
    		file_put_contents(TMP_PATH.'/cache_err.log', var_export($key,true),FILE_APPEND);
    	}
    	$flag = apc_add($key,$value,$life);
    	return $flag;
    	 
    }


    /**
     * 将数据存入缓存中
     * @param string $key   key
     * @param mix    $value 要存入的数据
     * @param int    $life  存活时间
     * @return bool
     */
    public function set($key,$value,$life=0)
    {
        if(!$this->connect()) return false;
        if(!is_string($key)){
            file_put_contents(TMP_PATH.'/cache_err.log', var_export($key,true),FILE_APPEND);
        }
        $flag = apc_store($key,$value,$life);
        return $flag;
        	
    }
    
    /**
     * 将数据存入缓存中
     * @param string $keys    包含键值和数值
     * @param int    $life  存活时间
     * @return bool
     */
    public function mset( $keys, $life=0 )
    {
    	if(!$this->connect()) return false;
    	$flag = false;
    	foreach ($keys as $key=>$v)
    	{
    		$flag = apc_store($key,$v,$life);
    	}
    	
    	return $flag;
    }
    
    
    /**
     * 将数据存入缓存中
     * @param string $key   key
     * @param mix    $value 要存入的数据
     * @param int    $life  存活时间
     * @return bool
     */
    public function append($key,$value,$life=0)
    {
        if(!$this->connect()) return false;
        $flag = false;
        if( apc_exists($key))
        {
            $flag = apc_store($key,$value,$life);
        }
        return $flag;
        	
    }
    /**
     * 将数据更新到缓存中，如果存在缓存
     * @param string $key   key
     * @param mix    $value 要存入的数据
     * @param int    $life  存活时间
     * @global $this->link 为Memache类的对象
     * @return bool
     */
    public function replace($key,$value,$life=0)
    {
        if(!$this->connect()) return false;
        $flag = apc_store($key,$value,$life);
        return $flag;
    }

    /**
     * 将数据更新到缓存中，如果存在缓存
     * @param string $key   key
     * @param mix    $value 要存入的数据
     * @param int    $life  存活时间
     * @global $this->link 为Memache类的对象
     * @return bool
     */
    public function delete($key)
    {
        if(!$this->connect()) return false;
        $flag = apc_delete($key);
        return $flag;
    }

    /**
     * 清除公共缓存
     * @param $userid
     */
    public function clearCache($keys)
    {
        if(!$this->connect()) return false;
        if(is_array($keys))
        {
	        foreach ( $keys as $key)
	        {
	            $flag = apc_delete($key);
	        }
        }
        return true;
    }
    
    /**
     * 清除所有缓存
     */
    public function flush()
    {
        if(!$this->connect()) return false;
	    $flag = apc_clear_cache('user');
        return $flag;
    }
    
}
