<?php
/**
　* Mysql数据库PDO操作类 
　*/
 
class PdoDriver 
{ 
	private $PDOStatement = null;
        
    /**
    * 是否使用永久连接
    * @var bool
    * @access private
    */
    private $pconnect  = false;
 
    /**
    * 错误信息
    * @var string
    * @access private
    */
    private $error = '';
 
    /**
    * 单件模式,保存Pdo类唯一实例,数据库的连接资源
    * @var object
    * @access private
    */
    public $link;
 
    /**
    * 当前SQL语句
    * @var string
    * @access private
    */
    public $queryStr = '';
 
    /**
    * 最后插入记录的ID
    * @var integer
    * @access public
    */
    private $lastInsertId = null;
 
    /**
    * 返回影响记录数
    * @var integer
    * @access private
    */
    private $numRows = 0;
 
    // 事务指令数
    private $transTimes = 0;
 
    private $dbConfig = array();
    /**
     * 构造函数，
     * @param $dbconfig 数据库连接相关信息
     */
    public function __construct( $dbConfig, $PERSISTENT=false )
    {
    	$this->dbConfig = $dbConfig;
        if (!class_exists('PDO')){
        	throw new PDOException( "不支持:PDO",3000 );
        }
        
        //若没有传输任何参数，则使用默认的数据定义
       
        if(empty($this->dbConfig['host']))
        {
        	throw new PDOException( "没有定义数据库配置",3001 );
        }
        	  
        $this->dbConfig['dsn'] 	=  sprintf("%s:host=%s;port=%s;dbname=%s", 
										$this->dbConfig['driver'], 
										$this->dbConfig['host'],
										$this->dbConfig['port'],
										$this->dbConfig['dbname']
										);   
        $this->dbConfig['params'] = array(
						            PDO::MYSQL_ATTR_INIT_COMMAND=> "SET NAMES {$this->dbConfig['charset']}",
						            PDO::ATTR_CASE=> PDO::CASE_NATURAL,
						            PDO::ATTR_ERRMODE=> PDO::ERRMODE_EXCEPTION,
						            PDO::MYSQL_ATTR_USE_BUFFERED_QUERY => true,
						            PDO::ATTR_PERSISTENT => $PERSISTENT,
						            PDO::ATTR_TIMEOUT=>10
            						);

    }
    
    /**
     * 
     * 数据库连接，
     */
    public function connect()
    {
    	 if (!isset($this->link) ) 
         {
            try 
            { 
			    // 默认参数
				$params = array(
					PDO::MYSQL_ATTR_INIT_COMMAND => "SET NAMES 'UTF8'",
					//PDO::ATTR_AUTOCOMMIT=>false,
					PDO::ATTR_TIMEOUT =>3,
					PDO::ATTR_PERSISTENT=>false,
					PDO::MYSQL_ATTR_USE_BUFFERED_QUERY => true,
					
				);
				
				if(!empty($this->dbConfig['params'])){
					$params = $this->dbConfig['params'];
				}
				
				$this->link = new PDO( $this->dbConfig['dsn'], $this->dbConfig['user'], $this->dbConfig['passwd'], $params);
			
            } 
            catch (PDOException $e) 
            {
            	//echo $e->getMessage();
            	error_log(date('H:i:s').$e->getMessage()."\n",3,APP_PATH.'/tmp/'.date('Y-m-d').'sql_connect_error.log');
				throw new PDOException( $e->getMessage(),3002 );
				return false;
            }
             
    	    if(!$this->link) 
    	    {
                throw new PDOException('PDO CONNECT ERROR',3002 );
                return false;
            }
        }
    }
    
    
    /**
    * 释放查询结果
    * @access public
    */
    public function free() 
    {
        $this->PDOStatement = null;
    }
 
    
    /**
    * 获得所有的查询数据
    * @access public
    * @param string $sql  SQL指令
    * @param boolean 是否以主键为下标
    * @return array
    */
    public function getAll($sql=null, $primkey=false) 
    {               
		$this->query($sql,'1');
		if($primkey)
		{
        	$results =  $this->PDOStatement->fetchall( PDO::FETCH_GROUP|PDO::FETCH_ASSOC );		//查询数据并形成数组        
         	$results = array_map('reset', $results);   
         	$this->free();
         	return $results;
		}        
        //返回数据集
        $result = $this->PDOStatement->fetchAll(constant('PDO::FETCH_ASSOC'));
        $this->free();
        return $result;
    }
 
 
     /**
    * 获得所有的查询数据
    * @access public
    * @param string $sql  SQL指令
    * @param boolean 是否以主键为下标
    * @return array
    */
    public function fetchAll($sql=null, $primkey=false) 
    {               
		 return $this->getAll( $sql, $primkey );
    }
	
