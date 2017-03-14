<?php
    
namespace php;

    use    php\engine as engine;
    
    // 定义项目主目录常量
    define('CUR_PATH', realpath(dirname(__FILE__)).DIRECTORY_SEPARATOR ); 
    require_once CUR_PATH.'globals.php';
    
    
    $srv = new SocketServer('0.0.0.0',8002);
    
    
    
	/*!	@class		SocketServer
		@author		sven
		@abstract 	A Framework for creating a multi-client server using the PHP language.
	 */
	class SocketServer
	{
		/*!	@var		config
			@abstract	Array - an array of configuration information used by the server.
		 */
		protected $config;

		/*!	@var		hooks
			@abstract	Array - a dictionary of hooks and the callbacks attached to them.
		 */
		protected $hooks;

		/*!	@var		master_socket
			@abstract	resource - The master socket used by the server.
		 */
		protected $socket;
		
		/*!	@var		max_clients
			@abstract	unsigned int - The maximum number of clients allowed to connect.
		 */
		public $max_clients = 1024;

		/*!	@var		max_read
			@abstract	unsigned int - The maximum number of bytes to read from a socket at a single time.
		 */
		public $max_read = 4096;

		/*!	@var		clients
			@abstract	Array - an array of connected clients.
		 */
		public $clients;

        /*!	@var		connections
			@abstract	Array - an array of connected clients.
		 */
		public $connections = array();
        
        /*!	@var		buffers 
		 */
		public $buffers = array();
        
        public $builder = NULL;
        
        public $data_obj = NULL;
        
        public $ittaphp = NULL;
        
		/*!	@function	__construct
			@abstract	Creates the socket and starts listening to it.
			@param		string	- IP Address to bind to, NULL for default.
			@param		int	- Port to bind to
			@result		void
		 */
		public function __construct($bind_ip,$port)
		{
			set_time_limit(0);
			$this->hooks = array();

			$this->config["ip"] = $bind_ip;
			$this->config["port"] = $port; 
             
            
            $server_addr  = 'tcp://'.$this->config["ip"].':'.$this->config["port"];
            $this->socket = stream_socket_server ( $server_addr, $errno, $errstr);
            stream_set_blocking($this->socket, 0);
            $base = event_base_new();
            $event = event_new();
            event_set($event, $this->socket, EV_READ | EV_PERSIST, [$this,'ev_accept'], $base);
            event_base_set($event, $base);
            event_add($event);
            event_base_loop($base);
            
            
 
		}
        
                
        function ev_accept($socket, $flag, $base) {
            static $id = 0;
            
            $connection = stream_socket_accept($socket);
            stream_set_blocking($connection, 0);
            
            $id += 1;
            
            $buffer = event_buffer_new($connection, [$this,'ev_read'], [$this,'ev_write'], [$this,'ev_error'], $id); 
            event_buffer_base_set($buffer, $base);
            event_buffer_timeout_set($buffer, 30, 30);
            event_buffer_watermark_set($buffer, EV_READ, 0, 0xffffff);
            event_buffer_priority_set($buffer, 10);
            event_buffer_enable($buffer, EV_READ | EV_PERSIST);
            
            // we need to save both buffee and connection outside
            $this->connections[$id] = $connection;
            $this->buffers[$id] = $buffer;
        }

        function ev_error($buffer, $error, $id) {
            //var_dump( $error );
            event_buffer_disable($this->buffers[$id], EV_READ | EV_WRITE);
            event_buffer_free($this->buffers[$id]);
            fclose($this->connections[$id]);
            unset($this->buffers[$id], $this->connections[$id]);
        }

        
        function ev_read($buffer, $id) {
            $ctx_data = '';
            while ($read = event_buffer_read($buffer, $this->max_read )) {
                
                $ctx_data .=$read ;  
                
            }
            
            $param_arr = array();
            list( $type, $cmd, $sid, $req_id,$data ) = explode( '||',$ctx_data ); 
            $resp_data = $this->dispatch( $sid, $cmd, $req_id, $buffer, $data  );  
            $resp_type = 2;   
            $response_str = "{$resp_type}||{$cmd}||{$sid}||{$req_id}||{$resp_data}\n";
            //var_dump($response_str); 
            event_buffer_write( $buffer , $response_str ,strlen($response_str));
 
        }
        
        
        function dispatch( $sid, $cmd, $req_id, $buffer, $data ){ 
            
            $this->ittaphp = new engine\ittaphp(); 
            return json_encode($this->ittaphp->route( $sid, $cmd, $req_id, $buffer, $data ));
            
        }
 

        function ev_write($buffer, $id)
        {
            // echo "[$id] " . __METHOD__ .PHP_EOL;
        }

		/*!	@function	hook
			@abstract	Adds a function to be called whenever a certain action happens.  Can be extended in your implementation.
			@param		string	- Command
			@param		callback- Function to Call.
			@see		unhook
			@see		trigger_hooks
			@result		void
		 */
		public function hook($command,$function)
		{
			$command = strtoupper($command);
			if(!isset($this->hooks[$command])) { $this->hooks[$command] = array(); }
			$k = array_search($function,$this->hooks[$command]);
			if($k === FALSE)
			{
				$this->hooks[$command][] = $function;
			}
		}

		/*!	@function	unhook
			@abstract	Deletes a function from the call list for a certain action.  Can be extended in your implementation.
			@param		string	- Command
			@param		callback- Function to Delete from Call List
			@see		hook
			@see		trigger_hooks
			@result		void
		 */
		public function unhook($command = NULL,$function)
		{
			$command = strtoupper($command);
			if($command !== NULL)
			{
				$k = array_search($function,$this->hooks[$command]);
				if($k !== FALSE)
				{
					unset($this->hooks[$command][$k]);
				}
			} else {
				$k = array_search($this->user_funcs,$function);
				if($k !== FALSE)
				{
					unset($this->user_funcs[$k]);
				}
			}
		}
 

		/*!	@function	disconnect
			@abstract	Disconnects a client from the server.
			@param		int	- Index of the client to disconnect.
			@param		string	- Message to send to the hooks
			@result		void
		*/
		public function disconnect($id,$message = "")
		{ 
            if( $message ) { 
                fwrite( $this->connections[$id], $message."\n" ); 
            } 	
            event_buffer_disable($this->buffers[$id], EV_READ | EV_WRITE);
            event_buffer_free($this->buffers[$id]);
            fclose($this->connections[$id]);
            unset($this->buffers[$id], $this->connections[$id]);            
		}

		/*!	@function	trigger_hooks
			@abstract	Triggers Hooks for a certain command.
			@param		string	- Command who's hooks you want to trigger.
			@param		object	- The client who activated this command.
			@param		string	- The input from the client, or a message to be sent to the hooks.
			@result		void
		*/
		public function trigger_hooks($command,&$client,$input)
		{
			if(isset($this->hooks[$command]))
			{
				foreach($this->hooks[$command] as $function)
				{
					SocketServer::debug("Triggering Hook '{$function}' for '{$command}'");
					$continue = call_user_func($function,$this,$client,$input);
					if($continue === FALSE) { break; }
				}
			}
		}
 
		/*!	@function	debug
			@static
			@abstract	Outputs Text directly.
			@discussion	Yeah, should probably make a way to turn this off.
			@param		string	- Text to Output
			@result		void
		*/
		public static function debug($text)
		{
			echo("{$text}\r\n");
		}

		/*!	@function	socket_write_smart
			@static
			@abstract	Writes data to the socket, including the length of the data, and ends it with a CRLF unless specified.
			@discussion	It is perfectly valid for socket_write_smart to return zero which means no bytes have been written. Be sure to use the === operator to check for FALSE in case of an error. 
			@param		resource- Socket Instance
			@param		string	- Data to write to the socket.
			@param		string	- Data to end the line with.  Specify a "" if you don't want a line end sent.
			@result		mixed	- Returns the number of bytes successfully written to the socket or FALSE on failure. The error code can be retrieved with socket_last_error(). This code may be passed to socket_strerror() to get a textual explanation of the error.
		*/
		public static function socket_write_smart(&$sock,$string,$crlf = "\r\n")
		{
			SocketServer::debug("<-- {$string}");
			if($crlf) { $string = "{$string}{$crlf}"; }
			return fwrite($sock,$string,strlen($string)); 
      
		}

		/*!	@function	__get
			@abstract	Magic Method used for allowing the reading of protected variables.
			@discussion	You never need to use this method, simply calling $server->variable works because of this method's existence.
			@param		string	- Variable to retrieve
			@result		mixed	- Returns the reference to the variable called.
		*/
		function &__get($name)
		{
			return $this->{$name};
		}
	}
 
 

?>
