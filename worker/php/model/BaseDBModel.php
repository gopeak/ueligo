<?php

namespace php\model;


/**
 * 模块的基类
 * @author seven 
 */
class BaseDBModel extends BaseModel
{
	 
  
	public     $dbConfig; 
	
	public     $curDbConfig;
 
	/**
     * Slave 数据库PDO的对象
     * @see PDO 
     * @var object 
     */
	public     $slaveDB;
	
	/**
     * Master 数据库PDO的对象
     * @see PDO 
     * @var object 
     */
	public     $masterDB;
	
	/**
	 * 用于实现slavedb单例模式
	 * @var self
	 */
	protected static $_slavedbInstance;
	
	
	/**
	 * 用于实现masterdb单例模式
	 * @var self
	 */
	protected static $_masterdbInstance;
	
 
    /**
     * 获取分表中的表名，获取数据库抽象类对象，获取Memcached抽象对象
     * @param $pid
     */
	function __construct( $pid ='',$persister=false )
	{
		parent::__construct( $pid ); 
		
		$this->dbConfig = BaseModel::getConfigVar( 'database' );
		//v($this->dbConfig);
		$this->setTable();
	    $dbName          = "default";
	    $child_classname = get_class( $this );
        if( isset( $this->dbConfig['config_map_class'] ) ){
            foreach ( $this->dbConfig['config_map_class'] as $key => $row) {
                if(in_array($child_classname,$row)){
                    $dbName = $key;
                    break;
                }
            }
        }
		$this->dbName = $dbName;
		//Slave数据库预连接 
		if( empty($this->slaveDB) ) {
			$this->slaveDB = $this->dbConnect( $dbName,$persister );	 
		}
		//var_dump($this->slaveDB);
		if( empty($this->masterDB) ) {
			$this->masterDB = $this->dbMasterConnect($dbName,$persister);	 //Master数据库预连接
		} 
		 
	} 
 
 
	
