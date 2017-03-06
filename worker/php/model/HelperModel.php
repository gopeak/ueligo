<?php

namespace php\model;

/**
 * 协助类
 *
 */
class HelperModel {


	/**
	 * 选择其他数据库
	 */
	static public function selectDB($dbname){
	    $_db_cfg = array(
	        'driver' =>'mysql',
	        'host' =>'192.168.3.10',
	        'port'=>'3306',
	        'user' =>'root',
	        'passwd' =>'max123',
	        'dbname' => $dbname,
	        'charset'=>'utf8',
	        'weight'=>10,
	    );
	    return new PdoModel( $_db_cfg );
	}	

	

	/**
	 * Simple function to replicate PHP 5 behaviour
	 */
	static function microtime_float()
	{
		list ( $usec, $sec ) = explode ( " ", microtime () );
		$usec  =  substr(str_replace ( '0.', '', $usec ),0,-2) ;
		$final =  $sec.$usec ;
		return $final;
	}

	/**
	 * 根据数值表计算用户等级
	 * @param $data
	 * @param $exp
	 */
	static function getLevel( $data,$exp )
	{
		if(empty($data))
		{
			return 1;
		}
		krsort($data);
		foreach ($data as $key => $value)
		{
			if($exp >= $value['exp'])  return $key;
		}
		return 1;
	}

	/**
	 * 将debug信息放置于全局变量中，以便于在框架中捕获
	 * Enter description here ...
	 * @param string $message
	 */
	static function debug($message) {
		$GLOBALS['__debugs'][] = $message;

	}


