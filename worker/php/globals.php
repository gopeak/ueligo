<?php

/**
 *  入口文件加载开发框架，并设定项目目录和项目状态
 *  @package    /
 *  @author     seven
 *
 */
 


// 定义项目主目录常量
define('APP_PATH', realpath(dirname(__FILE__)).DIRECTORY_SEPARATOR );

define('LIBS_PATH', APP_PATH.'/libs' );

// 引入主配置文件
include_once  APP_PATH."config/app.cfg.php";

 
include_once(APP_PATH.'/libs/function.php');


spl_autoload_register('itta_autoload');