	/**
	 * Slave 连接数据库
	 * @param $name
	 * @param $persister 是否持久连接
	 * @todo取消全局变量
	 *
	 */
	public function dbConnect( $name = 'default', $persister=false  )
	{
		
		if( isset(self::$_slavedbInstance[$name]) && !is_null(self::$_slavedbInstance[$name]) ) 
		{
			return self::$_slavedbInstance[$name];
			
		}else{
			
			//计算权重
			$weight = 0;
			$tempdata = array();
			foreach ( $this->dbConfig['database'][$name] as $db_key => $one) {
	
				$weight += $one['weight'];
				for ($i = 0; $i < $one['weight']; $i ++) {
					$tempdata[] = $db_key;
				}
			}
			$use = strval(rand(0, $weight-1));
			$slave_name = $tempdata[$use];
			$_database_cfg = $this->dbConfig['database'][$name][$tempdata[$use]];
			unset($tempdata);
			if(!$_database_cfg)
			{
				$msg =  '[CORE] 数据库配置错误';
				throw new \Exception($msg, 500);
				die();
			}
	
			$_path = LIBS_PATH.'/db/PdoDriver.php';
			if(!file_exists( $_path ))
			{
				throw new \Exception('所在载入的类库:'.$_path.' 不存在');
			}
			require_once( $_path );
		  
			self::$_slavedbInstance[$name] =  \PdoDriver::getInstance('slave',$name, $_database_cfg ,$persister );
			return self::$_slavedbInstance[$name];
		}
	
		
	}
	
	
	/**
	 * Master数据库连接
	 * @param $name
	 * @param $PERSISTENT 是否持久连接
	 * @todo取消全局变量
	 */
	public function dbMasterConnect( $name = 'default', $PERSISTENT=false  )
	{
	
		if( isset(self::$_masterdbInstance[$name]) && is_object(self::$_masterdbInstance[$name]) ) 
		{
			return self::$_masterdbInstance[$name];
			
		}else{
			 
			$_database_cfg = $this->dbConfig['database'][$name]['master'];
			unset($tempdata);
			if(!$_database_cfg)
			{
				$msg =  '[CORE] 数据库配置错误';
				throw new \Exception($msg, 500);
				die();
			}
	
			$_path = LIBS_PATH.'/db/PdoDriver.php';
			if(!file_exists( $_path ))
			{
				throw new \Exception('所在载入的类库:'.$_path.' 不存在');
			}
			require_once( $_path );
			
			$_masterdbInstance[$name] =  \PdoDriver::getInstance('master',$name, $_database_cfg ,$PERSISTENT  );
		 
			return $_masterdbInstance[$name];
		}
	
		 
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
    protected function getMultiy($table, $fileds, $where) 
    { 
        $sql	 = "SELECT $fileds  FROM  $table   ".$where;
        $finally = array();
        
        $rows    = $this->slaveDB->getAll($sql, $primkey);
        
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
    protected function insertByInfo( $table, $insertInfo  ) 
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
    protected function replaceByInfo( $table, $insertInfo  )
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
    	 
    	return (boolean)$insert_flag;
    }
 
 

    /**
     * 更新一条记录的信息 
     * @param $userid int 用户ID
     * @param $updateinfo array  可同时更新多个字段值,如 array('u_name'=>'马柱国','u_headerurl'=>100)
     * @return bool
     */
    protected function updateInfo( $table, $where, $updateinfo )
    {
        $sql  =" Update $table Set ";
        //$sql .= self::arrayToSqlSet($updateinfo);
        $sql .= $this->masterDB->parseSets( $updateinfo );
        //write_file($sql, "updateinfo1.txt");
        $sql .="  ".$where;
        if( $this->wtriteThroughBehind($sql) )
        {    
            return true;
        }else{
            return false;
        }
    }

    /**
     * 用于更新的SQL语句是否后台执行
     * @param $sql
     * @return boolean
     */
    public function wtriteThroughBehind($sql)
    {
        //error_log(date('H:i:s').'_'.substr($sql,0,40)."\n",3,TMP_PATH.'/'.date('Y-m-d').'sql_update.log');
		//if( USE_WTRITE_THROUGH_BEHIND == FALSE   )
        //{
            //var_dump($this->masterDB);
			return $this->masterDB->query($sql);
        //}
 
        
        if( WTRITE_THROUGH_BEHIND_HANDLER=='File' )
        {
            $session_id = session_id();
			return file_put_contents(TMP_PATH.'/session_sqls/'.$session_id.'.sql', $sql, FILE_APPEND);
        }
        
        if( WTRITE_THROUGH_BEHIND_HANDLER=='zeromq' && !empty($this->zeroMQ) )
        {       
	       	/*
			$this->zeroMQ->doBackground("gearman_sql",  serialize(array($this->pid, $this->dbConfig,$sql)) );
	        if ($this->zeroMQ->returnCode() != GEARMAN_SUCCESS)
	        {
	              error_log(date('Y-m-d H:i:s').':'.substr($sql,0,40)."\n",3,TMP_PATH.'/'.date('Y-m-d').'gearman_doBackground_fail.log');
				 return $this->masterDB->query($sql);
	        } 
			*/
	        return true;
        }
        
        return false;

    }
	
	
	 /**
     * 删除
     * @param $userid
     * @param $msg_id
     * @return array
     */
    protected function delete( $table, $where ) 
    {
        $sql = " Delete from $table  $where ";
        //echo $sql;//die;
        if( $this->wtriteThroughBehind($sql) )
        { 
            return true;
        }else{
            return false;
        }
    }
   

 
    
    /**
     * 将数组转化成SQL的SET 列表
     * @param $arr 一个一维数组
	 * @todo应该使用循环，速度更快
     */
    static public function arrayToSqlSet($arr)
    {
        $arr_str = var_export($arr, true);
        $patten  = array( "/'([^']+)'\s\=\>/im", "/array \(/im", "/\,\n^\)/im", "/\n/im" );
        $replace = array( "\\1=","","", "" );
        $arr_str = preg_replace( $patten, $replace, $arr_str);
        return $arr_str;
    }
    
    
    static  public function makeInsertSql( $table ,$rows )
    {
    	$sql = "INSERT INTO  $table ";
    	$sql1="( ";
    	$i=0;
    	$count = count($rows);
    	$row   = current( $rows );
    	foreach( $rows as $field=> $value )
    	{
    		$i++;
    		$end = ",";
    		if( $i==$count )
    		{
    			$sql1.=" `$field`".$end;
    		}
    	}
    	$sql1.=")";
    	
    	$j= 0;
    	$countc = count( $rows );
    	foreach( $rows as $k=>$msg )
    	{
    		$sql2="( ";
    		$i=0;
    		$count = count($msg);
    		foreach( $msg as $field=> $value )
    		{
    			$i++;
    			$end = ",";
    			if( $i==$count )
    			{
    				$sql2.=" '$value'".$end;
    			}
    		}
    		$dot = ",";
    		$j++;
    		if( $countc ==$j ) $dot = "";
    		$sql2.=")".$dot;
    		$sql.=$sql2;
    	}
    		
    	$sql.=";";
    	return $sql;
    	
    }
    

    /**
     * 合并两个数组，类似 array_merge
     */
    static public function my_array_merge($arr1,$arr2)
    {
        if(!empty($arr1) && is_array($arr1) && is_array($arr2))
        {
            $tmp_keys=array_keys($arr1);
            foreach ($arr2 as $up_key => $vv) 
            {
                if(in_array($up_key,$tmp_keys))
                {
                    $arr1[$up_key]=$vv;
                }
            }
            return $arr1;
        }
        return $arr1;
    }


    /**
     * 计算表名
     * @return Model需要连接的表名称
     */
    public function setTable()
    {
       // v($this->dbConfig);
	   
	    if( !$this->dbConfig['enable_parttion'] ) 
        {
        	return;
        }
        if( !empty($this->table) )
        {
            $partition = 50;
            if( isset($this->dbConfig['table_partition'][$this->table]) )
            {
                $partition = $GLOBALS['table_partition'][$this->table];
            }
            $this->table = self::partitionByHash($this->pid,$partition);
        } 
        
    }
    
    /**
     * 返回最后插入的ID	
     */
    public function getLastInsertId(){
    	
    	return $this->masterDB->getLastInsId();
    	
    }

    /**
     * 分表算法
     * @param   $pid  用户ID
     * @param   $partitions 分为几个表
     */
    static  function partitionByHash( $pid, $partitions=50 )
    {
        $h  = sprintf("%u", crc32($pid));
        $h1 = intval(fmod($h, $partitions));
        return $h1;
    }


}