	/**
	 * 获取随机字符串
	 * @param $len
	 * @param $chars
	 */
	static function rand_string($len, $chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789')
	{
		$string = '';
		for ($i = 0; $i < $len; $i++)
		{
			$pos = rand(0, strlen($chars)-1);
			$string .= $chars{$pos};
		}
		return $string;
	}
 

	/**
	 * 校验手机号
	 * @param string $phone
	 * @return boolean
	 */
	static  public function checkPhone( $phone )
	{
		if(preg_match("/^17[0-9]{9}$|13[0-9]{9}$|15[0-9]{9}$|18[0-9]{9}$/",$phone)){
			//验证通过
			return true;
	
		}else{
			//手机号码格式不对
			return false;
			 
		}
	}

	static public function checkEmail($email)
	{
		if(preg_match("/^([a-zA-Z0-9])+([a-zA-Z0-9\._-])*@([a-zA-Z0-9_-])+([a-zA-Z0-9\._-]+)+$/",
				$email)){
			/*	list($username,$domain)=preg_split('/@/',$email);
			 	
			if(!checkdnsrr($domain,'MX')) {
			return false;
			}
			*/
			return true;
		}
		return false;
	}

	 
	 
	/**
	 *
	 * 9位不重复数字
	 * @param int $begin
	 * @param int $end
	 * @param int $limit
	 */
	static function uniquieNumRand($begin=0,$end=9,$limit=9)
	{
		$rand_array = range($begin,$end);
		shuffle($rand_array);//调用现成的数组随机排列函数
		return array_slice($rand_array,0,$limit);//截取前$limit个
	}

  
 
   

  

	/**
	 * 将秒数转成日期格式
	 * @param int $time
	 * @return array|boolean
	 */
	static function sec2Time( $time ){

		$time = intval( $time );
		if( $time ){
			$value = array(
					"years" => 0, "days" => 0, "hours" => 0,
					"minutes" => 0, "seconds" => 0,
			);
			/*
			 if($time >= 31556926){
			$value["years"] = floor($time/31556926);
			$time = ($time%31556926);
			}
			*/
			if($time >= 86400){
				$value["days"] = floor($time/86400);
				$time = ($time%86400);
			}
			if($time >= 3600){
				$value["hours"] = floor($time/3600);
				$time = ($time%3600);
			}
			if($time >= 60){
				$value["minutes"] = floor($time/60);
				$time = ($time%60);
			}
			$value["seconds"] = floor($time);

			$hour = '';
			if( !empty($value['hours']) ) $hour= $value['hours'].'时';
				
			$need_time =    $hour.$value['minutes'].'分'.$value['seconds'].'.';
			if( $value['days']>0 )
			{
				$need_time = $value['days'].'天'.$need_time;
			}

			return $need_time;
		}else{
			return  '-';
		}
	}

 



	/**
	 * 检查用户具备某一权限
	 * @param string $priv_str 要被检测的权限
	 * @return boolean
	 */
	static function checkAdminPriv( $priv_str ) {

		if( !isset($_SESSION['priv']) ) return false;
		if ( $_SESSION['priv']== '*' ) {
			return true;
		}
		$arr = explode( ',' , $_SESSION['priv']);

		foreach( $arr  as $p){
			if( $p =='*') return true;
			if( $p==$priv_str ) return true;
			list( $module , $act ) = explode( '.' , $p );
			list( $needModule , $needAct ) = explode( '.' , $priv_str );
			if( $module==$needModule && $act=='*' ) return true;

		}

		return false;
	}


	static function substr($string, $length, $option = array()) {
		$strcut = '';
		$strLength = 0;
		$i_option = array('add_dot'=>true, 'charset'=>'utf-8', 'char_len'=>false);
		$option = array_merge($i_option, $option);
		if(strlen($string) > $length) {
			//将$length换算成实际UTF8格式编码下字符串的长度
			for($i = 0; $i < ($length-($option['add_dot']?3:0)); $i++) {
				if ( $strLength >= strlen($string) )
					break;
				//当检测到一个中文字符时
				if( ord($string[$strLength]) > 127 ) {
					if ($option['char_len'] || ++$i < ($length-($option['add_dot']?3:0))) {
						$strLength += (($option['charset'] == 'utf-8')?3:2);
					}
				}
				else
					$strLength += 1;
			}
			return substr($string, 0, $strLength).($option['add_dot']?'...':'');
		} else {
			return $string;
		}
	}

  
	
	/**
	 * 删除一个目录及文件夹
	 */
	static public function deldir($dir) {
		$dh=@opendir($dir);
		while ($file=@readdir($dh)) {
			if($file!="." && $file!="..") {
				$fullpath=$dir."/".$file;
				if(!is_dir($fullpath)) {
					unlink($fullpath);
				} else {
					deldir($fullpath);
				}
			}
		}
	
		closedir($dh);
	
		if(rmdir($dir)) {
			return true;
		} else {
			return false;
		}
	}
	
	
    static function makeDir($destFolder)
	{
		if (! is_dir($destFolder) && $destFolder != './' && $destFolder != '../')
		{
			$dirname = '';
			$folders = explode( '/', $destFolder);
			foreach ($folders as $folder)
			{
				$dirname .= $folder . '/';
				if ($folder != '' && $folder != '.' && $folder != '..' && ! is_dir($dirname))
				{
					mkdir($dirname);
				}
			}
	
			// chmod($destFolder,0777);
		}
	}

	
 
	static public function  timeIntervalWithStartDate( $time )
	{
	    
	    $SECOND = 1;
		$MINUTE = (60 * $SECOND);
	    $HOUR   = (60 * $MINUTE);
		$DAY    = (24 * $HOUR);
	    $MONTH  = (30 * $DAY);
		
	    $delta = time()-$time  ;
	    
	    if ($delta < 1 * $MINUTE)
	    {
	        return $delta. "秒前";
	    }
	    if ($delta < 2 * $MINUTE)
	    {
	        return "1分钟前";
	    }
	    if ($delta < 45 * $MINUTE)
	    {
	        $minutes = floor($delta/$MINUTE);
	        return $minutes."分钟前";
	    }
	    if ($delta < 90 * $MINUTE)
	    {
	        return "1小时前";
	    }
	    if ($delta < 24 * $HOUR)
	    {
	        $hours = floor($delta/$HOUR);
	        return $hours."小时前";
	    }
	    if ($delta < 48 * $HOUR)
	    {
	        return "昨天";
	    }
	    if ($delta < 30 * $DAY)
	    {
	        $days = floor($delta/$DAY);
	        return $days."天前";
	    }
	    if ($delta < 12 * $MONTH)
	    {
	        $months = floor($delta/$MONTH);
	        return $months."个月前";
	    }
	    else
	    {
	        $years = floor($delta/$MONTH/12.0);
	        return $years."年前";
	    }
	}
  
	 
}