    /**
    * 获得一条查询结果
    * @access public 
    * @param string $sql  SQL指令
    * @return array
    */
    public function getRow($sql=null) 
    {               
        $this->query($sql,'1');
        // 返回数组集
        $result = $this->PDOStatement->fetch(constant('PDO::FETCH_ASSOC'),constant('PDO::FETCH_ORI_NEXT'));
		if( $result===false ) $result = array();
		$this->free();
        return $result;
    }
	
   /**
    * 获得一条查询结果
    * @access public 
    * @param string $sql  SQL指令
    * @return array
    */
    public function fetchRow($sql=null) 
    {     
        return $this->getRow( $sql );
    }
	
        
    /**
    * 获得一条查询结果一列的一个值
    * @access public
    * @param string $sql  SQL指令
    * @return array
    */    
    public function getOne($sql=null) 
    {               
        $this->query($sql,'1');
        //返回数据集
        $result = $this->PDOStatement->fetchColumn();
        $this->free();
        return $result;
    }
	
	
    /**
    * 获得一条查询结果一列的一个值
    * @access public
    * @param string $sql  SQL指令
    * @return array
    */    
    public function fetchOne($sql=null) 
    {       
        return $this->getOne( $sql );
    }
  
 
 
 	/**
    * 执行查询 
    * @access public
    * @param string $sql sql指令
    * @param string $type = 0:INSERT, UPDATE 以及DELETE  $type = 1:主要针对 SELECT, SHOW 等指令 
    * @todo 是否重写异常的抛出
    * @return mixed
    */
    public function query( $sql='', $type = '0' ) 
    {
    	$this->connect();
    	 
    	if(empty($this->link))
    	{
    		return false;
    	}
    		
    	$this->queryStr = $sql;
    	//释放前次的查询结果
    	if ( !empty($this->PDOStatement) ){
    		$this->free();
    	}
    	switch ($type)
    	{
    		case '0':
    			//echo $sql;
    			$result = $this->link->exec($this->queryStr);
    			// 有错误则抛出异常
    			if ( false === $result)
    			{
    				//echo $this->queryStr."\n<br>";
    				//$errorInfo = $this->link->errorInfo();
    				//throw PDOException($errorInfo[1].$errorInfo[2]."\n<br>".$this->queryStr,$errorInfo[0]);
    				//die;
    				//error_log(date('H:i:s').'-->'.$this->queryStr."\n",3,TMP_PATH.'/'.date('Y-m-d').'pdo_exec_error.log');
    				return false;
    			}
    			else
    			{
    				$this->numRows = $result;
    				if(strripos($this->queryStr,'insert')!==false)
    				{
    					//@todo下面方法是否必须,当前项目不需要$this->lastInsertId
    					 $this->lastInsertId = $this->link->lastInsertId();
    				}
    				return $this->numRows;
    			}
    			break;
    		case '1':
    		 
    			$this->PDOStatement = $this->link->prepare($this->queryStr);
    			//var_dump($this->link);
    			if(empty($this->PDOStatement))
    			{
    				//error_log(date('H:i:s').'-->'.$this->queryStr."\n",3,TMP_PATH.'/'.date('Y-m-d').'pdo_prepare_error.log');
    				return false;
    			}
    			 
    			$bol = $this->PDOStatement->execute();
    			 
    			return $bol;
    			break; 
    	} 
       
    }
	
