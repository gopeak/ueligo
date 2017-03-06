<?php

/**
 * This file is part of zeroman.
 *
 * Licensed under The MIT License
 * For full copyright and license information, please see the MIT-LICENSE.txt
 * Redistributions of files must retain the above copyright notice.
 *
 * @author sven@
 * @link http://www.zeromore.net/ 
 */
define ( 'SDK_PATH', realpath ( dirname ( __FILE__ ) ) );
class Application {
	public $config;
	
	/**
	 * 用于实现当前类的单例模式
	 *
	 * @var self
	 */
	protected static $_instance;
	
	/**
	 * 用于实现ZMQ_SOCKET的单例模式
	 *
	 * @var self
	 */
	protected static $_zmq_instance;
	
	/**
	 * 用于实现 SOCKET pfsockopen 的单例模式
	 *
	 * @var self
	 */
	protected static $_fd_instance;
	
	/**
	 * 创建一个ZMQ_SOCKET的对象
	 *
	 * @return self
	 */
	public static function getFdInstance($config, $force_connect = false) {
		// 无法实现单例(持久对象)模式，否则会抛出异常
		if (is_null ( self::$_fd_instance ) || ! is_resource ( self::$_fd_instance ) || $force_connect) {
			
			$hub_host = $config ['hub'] ['host'];
			$hub_port = $config ['hub'] ['port'];
			echo "fsockopen:" . $hub_host . $hub_port . "\n";
			$fp2 = fsockopen ( $hub_host, $hub_port, $errno, $errstr, 5 );
			if (! $fp2) {
				echo "fsockopen :$errstr ($errno)<br />\n";
			}
			self::$_fd_instance = $fp2;
		}
		// var_dump( 'self::$_fd_instance:', $fp2 );
		return self::$_fd_instance;
	}
	
	/**
	 * 创建一个自身的单例对象
	 *
	 * @return self
	 */
	public static function getInstance() {
		if (! isset ( self::$_instance ) || ! is_object ( self::$_instance )) {
			
			self::$_instance = new self ();
		}
		return self::$_instance;
	}
	public function __construct() {
		$this->config = $this->getServerConfig ();
	}
	
	/**
	 * 获取服务器的配置文件
	 *
	 * @return array
	 */
	public function getServerConfig() {
		
		$cfgArr = json_decode ( file_get_contents ( SDK_PATH . '/../../../config.json' ), true );
		if (empty ( $cfgArr ))
			$cfgArr = array ();
		
		return $cfgArr;
	}
	
	/**
	 * 获取服务器的绝对路径
	 *
	 * @return string
	 */
	public function getServerBase() {
		return realpath ( '../../../' );
	}
	
	/**
	 * 获取php的绝对路径
	 */
	public function getPhpBase() {
		return realpath ( '../' );
	}
	
	/**
	 * 获取服务器状态
	 */
	public function getEnable() {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		$final ['data'] = '';
		
		$app = Application::getInstance ();
		
		$msg = array ();
		$msg ['cmd'] = 'get_enable';
		$data = $app->request_hub ( json_encode ( $msg ) );
		
		$final ['data'] = $data;
		$final ['msg'] = '';
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 设置服务器可用
	 *
	 * @return array
	 */
	public function enabled() {
		$final ['ret'] = 0;
		
		$fd = Application::getFdInstance ( $this->config );
		$msg ['cmd'] = 'enabled';
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'msg sended!';
		$final ['ret'] = 1;
		
		return $final;
	}
	
	/**
	 * 设置服务器不可用
	 *
	 * @return array
	 */
	public function disabled() {
		$final ['ret'] = 0;
		
		$fd = Application::getFdInstance ( $this->config );
		$msg ['cmd'] = 'disabled';
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'msg sended!';
		$final ['ret'] = 1;
		
		return $final;
	}
	
	/**
	 * 获取服务器的所有场景列表
	 *
	 * @return array
	 */
	public function request_hub($data, $timeout = 5, $blocking = true) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		$final ['data'] = '';
		$config = $this->config;
		$hub_host = $config ['hub'] ['host'];
		$hub_port = $config ['hub'] ['port'];
		
		$fp3 = fsockopen ( $hub_host, $hub_port, $errno, $errstr );
		if (! $fp3) {
			$final ['msg'] = $errno . " " . $errstr;
			$final ['ret'] = 0;
			return $final;
		}
		
		// echo "request_hub request: ".$data."\n";
		fwrite ( $fp3, $data . "\n" );
		
		// tream_set_blocking( $fp3, $blocking );
		stream_set_timeout ( $fp3, $timeout );
		$ret = '';
		// sleep( 1 );
		while ( ! feof ( $fp3 ) ) {
			$ret .= fgets ( $fp3, 4096 );
			break;
		}
		fclose ( $fp3 );
		// echo "request_hub response: ".$ret."\n";
		return $ret;
	}
}

/**
 * 场景控制类
 *
 * @author Administrator
 *        
 */
