<?php
 
require_once realpath( dirname(__FILE__) ).'/globals.php';
require_once LIBS_PATH.'/Process.php';

$g_files = array(); 
$g_init_files_time  = array();
 

function monitor_dir(){
    
    $monitor_service_dir = APP_PATH.'service/'; 
    $monitor_model_dir   = APP_PATH.'model/';  
    //echo  "monitor_dir function is called\n" ;
    //echo time()."\n";
    if( check_files_change( $monitor_service_dir )   )
    {
        stop_process_by_key( 'workers' );
        start_process_by_key( 'workers' ,'restart',10 );
        echo  "worker process restarted \n" ;
    }
    
}
$monitor_service_dir = APP_PATH.'service/'; 
$monitor_model_dir   = APP_PATH.'model/'; 

init_files_change( $monitor_service_dir );
init_files_change( $monitor_model_dir );
 

 // Set up a tick handler
 register_tick_function ( "monitor_dir" );

 // Initialize the function before the declare block 
 // Run a block of code, throw a tick every 2nd statement
 declare( ticks = 2 ) {
  
    while( true ) 
    { 
        for ( $x  =  0 ;  $x  <  10 ; ++ $x ) {
            sleep(1);
        }
        
    }
}
 