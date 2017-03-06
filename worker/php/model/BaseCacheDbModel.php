<?php

namespace php\model;

/**
 * 缓存操作的基类
 * @author seven@mcross.cn;sky@mcross.cn
 */
class BaseCacheDbModel extends BaseDBModel
{
  
	
	/**
	 * 
	 * PHP Extendsion Memcached 对象
	 * @var object
	 */
	public    $cacheDB;
 
 
    
    /**
     * 缓存的默认过期时间,单位:秒
     * @var int
     */
    public $key_exprie = 0;
    
    
    public $uid = '';
    
    
    
    /**
     *  继承数据库操作基类并预连接Memcached抽象对象
     * @param $pid
     */
	function __construct( $pid ='',$PERSISTENT=false )
	{
		parent::__construct( $pid,$PERSISTENT ); 
		
		$this->cacheDB = $this->memConnect();		 //缓存预连接
		 
		 
	}
	
	 
	/**
	 * 实例化Memcache封装的对象，实际上并未真正的连接到Memcache中
	 * @param serverkey $name
	 * @throws \Exception
	 * @return  object memcache封装的对象
	 */
	protected function memConnect( $name = 'server' )
	{
		static $_mems;
	
		if( !isset($_mems[$name]) )
		{
			$cacheHandle = strtolower(CACHE_HANDLER);
			$cacheConfig = BaseModel::getConfigVar( 'cache' );
			$_cache_cfg = $cacheConfig[$cacheHandle]['server'];
			
			if( !$_cache_cfg && $cacheHandle!='apc' )
			{
				$msg = '[CORE] 未配置 缓存 参数';
				throw new \Exception($msg, 500);
				die();
			}
			 
			if($cacheHandle=='redis')
			{
				$this->loadLib("cache.Icache");
				$this->loadLib("cache.MyRedis");
				$_mems[$name] = new \MyRedis( $_cache_cfg, ENABLE_CACHE);
			}
			if($cacheHandle=='apc')
			{
				$this->loadLib("cache.Icache");
				$this->loadLib("cache.Apc");
				$_mems[$name] = new \Apc( array(), ENABLE_CACHE);
			}
		 
		}
		return $_mems[$name];
	}
	
 
	
	/**
	 * 从一个表里查询出数据
	 * @param $table
	 * @param $fileds
	 * @param $where
	 * @param $primkey
	 * @param $key
	 * @return array
	 */
	protected function getRowByKey($table, $fileds, $where, $key)
	{
	 
		//使用缓存机制
		if(!empty($key))
		{
			//$this->cacheDB->delete($key);
			$memflag=$this->cacheDB->get($key);
			if($memflag!==false)
			{
				return $memflag;
			}
		}
		$sql	 = "SELECT $fileds  FROM  $table   ".$where;
		//  echo $sql;
		//error_log(date('H:i:s').'_'.$sql."\n",3,TMP_PATH.'/'.date('Y-m-d').'sql_select_row.log');
		$finally = array();
		$rows    = $this->slaveDB->getRow($sql);
		 
		if( $rows!==false  && !empty($key) )
		{
			$replace_flag = $this->cacheDB->set($key, $rows, $this->key_exprie);
		}
		if(!empty($rows))
		{
			$finally = $rows;
		}
		return  $finally;
	}


	/**
	 * 从一个表里查询出数据
	 * @param $table
	 * @param $fileds
	 * @param $where
	 * @param $primkey
	 * @param $key
	 * @todo 将getTableFromcacheDBByKey重新命名为getRowsFromcacheDBByKey
	 * @return array
	 */
	protected function getRowsByKey($table, $fileds, $where, $key ,$primkey_assoc = false )
	{
		//使用缓存机制
		if(!empty($key))
		{
			$memflag = $this->cacheDB->get($key);
			if($memflag!==false)
			{
				return $memflag;
			}
		}
		$sql	 = "SELECT $fileds  FROM  $table   ".$where;
	
		//error_log(date('H:i:s').'_'.$sql."\n",3,TMP_PATH.'/'.date('Y-m-d').'sql_select_rows.log');
		$finally = array();
		$rows    = $this->slaveDB->getAll( $sql,$primkey_assoc );
	
		if(  $rows!==false   && !empty($key) )
		{
			$replace_flag = $this->cacheDB->set($key, $rows, $this->key_exprie);
		}
		if(!empty($rows))
		{
			$finally = $rows;
		}
		return  $finally;
	}

    /**
     * 从一个表里查询取得一个字段的值，并关联缓存
     * @param $table
     * @param $fileds
     * @param $where
     * @param $primkey
     * @param $key
     * @return array
     */
    protected function getOneByKey($table, $fileds, $where, $key) 
    {
 
    	//使用缓存机制
        if(!empty($key))
        {
            $memflag=$this->cacheDB->get($key);
            if($memflag!==false) 
            {
            	return $memflag;
            }
        }
        $sql	 = "SELECT $fileds  FROM  $table   ".$where;
        // echo $sql;
		//error_log(date('H:i:s').'_'.$sql."\n",3,TMP_PATH.'/'.date('Y-m-d').'sql_select_row.log');
        $finally = array();
        $one     = $this->slaveDB->getOne($sql);
       
        if( $one!==false  && !empty($key) )
        {
       		$replace_flag = $this->cacheDB->set($key, $one, $this->key_exprie);
        }
        
        return  $one;
    }

