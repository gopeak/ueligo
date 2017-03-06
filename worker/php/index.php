<?php

/**
 *  页面访问入口文件
 *  @package    weiopen
 *  @author     seven
 *
 */
namespace php;



use    php\engine as engine;


require_once realpath( dirname(__FILE__) ).'/globals.php';


$ittaphp = new engine\ittaphp(); 

$obj = json_decode( '{"token":"session_token","cmd":"user.getUser","params":1}');

v($obj);
 
$ret = $ittaphp->route( $obj );


print_r( $ret );
 