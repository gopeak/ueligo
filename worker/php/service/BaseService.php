<?php

namespace php\service {

	/**
	 * service 的基类文件主要用于接口的基础服务
	 *
	 * @author seven
	 *
	 */
	class BaseService
	{
 
		public $sid;
        
        public $cmd;
        
        public $req_id;
        
        public $conn;

		function __construct( $sid, $cmd, $req_id, $conn )
		{
			$this->sid =  $sid;
            $this->cmd =  $cmd;
            $this->req_id =  $req_id;
            $this->conn =  $conn;
            
		}

 
	 

	}

}