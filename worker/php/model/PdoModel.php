<?php
namespace php\model;
use PDO;
/***********************************************************************************
 * xphp框架，简单易用的微型PHP框架
 * xphp框架PDO数据库驱动
 * -------------------------------------------------------------------------------
 * CopyRight By Seven & 秋士悲
 * 您可以自由使用该源码，但是在使用过程中，请保留作者信息。
 */
class PdoModel{
	/**
	 * PDO预处理类的保存变量
	 * @var object
	 * @access private
	 */
	private $PDOStatement = null;

	/**
	 * 保存PDO数据连接的唯一实例，避免重复连接
	 * @var object
	 * @access public
	 */
	public $link;

	/**
	 * 当前执行的SQL语句
	 * @var string
	 * @access private
	 */
	private $queryStr = '';

	/**
	 * 数据库配置数组，在PdoDriver构造函数中初始化
	 * @var array
	 * @access private
	 */
	private $dbConfig = array();
	
	/**
	 * PdoDriver类构造函数
	 * @param array $dbConfig 数据库连接配置，在app.cfg.php中配置
	 * @param boolean $PERSISTENT 是否持久连接，BaseModel调用时传入
	 * @throws PDOException 抛出的PDO错误
	 * @access public
	 */
	public function __construct( $dbConfig, $PERSISTENT=false ){
		$this->dbConfig = $dbConfig;
		
		//判断服务器环境是否支持PDO
		if(!class_exists('PDO')){
			throw new PDOException( "当前服务器环境不支持PDO，访问数据库失败。",3000 );
		}

		//判断是传入了正确的数据库配置参数
		if(empty($this->dbConfig['host'])){
			throw new PDOException( "没有定义数据库配置，请在配置文件中配置。",3001 );
		}
		
		$names = (isset($this->dbConfig['charset']) && !empty($this->dbConfig['charset'])) ? $this->dbConfig['charset'] : 'utf8';

		//生成数据库配置
		$this->dbConfig['dsn'] 		= sprintf("%s:host=%s;port=%s;dbname=%s", 
										$this->dbConfig['driver'], 
										$this->dbConfig['host'],
										$this->dbConfig['port'],
										$this->dbConfig['dbname']
										); 
		$this->dbConfig['params'] 	= array(
										PDO::MYSQL_ATTR_INIT_COMMAND=> "SET NAMES {$names}",
										PDO::ATTR_CASE=> PDO::CASE_NATURAL,
										PDO::ATTR_ERRMODE=> PDO::ERRMODE_EXCEPTION,
										PDO::MYSQL_ATTR_USE_BUFFERED_QUERY => true,
										PDO::ATTR_PERSISTENT => $PERSISTENT,
										PDO::ATTR_TIMEOUT=>10
										);
	}

	/**
	 * 数据库预连接，保存到$link变量中
	 * @throws PDOException
	 * @access public
	 */
	private function connect(){
		if (!isset($this->link)){
			try{
				$this->link = new PDO($this->dbConfig['dsn'], $this->dbConfig['user'], $this->dbConfig['passwd'], $this->dbConfig['params']);
			} catch (PDOException $e){
				$message	= date('H:i:s') . '--' . $e->getMessage() . PHP_EOL;
				$file		= XPHP_LOG_PATH . 'pdo' . DIRECTORY_SEPARATOR . date('Y-m-d').'_sql_connect_error.log';
				error_log($message, 3, $file);
				throw new PDOException($e->getMessage(), 3002);
			}

			if(!$this->link){
				throw new PDOException('PDO CONNECT ERROR', 3003);
			}
		}
	}

	/**
	 * 执行更新性的SQL语句
	 * @param string $sql 要执行的SQL指令。
	 * @return integer
	 * @access public
	 * @tutorial 返回受修改或删除 SQL语句影响的行数。如果没有受影响的行，则返回 0。失败返回false
	 */
	public function exec($sql=''){
		if (empty($sql)) {
			throw new PDOException('要执行的SQL语句为空。',3002);
		}

		$this->connect();
		
		if(empty($this->link)){
			throw new PDOException('无法连接数据库。',3005);
		}
	
		$this->queryStr = $sql;

		$result = $this->link->exec($sql);
		return $result;
	}
	
	/**
	 * 执行查询性的SQL查询，准备PDOStatement
	 * @param string $sql 要执行的SQL指令。
	 * @return 返回true或者false
	 * @access private
	 * @todo 是否重写异常的抛出
	 */
	private function query($sql=''){
		if (empty($sql)) {
			throw new PDOException('要执行的SQL语句为空。',3002);
		}
	
		$this->connect();
	
		if(empty($this->link)){
			throw new PDOException('无法连接数据库。',3005);
		}
	
		$this->queryStr = $sql;

		//释放前次的查询结果
		if(!empty($this->PDOStatement)){
			$this->PDOStatement = null;
		}
		
		$this->PDOStatement = $this->link->prepare($sql);
		if(empty($this->PDOStatement)){
			return false;
		}
		
		$result = $this->PDOStatement->execute();
		return $result;
	}

	/**
	* 获得所有的查询数据
	* @param string $sql  要执行的SQL指令
	* @param boolean $primkey 是否以主键为下标。使用主键下标，可以返回以数据库主键的值为下标的二维数组
	* @return array 查询得到的数据集，失败返回false
	* @access public
	*/
	public function getAll($sql, $primkey=false){
		$this->query($sql);
		if($primkey){
			$result = $this->PDOStatement->fetchAll(PDO::FETCH_GROUP|PDO::FETCH_ASSOC);
			$result = array_map('reset', $result);
		}else{
			$result = $this->PDOStatement->fetchAll(PDO::FETCH_ASSOC);
		}
		$this->PDOStatement = null;			
		return $result;
	}

	/**
	 * 获得一条查询数据
	 * @param string $sql  要执行的SQL指令
	 * @return 一条查询数据，失败返回 FALSE。
	 * @access public
	 */
	public function getRow($sql){
		$this->query($sql);
		$result = $this->PDOStatement->fetch(PDO::FETCH_ASSOC,PDO::FETCH_ORI_NEXT);
		$this->PDOStatement = null;
		return $result;
	}
	
	/**
	 * 获得一条查询结果一列的一个值，没有数据则返回false
	 * @param string $sql 要执行的SQL指令
	 * @return 获得一条查询结果一列的一个值，没有数据则返回false
	 * @access public
	 */
	public function getOne($sql){
		$this->query($sql);
		$result = $this->PDOStatement->fetchColumn();
		$this->PDOStatement = null;
		return $result;
	}
	
	/**
	* 获取最近一次查询的sql语句
	* @return string
	* @access public
	*/
	public function getLastSql(){
		return $this->queryStr;
	}

	/**
	* 获取最后插入的ID
	* @return integer 最后插入时的数据ID
	* @access public
	*/
	public function getLastInsId(){
		$this->connect();
		return $this->link->lastInsertId();
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
	* 关闭数据库PDO连接
	* @access public
	*/
	public function close(){
		$this->link = null;
	}
}
