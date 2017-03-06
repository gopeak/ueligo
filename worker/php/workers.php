<?php
/*
 *  Least-recently used (LRU) queue device
 *  Demonstrates use of the zmsg class
 *  ps -ef |grep "workers.php start" |awk '{print $2}'|xargs kill -9
 * @author Ian Barber <ian(dot)barber(at)gmail(dot)com>
 */
 

define('CUR_PATH', realpath(dirname(__FILE__))); 
require_once CUR_PATH.'/globals.php';
include_once CUR_PATH.'/libs/sdk.php';
require_once CUR_PATH.'/engine/ittaphp.php';
$cfgArr = json_decode( file_get_contents(realpath( CUR_PATH.'/../../').'/config.json') ,true );
$i=0;

 
$gworker_nbr = "0"; 
$conns = array();
worker_thread2($argv[1]);die;

$cli_argv = $argv;
main( $cli_argv );


 
function main( $argv )
{   
    global $cfgArr;
 
    
	if ( !isset( $argv[1] ) ) {
		exit("请提供参数");
	}
	$act = $argv[1];
    
    $nbr_workers = intval( $cfgArr['worker']['worker_num'] );
   // $nbr_workers = isset($argv[2]) ? intval($argv[2]) : $nbr_workers; 
	if( $act =='stop' ) {
		
		$exec_str = "ps -ef |grep workers.php |awk '{print $2}'|xargs kill -9";
		exec( $exec_str ); 
	}
	if( $act =='start' || $act=='restart' ){
		//v($nbr_workers);
		//$exec_str = "ps -ef |grep workers.php |awk '{print $2}'|xargs kill -9";
		//exec( $exec_str );
     
		for ($worker_nbr = 0; $worker_nbr < $nbr_workers; $worker_nbr++) {
		
			$pid = pcntl_fork();
		 
			if ($pid ==0) { 
				 //子进程得到的$pid为0, 所以这里是子进程执行的逻辑。
				worker_thread2( $worker_nbr );
			}
		}
	}
	 
    sleep(1);
}


function worker_thread2(  $worker_nbr ) { 
	global $gworker_nbr ,$cfgArr  ;
	
	error_reporting( E_ALL );
	
	$ittaphp = new php\engine\ittaphp();
	$retries = 10;
	
	$host = $cfgArr['worker']['host'];
	$port = $cfgArr['worker']['port'];
 	
	$fp = @fsockopen( $host, $port , $errno, $errstr, 5 );
	if (!$fp) {
		echo "fsockopen: $errstr ($errno)<br />\n";
	} else {
		$worker_idf = 'php_'.mt_rand(100000,9999999);
		$worker_ready = sprintf("worker.connect||%s||%s||%s", "", $worker_idf, "");
		
		fwrite($fp, $worker_ready."\n");
		usleep(10000);
		$sid = "";
		while ( true ) {
			if( is_resource($fp) &&  !feof($fp) ) { 
				if ( $retries <= 0 ) {
					echo "worker connect failed!\n";
					break;
				}
				
				$resp = fgets($fp, 4096);
				file_put_contents( "d:".$worker_idf.'.txt' , $resp."\n", FILE_APPEND  );
				//echo $resp."\n";
				if( $resp ) {
					$arr = explode( '||', $resp );
					if( count($arr)<4 ) {
						continue;
					}
					$cmd = $arr[0];
					$client_idf  = $arr[1];
					$worker_idf  = $arr[2];
					$worker_data = $arr[3];
					if( $cmd == "worker.ping" ) {
						
						// 心跳检测
						//echo "worker.ping ".$worker_data."\n";
						$up_time = time();
						
					}
					if ( $cmd == "req.msg" ) {
						$_REQUEST = json_decode( $worker_data ,true);
						if( !is_array($_REQUEST) ){
							$ret['code']  = 500;
							$ret['msg']  = 'json格式错误';
							$ret['cmd'] = '';
							$ret['req_id'] = 0;
						}else{
							$req_id = 0;
							if( isset($_REQUEST['req_id']) ) $req_id = $_REQUEST['req_id'];
							$_REQUEST['client_idf'] = $client_idf;
							$_REQUEST['worker_idf'] = $worker_idf;
							$ret = $ittaphp->route( $_REQUEST );
							$ret['cmd'] = $_REQUEST['cmd'];
							$ret['req_id'] = $req_id;
						}
						
						$json = json_encode( $ret );
						$data =  Sprintf("worker.reply||%s||%s||%s", $client_idf, $worker_idf, $json);
						fwrite( $fp , $data."\n" );
					}
					
					if( $cmd== "worker.connect" ) {
						echo sprintf( "worker %s readey!" , $worker_idf )."\n" ;
					}
				}
				usleep ( 1000 );
			}else{
				
				echo " connect closed!\n";
				$fp = null;
				$fp = @fsockopen( $host, $port , $errno, $errstr, 5 );
				if( $fp ) {
					$worker_idf = 'php_'.mt_rand(100000,9999999);
					$worker_ready = sprintf("worker.connect||%s||%s||%s", "", $worker_idf, "");
					fwrite($fp, $worker_ready."\n");
					$retries = 10;
				}else{
					$retries--;
				}
				sleep( 2 );
			}
		}
		
		fclose($fp);
	}
    
}

 


