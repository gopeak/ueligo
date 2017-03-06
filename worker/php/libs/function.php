<?php

/**
 * 开发过程中常用的函数
 * @author  
 */
 
/**
 +----------------------------------------------------------
 * 系统自动加载类库
 +----------------------------------------------------------
 * @param string $classname 对象类名
 +----------------------------------------------------------
 * @return void
 +----------------------------------------------------------
 */
function itta_autoload( $class )
{ 
	if( strpos( $class, 'php\\' )!==false ){
		$class = str_replace( array(  'php\\', '\\'), array( '' ,'/' ), $class ) ;
		$file = APP_PATH.$class . '.php';
		if (is_file($file))
		{
			require_once $file;
			return;
		} 
	} 
	
	$file = APP_PATH.'/service/'.$class . '.php'; 
	if (file_exists($file))
	{
		require_once $file;
		return;
	} 
	
	$file = APP_PATH.'/model/'.$class . '.php'; 
	if (file_exists($file))
	{
		
		require_once $file;
		return;
	} 
 
}

/**
 * 简化var_dump
 */ 
function v( $v1  ) {
    
    var_dump( $v1  );
}
 
/**
 * 简化版 file_put_contents
 */ 
function f( $filename ,  $data ,  $flags = 0    ) { 
    file_put_contents (  $filename ,  $data ,  $flags    );
}

/**
 * 去除空格函数
 * @param $str
 * @return string
 */ 
function trimStr( $str )
{
    $str = trim($str);
    $ret_str = '';
    for($i=0;$i < strlen($str);$i++)
    { 
        if(substr($str, $i, 1) != " ")
        { 
            $ret_str .= trim(substr($str, $i, 1)); 
        }
        else
        {
            while(substr($str,$i,1) == " ")
            {
                $i++;
            }
            $ret_str.= " ";
            $i--; // ***
        }
    }
    return $ret_str;
}


    /**
     * 获取随机字符串
     * @param $len
     * @param $chars
     * @return string
     */
   function rand_string($len, $chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789')
    {
        $string = '';
        for ($i = 0; $i < $len; $i++)
        {
            $pos = rand(0, strlen($chars)-1);
            $string .= $chars{$pos};
        }
        return $string;
    }
 
 
    
    /**
     * mysql_escape_string简写
     * @param string  $str
     * @return string
     */
    function safe_str( $str )
    {
        return mysql_escape_string( $str );
    }
    
    /**
     * mysql_escape_string简写
     * @param string  $str
     * @return string
     */
    function es( $str )
    {
        return mysql_escape_string( $str );
    }
    
    
    // 检查文件是否有修改
    function init_files_change( $monitor_dir )
    { 
        global $g_init_files_time;
        // 递归遍历目录
        $dir_iterator = new RecursiveDirectoryIterator($monitor_dir);
        $iterator = new RecursiveIteratorIterator($dir_iterator);
        $time_now = time();
        foreach ($iterator as $file)
        {
            // 只监控php文件
            if(pathinfo($file, PATHINFO_EXTENSION) != 'php')
            {
                continue;
            } 
            $g_init_files_time[  $monitor_dir.$file->getFilename()  ] =  $file->getMTime(); 
           
        }
        return false;

    }
    
    // 检查文件是否有修改
    function check_files_change( $monitor_dir )
    { 
        global $g_init_files_time; 
        //print_r( $g_init_files_time );
        // 递归遍历目录
        $dir_iterator = new RecursiveDirectoryIterator($monitor_dir);
        $iterator = new RecursiveIteratorIterator($dir_iterator);
        $time_now = time();
        foreach ($iterator as $file)
        {
            // 只监控php文件
            if(pathinfo($file, PATHINFO_EXTENSION) != 'php')
            {
                continue;
            }
            // 如果最近有修改 
            //v( date( 'Y-m-d:H:i:s',$file->getMTime() ) );
            
            if( isset( $g_init_files_time[ $monitor_dir.$file->getFilename() ] ) ) 
            {
                $last_up_time = $g_init_files_time[ $monitor_dir.$file->getFilename() ];
                if(  ( $file->getMTime() -  $last_up_time ) >2     )
                {  
                    $g_init_files_time[ $monitor_dir.$file->getFilename() ] = $file->getMTime();
                    return true;
                }
            }else{
                return true;
            }
            
        }
        return false;

    }
    
    function stop_process_by_key(  $key  )
    {
        
    	$cmd = "taskkill /F /im php7.exe ";
        pclose(popen("start /B ". $cmd, "r"));
        return;
    	
    	$process_ctrl = new Process(  );
          
        $process_key =  $key.'.php';
    
        	
        	if (substr(php_uname(), 0, 7) == "Windows"){
        	 
        		$cmd= 'D:\php7\php7.exe  '.APP_PATH.'close_proc.php';
        		echo $cmd."\n";
        		pclose(popen("start /B ". $cmd, "r"));
        		
        	} else {
        		$cmd= 'ps -ef |grep "'.escapeshellcmd ( $process_key) .'"  |awk \'{print $2}\' |xargs -i kill -9 {}';
            	system ( $cmd );
        	}
         
     }                        
 
    function start_process_by_key( $key ,$cmd, $num )
    {
       
    	stop_process_by_key( $key );
        
        $cmd = "D:\php7\php7.exe  ".CUR_PATH."/workers.php ".time();
        
        for( $i=1;$i<10; $i++ ) {
        	//execInBackground ( $cmd );
        	pclose(popen("start /B ". $cmd, "r"));
        }
        
          
     }

     
     
     
     
     
     
     
     
     
     
     
     
     
    