    /**
     * 从一个表里查询出数据
     * @param $table
     * @param $fileds
     * @param $where
     * @param $primkey
     * @param $key
	 * @todo 将getTableFromcacheDBByKey重新命名为getRowsFromcacheDBByKey
     * @return array
     */
    protected function getMultipleRowByKey($table, $fileds, $where, $key ) 
    {
        //使用缓存机制
        if(!empty($key))
        {
            $memflag = $this->cacheDB->get($key);
            if($memflag!==false) 
            {
            	return $memflag;
            }
        }
        $sql	 = "SELECT $fileds  FROM  $table   ".$where;
	   
		//error_log(date('H:i:s').'_'.$sql."\n",3,TMP_PATH.'/'.date('Y-m-d').'sql_select_rows.log');
        $finally = array();
        $rows    = $this->slaveDB->getAll($sql);
        
        if(  $rows!==false   && !empty($key) )
        {
            $replace_flag = $this->cacheDB->set($key, $rows, $this->key_exprie);
        }
        if(!empty($rows)) 
        {
        	$finally = $rows;
        }
        return  $finally;
    }

    /**
     * 从一个表里查询出数据
     * @param $table
     * @param $fileds
     * @param $where
     * @param $primkey
     * @param $key
     * @return array
     */
    protected function getMultipleRowsByKeys($table, $fileds, $where, $keys) 
    {
        //使用缓存机制
        if(!empty($keys))
        {
        	$memflag=$this->cacheDB->mget($keys);
            if($memflag!==false) 
            {
            	return $memflag;
            }
        }
        $sql	 = "SELECT $fileds  FROM  $table   ".$where;
        $finally = array();
        
        $rows    = $this->slaveDB->getAll($sql, $primkey);
        if(!empty($keys))
        {
            $replace_flag = $this->cacheDB->set($key, $rows, $this->key_exprie);
        }
        if(!empty($rows)) 
        {
        	$finally=$rows;
        }
        return  $finally;
    }   

    
    /**
     * 构造插入的SQL语句目的利于缓存,能够同步缓存的数据
     * 该函数用于缓存中数据是一维数组的情况  
     * @param $table
     * @param $insertId
     * @param $insertInfo 必须是完整的数据
     * @return bool
     */
    protected function insertInfoByKey($table, $insertInfo, $key="" ) 
    {
        if(empty($table) or empty($insertInfo)) 
        {
            return false;
        }
        $new_arr     = array();
        $insert_flag = false;
        $sql  = "Insert  into  $table Set  ";
        $sql .= $this->masterDB->parseSets($insertInfo);
        //echo $sql; 
        $insert_flag = $this->wtriteThroughBehind($sql) ;
        if(!empty($key) && $insert_flag!=false)
        {
            $this->cacheDB->set($key, $insertInfo, $this->key_exprie);
        }
        return (boolean)$insert_flag;
    }
    
    /**
     * 构造插入的SQL语句目的利于缓存,能够同步缓存的数据
     * 该函数用于缓存中数据是一维数组的情况
     * @param $table
     * @param $insertId
     * @param $insertInfo 必须是完整的数据
     * @return bool
     */
    protected function replaceInfoByKey($table, $insertInfo, $key="" )
    {
    	if(empty($table) or empty($insertInfo))
    	{
    		return false;
    	}
    	$new_arr     = array();
    	$insert_flag = false;
    	$sql  = "Replace  into  $table Set  ";
    	$sql .= $this->masterDB->parseSets($insertInfo);
    	//echo $sql;
    	$insert_flag = $this->wtriteThroughBehind($sql) ;
    	if(!empty($key) && $insert_flag!=false)
    	{
    		$this->cacheDB->set($key, $insertInfo, $this->key_exprie);
    	}
    	return (boolean)$insert_flag;
    }
    
    
    /**
     * 构造插入的SQL语句目的利于缓存,能够同步缓存的数据
     * 该函数适用于缓存中数据是二维数组的情况  
     * @param $table
     * @param $insertId
     * @param $insertInfo
     * @return bool
     */
    protected function insertMultipleInfoByKey( $table, $insertId, $insertInfo, $key="", $life=0, $json_key=array() ) 
    {
        if(empty($table) or empty($insertInfo)) 
        {
        	return false;
        }
        $new_arr     = array();
        $insert_flag = false;
        $sql  = "INSERT into   $table Set  ";
        $sql .= $this->masterDB->parseSets($insertInfo);
       
        $insert_flag = $this->wtriteThroughBehind($sql) ;
        if(!empty($key) && $insert_flag!=false)
        {
        	$cacheFlag = $this->cacheDB->get($key);
        	if( empty($insertId) )
        	{
        	     $this->cacheDB->delete($key);
        	     return  (boolean)$insert_flag;
        	}
            if($cacheFlag === false) 
            {
                $this->cacheDB->set($key, array($insertId=>$insertInfo), $this->key_exprie);
            }else
            {
                $cacheFlag[$insertId] = $insertInfo;
                if( $this->cacheDB->set($key,$cacheFlag, $this->key_exprie)===false )
                {
                    $this->cacheDB->delete($key);
                }
            }
        }
        return (boolean)$insert_flag;
    }

