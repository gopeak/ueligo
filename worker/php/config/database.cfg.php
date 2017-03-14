<?php
  
  /**
   *  数据库设定,支持数据库切分，支持主从数据库，支持表垂直分割
   *  
   */ 
    
  //  主数据库设置
  $_config['database']['default']['master'] = array(
                                         'driver' =>'mysql',
                                         'host' =>'127.0.0.1',
                                         'port'=>'3306',
                                         'user' =>'root',
                                         'passwd' =>'123456',
                                         'dbname' =>'socket_db',
                                         'charset'=>'utf8',
										 'weight'=>10,
                                         ); 
    
 

return $_config;