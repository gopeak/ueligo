<?php

namespace php\model {


	/**
	 * 用户模块逻辑实现类
	 *
	 *
	 * @author Seven@mcross.cn
	 *
	 */
	class UserModel extends BaseCacheDbModel
	{

		public  $table = 'user_account';


		public $fields = ' * ';


		const  DATA_KEY = 'user_account/';

		const  REG_RETURN_CODE_OK    = 1;
		const  REG_RETURN_CODE_EXIST = 2;
		const  REG_RETURN_CODE_ERROR = 3;

		const  LOGIN_CODE_OK    = 1;
		const  LOGIN_CODE_EXIST = 2;
		const  LOGIN_CODE_ERROR = 3;


		public $uid = '';

		function __construct( $uid ='',$PERSISTENT=false )
		{
			parent::__construct( $uid,$PERSISTENT );

			$this->uid = $uid;

		}


		/**
		 * 取得一个用户的基本信息
		 * @param $noCache 是否必须从数据库中查询数据
		 * @return array
		 */
		public function getUser( $noCache=false )
		{

			$uid = $this->uid;
			$fileds	=	'*';
			$where	=	" WHERE uid='$uid'";
			$key	=	self::DATA_KEY.$uid;

			if($noCache==false)
			{
				$finally = $this->getRowByKey($this->table,$fileds,$where,$key);
			}
			else
			{
				//必须从数据库中查询以取得数据
				$sql = "SELECT  {$fileds} FROM {$this->table} $where ";
				$finally = $this->masterDB->getRow($sql);
			}

			return  $finally;
		}
		/**
		 * 根据uid查询用户
		 * @param int $uid
		 * @return Ambigous <multitype:, unknown>
		 * @author 秋士悲
		 */
		public function getUserById($uid){
			$table		= $this->table;
			$fileds		= '*';
			$where		= "WHERE uid='$uid'";
			$key		= '';
			$result		= parent::getRowByKey($table, $fileds, $where, $key);
			return $result;
		}

		/**
		 * 取得一个用户的基本信息,通过Email地址
		 * @param $email Email地址
		 * @return array
		 */
		public function getUserByEmail( $email )
		{

			$table	= $this->table;
			$fileds	=	'*,uid as k';
			$where	=	" Where `email`='$email'   ";
			$key	=	self::DATA_KEY.'email_'.$email;
			$user	=	parent::getRowByKey($this->table,$fileds,$where,$key);
			return  $user;
		}
		/**
		 * 取得一个用户的基本信息,通过帐号地址
		 * @param $username 帐号
		 * @return array
		 */
		public function getUserByUser( $username )
		{

			$table	= $this->table;
			$fileds	=	'*,uid as k';
			$where	=	" Where `username`='$username'   ";
			$key	=	self::DATA_KEY.'username/'.$username;
			$user	=	parent::getRowByKey($this->table,$fileds,$where,$key);
			return  $user;
		}



		/**
		 * 取得一个用户的基本信息,通过Email地址
		 * @param $email Email地址
		 * @return array
		 */
		public function getUserByOpenid( $open_id )
		{

			$table	= $this->table;
			$fileds	=	'*,uid as k';
			$where	=	" Where `open_id`='$open_id' ";
			$key	=	self::DATA_KEY."openid/$open_id";
			$user	=	parent::getRowByKey($this->table,$fileds,$where,$key);
			return  $user;
		}

		public function checkLoginByEmail( $email, $password )
		{

			$user	=	$this->getUserByEmail( $email );
			//var_dump($user);
			if( !isset($user['password']) )
			{
				return array( self::LOGIN_CODE_EXIST ,$user);
			}

			if(  $user['password']!=md5($password) )
			{
				return array( self::LOGIN_CODE_ERROR, $user );
			}

			return  array( self::LOGIN_CODE_OK, $user );

		}


		public function checkLoginByUsername( $username, $password )
		{

			$user	=	$this->getUserByUser( $username );
			//var_dump($user);
			if( !isset($user['password']) )
			{
				return array( self::LOGIN_CODE_EXIST ,$user);
			}

			if(  $user['password']!=md5($password) )
			{
				return array( self::LOGIN_CODE_ERROR, $user );
			}

			return  array( self::LOGIN_CODE_OK, $user );

		}





		/**
		 * 通过呢称检查密码是否正确
		 * @param string $nick
		 * @param string $password
		 * @return number
		 */
		public function checkLoginByNick( $nick, $password )
		{
			$table	= $this->table;
			$fileds	=	'password';
			$where	=	" Where `nick`='$nick'   ";
			$key	=	'';
			$user	=	parent::getRowByKey($this->table,$fileds,$where,$key);

			if( !isset($user['password']) )
			{
				return array( 2, $user );
			}

			if(  $user['password']!=$password )
			{
				return array( 3 ,$user );
			}

			return array( 1, $user );
		}





		/**
		 * 获得用户总数
		 * @return array
		 */
		public function getCountUser(  )
		{

			$table	= $this->table;
			$fileds	=	'count(uid) as _sum';
			$where	=	"   ";
			$key	=	self::DATA_KEY.'count';
			$user	=	parent::getRowByKey( $this->table, $fileds, $where, $key );
			 
			return intval( $user['_sum'] );
		}




