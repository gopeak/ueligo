<?php

namespace php\engine;

/**
 * 自定义一个异常处理类
 */
class GameException extends \Exception
{
    
    /**
     * 那些用户可以trace
     * @var ArrayIterator
     */
    protected $trace_uids;
    
    /**
	 *是否trace
     */
    protected $is_trace = true;
    
 
    
    public $msg = array();
    
    
    /**
     * 重定义构造器使 message 变为必须被指定的属性
     * @param array $message 结构示例：array('key'=>'CollectionsRewardCoins','content'=>array('coins'=>100));
     *                       或者采用如下结构array('CollectionsRewardCoins',array('coins'=>100));
     * @param $code 异常的代码
     */ 
    public function __construct($messageArr, $code = 0) 
    {
     
        $final['key'] = isset($messageArr['key']) ? $messageArr['key'] : $messageArr[0];
        $final['value'] = isset($messageArr['content']) ? $messageArr['content'] : $messageArr[1];
       // $langModel = new LanguageModel();
       // $lang = $langModel->getLangKey($final['key']);
        if( !isset($lang['level']) ) $lang['level'] = 'none';
        //$final['level'] = $lang['level'];
        $this->msg = $final;
        $this->code = $code;
     
    }
 
}
