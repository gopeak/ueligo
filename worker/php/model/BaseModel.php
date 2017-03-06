<?php

namespace php\model;


/**
 * 模块的基类
 * @author seven 
 */
class BaseModel 
{
	public  $pid;
	
	 
    
    /**
     * 获取分表中的表名，获取数据库抽象类对象，获取Memcached抽象对象
     * @param $pid
     */
	function __construct( $pid =''  )
	{
		 
	    $this->pid = $pid;
		 
		 
	}
	
	
  static  public function getConfigVar( $file ){
		 
		 //v( $file );
		 $_file = APP_PATH.'config/'.$file .'.cfg.php';
		 if( file_exists(  $_file ) )
		 {
			 include  $_file ;
		 }else{
			
			return  array();
			 
		 }
	 
		 
		 return $_config;
		 
	 }
	 
	 
	 /**
	 * 载入类库
	 * @param 类库的路径，相对ittaphp下的libs目录  $package
	 *
	 */
	protected function loadLib( $package )
	{
		$file = str_replace( '.', '/', $package );
		$lib_file = LIBS_PATH."/{$file}.php";
		if(!file_exists( $lib_file ))
		{
			throw new \Exception('所在载入的类库:'.$lib_file.' 不存在');
		}
		require_once $lib_file;
	}
	
 


}
