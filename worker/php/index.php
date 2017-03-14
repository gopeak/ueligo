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

$obj = "1";
 
$cmd = "user.getUser"; 
$sid = md5(time()); 
$req_id = time(); 
$conn = NULL; 
$ret = $ittaphp->route( $sid, $cmd, $req_id, $conn, $obj );


print_r( $ret );
 