class ChannelService {
	
	/**
	 * 服务器的配置信息
	 *
	 * @var array
	 */
	public $config;
	
	/**
	 * 用于实现单例模式
	 *
	 * @var self
	 */
	protected static $_instance;
	
	/**
	 * 创建一个自身的单例对象
	 *
	 * @return self
	 */
	public static function getInstance($config) {
		if (! isset ( self::$_instance ) || ! is_object ( self::$_instance )) {
			
			self::$_instance = new self ( $config );
		}
		return self::$_instance;
	}
	
	/**
	 * 构造函数
	 *
	 * @param array $config        	
	 */
	public function __construct($config) {
		$this->config = $config;
	}
	
	/**
	 * 向某一场景内的所有用户发送一条广播消息
	 *
	 * @param mix $data        	
	 * @param string $channel
	 *        	场景名称
	 * @return array
	 */
	public function broadcast($data, $channel = '') {
		$final ['ret'] = 0;
		if (empty ( $channel ))
			$channel = $this->config ['area'] [0] ['id'];
		
		$fd = Application::getFdInstance ( $this->config );
		$msg = array ();
		$msg ['data'] = $data;
		$msg ['id'] = $channel;
		$msg ['cmd'] = 'broatcast';
		$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
		if ($wret === false) {
			$fd = Application::getFdInstance ( $this->config, true );
		}
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
		}
		
		$final ['msg'] = 'msg sended!';
		$final ['ret'] = 1;
		
