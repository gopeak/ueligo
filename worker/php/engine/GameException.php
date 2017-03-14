<?php

namespace php\engine;

/**
 * 自定义一个异常处理类
 */
class GameException extends \Exception
{
  
    
    
    /**
     * 重定义构造器使 message 变为必须被指定的属性 
     */ 
    public function __construct($message, $code = 0) 
    {
        parent::__construct($message,$code); 

        $this->message = $message; 
     
    }
 
}
