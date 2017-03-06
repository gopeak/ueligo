<?php

//时区设置
date_default_timezone_set('Asia/Shanghai');
 
error_reporting(E_ALL);
 

define( 'ENABLE_CACHE' , false );  
define( 'CACHE_HANDLER' , 'redis' ); 
 
// Xhprof设置
$GLOBALS['ENABLE_XHPROF'] = false ;   
$GLOBALS['XHPROF_RATE']   = 1;      //触发xhprof的几率
 
//Debug设置
$GLOBALS['debug']['is_debug'] = true;
$GLOBALS['debug']['is_trace'] = true; 


 

?>