		return $final;
	}
	
	/**
	 * 在服务器上创建一个场景
	 *
	 * @param string $name
	 *        	场景名称
	 * @param string $dsn
	 *        	场景绑定dsn,如tcp://10.0.1.2:8888 ipc://server_channel1
	 * @return array
	 */
	public function createChannel($name, $dsn = '') {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		
		$fd = Application::getFdInstance ( $this->config );
		// var_dump( $fd );
		$msg = array ();
		$msg ['name'] = $name;
		$msg ['dsn'] = $dsn;
		$msg ['cmd'] = 'create_channel';
		
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'ok';
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 在服务器上删除一个场景
	 *
	 * @param string $name        	
	 * @return array
	 */
	public function removeChannel($name) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		
		$fd = Application::getFdInstance ( $this->config );
		$msg = array ();
		$msg ['name'] = $name;
		$msg ['cmd'] = 'remove_channel';
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'ok';
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 使某个会话用户加入到场景中
	 *
	 * @param string $sid        	
	 * @param string $name        	
	 * @return array
	 */
	public function joinChannel($sid, $name) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		if (empty ( $sid ) || empty ( $name ))
			return $final;
		
		$fd = Application::getFdInstance ( $this->config );
		$msg = array ();
		$msg ['name'] = $name;
		$msg ['sid'] = $sid;
		$msg ['cmd'] = 'join_channel';
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'ok';
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 使某个会话用户离开场景
	 *
	 * @param string $sid        	
	 * @param string $name        	
	 * @return array
	 */
	public function leaveChannel($sid, $name) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		if (empty ( $sid ) || empty ( $name ))
			return $final;
		
		$fd = Application::getFdInstance ( $this->config );
		$msg = array ();
		$msg ['name'] = $name;
		$msg ['sid'] = $sid;
		$msg ['cmd'] = 'leave_channel';
		if ($fd) {
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'ok';
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 获取服务器的所有场景列表
	 *
	 * @return array
	 */
	public function getChannels() {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		$final ['data'] = '';
		
		$app = Application::getInstance ();
		
		$msg = array ();
		$msg ['cmd'] = 'get_channels';
		$data = $app->request_hub ( json_encode ( $msg ) );
		// v( "getChannels", $data );
		$final ['msg'] = $data;
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 获取会话用户已经加入的场景列表
	 *
	 * @param string $sid        	
	 * @return array
	 */
	public function getUserJoinChannels($sid) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		
		$msg = array ();
		$msg ['cmd'] = 'get_user_join_channels';
		$msg ['sid'] = $sid;
		$app = Application::getInstance ();
		$data = $app->request_hub ( json_encode ( $msg ) );
		
		$final ['msg'] = $data;
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 向某个会话用户发送一条消息
	 *
	 * @param string $sid        	
	 * @param string $msg        	
	 * @param string $from_sid
	 *        	消息发送者sid
	 * @param string $req_id
	 *        	请求资源标识
	 * @return number
	 */
	public function push($sid, $msg, $from_sid, $req_id) {
		$final ['ret'] = 0;
		
		try {
			
			$fd = Application::getFdInstance ( $this->config );
			$msg = array ();
			$msg ['data'] = $msg ? $msg : '';
			$msg ['sid'] = $sid;
			$msg ['from_sid'] = $from_sid;
			$msg ['req_id'] = $req_id;
			$msg ['cmd'] = 'push';
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} catch ( Exception $e ) {
			
			$final ['msg'] = $e->getMessage ();
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'ok';
		$final ['ret'] = 1;
		
		return $final;
	}
	
	/**
	 * 向多个会话用户发送一条消息
	 *
	 * @param array $sids        	
	 * @param string $msg        	
	 * @return array
	 */
	public function pushBySids($sids, $msg) {
		$final ['ret'] = 0;
		
		foreach ( $sids as $sid ) {
			
			$fd = Application::getFdInstance ( $this->config );
			$msg = array ();
			$msg ['data'] = $msg ? $msg : '';
			$msg ['sid'] = $sid;
			$msg ['from_sid'] = '';
			$msg ['req_id'] = 0;
			$msg ['cmd'] = 'push';
			
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
				if ($wret === false) {
					$fd = Application::getFdInstance ( $this->config, true );
				}
				if ($fd) {
					$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
				}
			} else {
				$final ['msg'] = 'hub socket null';
				$final ['ret'] = 0;
				return $final;
			}
			usleep ( 10000 );
		}
		
		$final ['msg'] = 'ok';
		$final ['ret'] = 1;
		
		return $final;
	}
}
class BackendDataService {
	public $redis;
	
	/**
	 * 用于实现单例模式
	 *
	 * @var self
	 */
	protected static $_instance;
	
	/**
	 * 创建一个自身的单例对象
	 *
	 * @return self
	 */
	public static function getInstance($redis) {
		if (! isset ( self::$_instance ) || ! is_object ( self::$_instance )) {
			
			self::$_instance = new self ( $redis );
		}
		return self::$_instance;
	}
	public function __construct($redis) {
		$this->redis = $redis;
	}
	public function set($key, $value, $expire = 0) {
	}
	public function get($key) {
	}
	public function delete($key) {
	}
}

/**
 * 会话操作类
 *
 * @author Administrator
 *        
 */
class SessionService {
	
	/**
	 * 服务器配置信息
	 *
	 * @var array
	 */
	public $config;
	
	/**
	 * 用于实现单例模式
	 *
	 * @var self
	 */
	protected static $_instance;
	
	/**
	 * 创建一个自身的单例对象
	 *
	 * @return self
	 */
	public static function getInstance($config, $prefix = 'REDISSESSID_') {
		if (! isset ( self::$_instance ) || ! is_object ( self::$_instance )) {
			
			self::$_instance = new self ( $config, $prefix );
		}
		
		return self::$_instance;
	}
	public function __construct($config) {
		$this->config = $config;
	}
	
	/**
	 * 获取所有的用户会话
	 *
	 * @return arary
	 */
	public function get_all_session() {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		
		$app = Application::getInstance ();
		
		$msg = array ();
		$msg ['cmd'] = 'get_all_session';
		$data = $app->request_hub ( json_encode ( $msg ) );
		
		$final ['msg'] = $data;
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 更新会话中的用户信息
	 *
	 * @param string $sid        	
	 * @param string $user
	 *        	json格式的用户信息,结构可以自定义
	 * @return array
	 */
	public function updateUserBySid($sid, $user) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		if (empty ( $sid ) || empty ( $name ))
			return $final;
		
		$fd = Application::getFdInstance ( $this->config );
		$msg = array ();
		$msg ['data'] = $user;
		$msg ['sid'] = $sid;
		$msg ['cmd'] = 'update_session';
		if ($fd) {
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
				if ($wret === false) {
					$fd = Application::getFdInstance ( $this->config, true );
				}
				if ($fd) {
					$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
				}
			} else {
				$final ['msg'] = 'hub socket null';
				$final ['ret'] = 0;
				return $final;
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'session info updated!';
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 获取会话用户的信息
	 *
	 * @param string $sid        	
	 */
	public function getBySid($sid) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		
		$app = Application::getInstance ();
		
		$msg = array ();
		$msg ['cmd'] = 'get_session';
		$msg ['sid'] = $sid;
		$data = $app->request_hub ( json_encode ( $msg ) );
		
		$final ['msg'] = $data;
		$final ['ret'] = 1;
		return $final;
	}
	
	/**
	 * 从服务器上剔除用户
	 *
	 * @param string $sid        	
	 */
	public function kickBySid($sid) {
		$final ['ret'] = 0;
		$final ['msg'] = '';
		if (empty ( $sid ) || empty ( $name ))
			return $final;
		
		$fd = Application::getFdInstance ( $this->config );
		$msg = array ();
		$msg ['sid'] = $sid;
		$msg ['cmd'] = 'kick';
		if ($fd) {
			
			$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			if ($wret === false) {
				$fd = Application::getFdInstance ( $this->config, true );
			}
			if ($fd) {
				$wret = fwrite ( $fd, json_encode ( $msg ) . "\n" );
			}
		} else {
			$final ['msg'] = 'hub socket null';
			$final ['ret'] = 0;
			return $final;
		}
		
		$final ['msg'] = 'sid kicked!';
		$final ['ret'] = 1;
		return $final;
	}
}