    /**
     * 执行更新的SQL语句，能够同步缓存的数据，适用于memcache的Key中为2维数组的情况 
     * @param $userid int 用户ID
     * @param $updateinfo array 镖师对应tb_user表的字段/值,可同时更新多个字段值,如 array('u_name'=>'马柱国','u_headerurl'=>100)
     * @return bool
     */
    protected function updateMultipleInfoByKey($table, $where, $updateinfo, $auto_id, $key="", $life=0, $json_key=array())
    {
        $sql  = " Update $table Set ";
        $sql .= $this->masterDB->parseSets($updateinfo);
        $sql .="  ".$where;
        $GLOBALS['__debugs']['sql'] = $sql;
        //echo $sql;die;
        if( $this->wtriteThroughBehind($sql) )
        {
            if(!empty($key))
            {
                $cacheFlag = $this->cacheDB->get($key);
                if($cacheFlag===false) 
                {
                	return true;
                }
                if(isset($cacheFlag[$auto_id]))
                {
                    $tmp=$cacheFlag[$auto_id];
                    if(!empty($tmp) && is_array($tmp))
                    {
                        $cacheFlag[$auto_id]=self::my_array_merge($tmp, $updateinfo);
                        //废弃array_merge函数,$cacheFlag[$auto_id]=array_merge($tmp,$updateinfo);
                        $replace_flag = $this->cacheDB->replace($key, $cacheFlag, $this->key_exprie);
                        if($replace_flag==false)
                        {
                            $this->cacheDB->delete($key);
                        }
                    }
                }else{
                	$this->cacheDB->delete($key);
                }
            }
            return true;
        }
        else
        {
            return false;
        }
    }

    /**
     * 更新一条记录的信息，能够同步缓存的数据，适用于memcache的Key中为1维数组的情况 
     * @param $userid int 用户ID
     * @param $updateinfo array 镖师对应tb_user表的字段/值,可同时更新多个字段值,如 array('u_name'=>'马柱国','u_headerurl'=>100)
     * @return bool
     */
    protected function updateInfoByKey($table, $where, $updateinfo, $key="")
    {
        $sql  =" Update $table Set ";
        //$sql .= self::arrayToSqlSet($updateinfo);
        $sql .= $this->masterDB->parseSets( $updateinfo );
        //write_file($sql, "updateinfo1.txt");
        $sql .="  ".$where;
        //echo $sql;
        if( $this->wtriteThroughBehind($sql) )
        {   
            if(!empty($key))
            {
                $cacheFlag=$this->cacheDB->get($key);

                if($cacheFlag===false)
                {
                    return true;
                }

                if(!empty($cacheFlag))
                {
                    foreach ( $updateinfo as $updateKey =>$info  )
                    {
                        if( isset($cacheFlag[$updateKey]) )
                        {
                            $cacheFlag[$updateKey] = $info;
                        }
                    }
                    $replace_flag = $this->cacheDB->replace($key, $cacheFlag, $this->key_exprie);
                    if($replace_flag == false)
                    {
                        $this->cacheDB->delete($key);
                    }
                }
            }
           
            return true;
        }
        else
        {
            return false;
        }
    }
 
	
	
	 /**
     * 删除
     * @param $table
     * @param $where
     * @return array
     */
    protected function deleteMultipleBykey($table,$where,$auto_id,$key="") 
    {
        $sql = " Delete from $table  $where ";
        //echo $sql;//die;
        if( $this->wtriteThroughBehind($sql) )
        {
            if(!empty($key))
            {
                $cacheFlag=$this->cacheDB->get($key);
                if($cacheFlag===false) 
                {
                	return true;
                }
                if(isset($cacheFlag[$auto_id]))
                {
                    unset($cacheFlag[$auto_id]);  
					$replace_flag = $this->cacheDB->replace($key, $cacheFlag, $this->key_exprie);
					if($replace_flag==false)
					{
						$this->cacheDB->delete($key);
					} 
                }
            }
            return true;
        }else{
            return false;
        }
    }
    

	

    /**
     * 删除
     * @param $userid
     * @param $msg_id
     * @return array
     */
    protected function deleteBykey($table,$where,$key="") 
    {
        $sql = " Delete from $table  $where ";
        //echo $sql;
        if( $this->wtriteThroughBehind($sql) )
        {
            if(!empty($key))
            { 
                
                $this->cacheDB->delete($key);
            }
            return true;
        }
        return false;
    }
    
    
 
 
 
 

}
