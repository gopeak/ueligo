<?php
// automatically generated, do not modify

namespace protocol;

class Type
{
    const Req = 1;
    const Reply = 2;
    const Push = 3;
    const Broadcast = 4;

    private static $names = array(
        "Req",
        "Reply",
        "Push",
        "Broadcast",
    );

    public static function Name($e)
    {
        if (!isset(self::$names[$e])) {
            throw new \Exception();
        }
        return self::$names[$e];
    }
}
