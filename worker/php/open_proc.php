<?php 

define('CUR_PATH', realpath(dirname(__FILE__)));
$cmd = "D:\php7\php7.exe  ".CUR_PATH."/workers.php ".time();

for( $i=1;$i<10; $i++ ) {
	//execInBackground ( $cmd );
	pclose(popen("start /B ". $cmd, "r"));
}

function execInBackground($cmd) {
	if (substr(php_uname(), 0, 7) == "Windows"){
		echo $cmd;
		pclose(popen("start /B ". $cmd, "r"));
	} else {
		exec($cmd . " > /dev/null &");
	}
}





?>