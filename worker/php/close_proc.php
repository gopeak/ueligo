<?php 

define('CUR_PATH', realpath(dirname(__FILE__)));
$cmd = "taskkill /F /im php7.exe ";

 
execInBackground( $cmd );
 

function execInBackground($cmd) {
	if (substr(php_uname(), 0, 7) == "Windows"){
		pclose(popen("start /B ". $cmd, "r"));
	} else {
		exec($cmd . " > /dev/null &");
	}
}





?>