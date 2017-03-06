<?php


/**
 *  api访问入口文件
 *  @package    weiopen
 *  @author     seven
 *
 */

namespace php;


use    php\engine as engine;

require_once 'globals.php';


$ittaphp = new engine\ittaphp(); 

$obj = json_decode( '{"token":"session_token","cmd":"user.getUser","params":1}');

 
echo  json_encode( $ittaphp->route( $obj ) );

 