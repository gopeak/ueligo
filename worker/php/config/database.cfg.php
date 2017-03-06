<?php
  
  /**
   *  数据库设定,支持数据库切分，支持主从数据库，支持表垂直分割
   *  
   */ 
    
  //  主数据库设置
  $_config['database']['default']['master'] = array(
                                         'driver' =>'mysql',
                                         'host' =>'192.168.3.10',
                                         'port'=>'3306',
                                         'user' =>'root',
                                         'passwd' =>'max123',
                                         'dbname' =>'kcvim',
                                         'charset'=>'utf8',
										 'weight'=>10,
                                         );

 $_config['database']['default']['slave1'] = array(
                                         'driver' =>'mysql',
                                         'host' =>'192.168.3.10',
                                         'port'=>'3306',
                                         'user' =>'root',
                                         'passwd' =>'max123',
                                         'dbname' =>'kcvim',
                                         'charset'=>'utf8',
										 'weight'=>0,
                                         );
 // 定义哪些PHP模型类使用哪些数据库
$_config['config_map_class']['default']   = array();   

// 是否启用分表策略
$_config['enable_parttion']  = false;

// 分表数量定义，使用hash算法
$_config['table_partition']   = array('user'=>50,'build'=>100,'task'=>100);

 

return $_config;