		/**
		 * 取得多个用户的所有信息
		 * @param $uids
		 * @return array
		 */
		public function getUserByids( $uids )
		{
			$uids_str = "''";
			if(empty($uids))
			{
				return array();
			}
			if( !empty($uids)  && is_array($uids))
			{
				$uids_str = implode(',', $uids);
			}
			if( is_string($uids) )
			{
				$uids_str = $uids;
				$uids = explode(',', $uids_str);
			}
			$fileds	=	'*';
			$where	=	" WHERE uid in ($uids_str)";
			$key	=	'users_'.$this->uid;
			$primkey=	"uid";
			$keys   = array();
			//为每个uid增加缓存的前缀

			$finally = $this->Cache->mget($keys);

			//如果返回False则直接查询数据库
			if( empty( $finally ) )
			{
				foreach ($uids as $key => $uid)
				{
					$this->uid = $uid;
					$finally[$uid] = $this->getUser();
				}

			} else {
				//补齐找不到的键值
				if( count($finally)<count($uids) )
				{
					foreach ($uids as $uid)
					{
						$uid = (string) $uid;
						if( !isset($finally[$uid]))
						{
							$this->uid = $uid;
							$finally[$uid] = $this->getUser();
						}
					}
				}
			}

			return  $finally;
		}



		/**
		 * 更新一个用户的信息
		 * @param $userid   用户ID
		 * @param $updateinfo array 可同时更新多个字段值,如 array('u_name'=>'马柱国','u_photo'=>'http://www')
		 * @return bool
		 */
		public function updateUser($updateinfo)
		{
			if(empty($updateinfo))
			{
				return false;
			}
			if(!is_array($updateinfo))
			{
				return false;
			}
			$uid  = $this->uid;
			$key  = self::DATA_KEY.$uid;
			$where= "  where `uid`='$uid'";
			$flag =$this->updateInfoByKey($this->table,$where,$updateinfo,$key);
			if($flag)
			{
				if( isset($_POST['_userinfo'][$uid]) )
				{
					$_POST['_userinfo'][$uid] = array_merge($_POST['_userinfo'][$uid],$updateinfo);
				}
				 
			}
			return 	$flag;
		}

		/**
		 * 更新用户信息
		 * @param integer $uid
		 * @param array $row
		 * @return boolean
		 */
		public function updateUserInfo($uid,$row){
			$table			= $this->table;
			$where			= "WHERE uid=$uid";
			$result			= parent::updateInfoByKey($table, $where, $row);
			return $result;
		}

		/**
		 * 增加用户未读消息的数量
		 * @param integer $uid
		 * @return boolean;
		 * @author 秋士悲
		 */
		public function plusUnreadMessage($uid){
			$table			= $this->table;
			$where			= "WHERE uid=$uid";
			$sql			= "UPDATE $table SET unreaded_msg=unreaded_msg+1 $where";
			$result			= $this->masterDB->query($sql);
			return true;
		}

		/**
		 * 减少用户未读消息的数量
		 * @param integer $uid
		 * @return boolean
		 * @author 秋士悲
		 */
		public function minusUnreadMessage($uid){
			$table			= $this->table;
			$where			= "WHERE uid=$uid";
			$sql			= "UPDATE $table SET unreaded_msg= CASE WHEN unreaded_msg>0 THEN unreaded_msg-1 END $where";
			$result			= $this->masterDB->query($sql);
			return true;
		}

		/**
		 * 只更新用户的某个字段
		 * @param $field   字段名
		 * @param $value   值
		 * @return bool
		 */
		public function updateUserField($field,$value)
		{
			if( empty($field) )
			{
				return false;
			}
			$uid   = $this->uid;
			$key   = self::DATA_KEY.$uid;
			$where = "  where `uid`='$uid'";
			$updateinfo = array($field=>$value);
			$flag = $this->updateInfoByKey($this->table,$where,$updateinfo,$key);
			if($flag)
			{
				if(isset($_POST['_userinfo'][$uid]))
				{
					$_POST['_userinfo'][$uid][$field] =  $value;
				}

			}

			return 	$flag;
		}

		/**
		 * 添加用户
		 * @param $userinfo   提交的用户信息
		 * @return bool
		 */
		public function addUser( $userinfo,$is_open_login=false )
		{
			if(empty($userinfo))
			{
				return array( self::REG_RETURN_CODE_ERROR,array() );
			}
			if( !$is_open_login ) {
				$user	=	$this->getUserByUser( $userinfo['username'] );
				//v($userinfo);
				if( isset($user['username'])  )
				{
					return array( self::REG_RETURN_CODE_EXIST, $user );
				}
			}

			$no   = mt_rand(10, 10000).substr( strval( time()) ,-4 );;
			$user = $userinfo;


			if( isset($userinfo['uid']) ) $user['uid'] = $userinfo['uid'];
			$flag = $this->insertInfoByKey( $this->table,$user, '' );
			//v($flag);
			if($flag)
			{
				//v($user);
				$uid = $this->masterDB->getLastInsId();
				$this->uid = $uid;
				$user = $this->getUser(true);


			}
			return  array(self::REG_RETURN_CODE_OK,$user);

		}


		/**
		 * 删除用户
		 * @param $uid
		 * @return bool
		 */
		public function deleteUser()
		{
			$uid = $this->uid;
			$key        =  self::DATA_KEY.$uid;
			$where 		=  " WHERE uid='$uid'";
			$flag       =  $this->deleteBykey($this->table,$where,$key);
			if($flag)
			{
				if(isset($_POST['_userinfo'][$uid]))
				{
					unset($_POST['_userinfo'][$uid]);
				}
			}
			return  $flag;
		}



		/**
		 * 为数组的每一个元素增加key前缀
		 * @param $prefixKey
		 * @param $value
		 */
		static function addUserPrefixKey($value)
		{
			return self::DATA_KEY.$value;
		}





	}

}