	public function exec(  $sql='', $type = '0'  )
    { 
    	return $this->query( $sql , $type  ) ;
    }
	
	
    
    /**
    * 获取最近一次查询的sql语句
    * @access function
    * @param
    * @return String 执行的SQL
    */
    public function getLastSql() 
    {
        return $this->queryStr;
    }
 
    /**
    * 获取最后插入的ID
    * @access public
    * @param
    * @return integer 最后插入时的数据ID
    */
    public function getLastInsId()
    {
    	$this->connect();
    	return $this->link->lastInsertId();
    }
	
	    /**
    * 获取最后插入的ID
    * @access public
    * @param
    * @return integer 最后插入时的数据ID
    */
    public function getLastInsertId()
    { 
    	return $this->getLastInsId();
    }
        
    /**
     * 开始一个事务
     * @access public
     */
    public function beginTransaction(){
    	$this->connect();
    	return $this->link->beginTransaction();
    }
    
    /**
     * 回滚一个事务
     * @access public
     */
    public function rollBack(){
    	$this->connect();
    	return $this->link->rollBack();
    }
    
    /**
     * 回滚一个事务
     * @access public
     */
    public function commit(){
    	$this->connect();
    	return $this->link->commit();
    }
    
    /**
    * 关闭数据库
    * @access public
    */
    public  function close() 
    {
        $this->link = null;
    }
 
    /**
    * SQL指令安全过滤
    * @access public
    * @param string $str  SQL指令
    * @return string
    */
    public  function escapeString($str) 
    {
        return addslashes($str);
    }
 
 
   
  
    /**
    * sets分析,在插入；更新数据时调用
    * @access private
    * @param mixed $values
    * @return string
    */
    public function parseSets( $sets ) 
    {

        $setsStr  = '';
        if(is_array($sets))
        {
            foreach ($sets as $key=>$val)
            {
                $key = $this->addSpecialChar($key);
                $val = $this->fieldFormat($val);
                $setsStr .= "$key = ".$val.",";
            }
            $setsStr = substr($setsStr,0,-1);
        }
        else if(is_string($sets)) 
        {
            $setsStr = $sets;
        }
        return $setsStr;
    }
 
    /**
    * 字段格式化
    * @access private
    * @param mixed $value
    * @return mixed
    */
    private function fieldFormat(&$value) 
    {
        if(is_int($value)) 
        {
            $value = intval($value);
        } 
        else if(is_float($value)) 
        {
            $value = floatval($value);
        } 
        elseif(preg_match('/^\(\w*(\+|\-|\*|\/)?\w*\)$/i',$value))
        {
            // 支持在字段的值里面直接使用其它字段
            // 例如 (score+1) (name) 必须包含括号
            $value = $value;
                
        }
        else if(is_string($value)) 
        {
            $value = '\''.$this->escapeString($value).'\'';
        }
        
        if( is_null($value) ) {
        	$value = "''";
        }
        
        
        return $value;
    }
 
    /**
    * 字段和表名添加` 符合
    * 保证指令中使用关键字不出错 针对mysql
    * @access private
    * @param mixed $value
    * @return mixed
    */
    private function addSpecialChar(&$value) 
    {
        if( '*' == $value ||  '`key`' == $value ||  false !== strpos($value,'(') || false !== strpos($value,'.') || false !== strpos($value,'`')) 
        {
            //如果包含* 或者 使用了sql方法 则不作处理
        } 
        elseif(false === strpos($value,'`') ) 
        {
            $value = '`'.trim($value).'`';
        }
        return $value;
    }
 
    /**
    * 去掉空元素
    * @access private
    * @param mixed $value
    * @return mixed
    */
    private function removeEmpty($value)
    {
        return !empty($value);
    }
 
 
 
}
