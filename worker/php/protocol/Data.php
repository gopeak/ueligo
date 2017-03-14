<?php
// automatically generated, do not modify

namespace protocol;

use \Google\FlatBuffers\Struct;
use \Google\FlatBuffers\Table;
use \Google\FlatBuffers\ByteBuffer;
use \Google\FlatBuffers\FlatBufferBuilder;

class Data extends Table
{
    /**
     * @param ByteBuffer $bb
     * @return Data
     */
    public static function getRootAsData(ByteBuffer $bb)
    {
        $obj = new Data();
        return ($obj->init($bb->getInt($bb->getPosition()) + $bb->getPosition(), $bb));
    }

    /**
     * @param int $_i offset
     * @param ByteBuffer $_bb
     * @return Data
     **/
    public function init($_i, ByteBuffer $_bb)
    {
        $this->bb_pos = $_i;
        $this->bb = $_bb;
        return $this;
    }

    /**
     * @return short
     */
    public function get_type()
    {
        $o = $this->__offset(4);
        return $o != 0 ? $this->bb->getShort($o + $this->bb_pos) : \protocol\Type::Req;
    }

    public function getCmd()
    {
        $o = $this->__offset(6);
        return $o != 0 ? $this->__string($o + $this->bb_pos) : null;
    }

    public function getSid()
    {
        $o = $this->__offset(8);
        return $o != 0 ? $this->__string($o + $this->bb_pos) : null;
    }

    /**
     * @return int
     */
    public function getReqId()
    {
        $o = $this->__offset(10);
        return $o != 0 ? $this->bb->getInt($o + $this->bb_pos) : 0;
    }

    public function getData()
    {
        $o = $this->__offset(12);
        return $o != 0 ? $this->__string($o + $this->bb_pos) : null;
    }

    public function getToken()
    {
        $o = $this->__offset(14);
        return $o != 0 ? $this->__string($o + $this->bb_pos) : null;
    }

    /**
     * @param FlatBufferBuilder $builder
     * @return void
     */
    public static function startData(FlatBufferBuilder $builder)
    {
        $builder->StartObject(6);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @return Data
     */
    public static function createData(FlatBufferBuilder $builder, $_type, $cmd, $sid, $req_id, $data, $token)
    {
        $builder->startObject(6);
        self::add_type($builder, $_type);
        self::addCmd($builder, $cmd);
        self::addSid($builder, $sid);
        self::addReqId($builder, $req_id);
        self::addData($builder, $data);
        self::addToken($builder, $token);
        $o = $builder->endObject();
        return $o;
    }

    /**
     * @param FlatBufferBuilder $builder
     * @param short
     * @return void
     */
    public static function add_type(FlatBufferBuilder $builder, $Type)
    {
        $builder->addShortX(0, $Type, 1);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @param StringOffset
     * @return void
     */
    public static function addCmd(FlatBufferBuilder $builder, $cmd)
    {
        $builder->addOffsetX(1, $cmd, 0);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @param StringOffset
     * @return void
     */
    public static function addSid(FlatBufferBuilder $builder, $sid)
    {
        $builder->addOffsetX(2, $sid, 0);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @param int
     * @return void
     */
    public static function addReqId(FlatBufferBuilder $builder, $reqId)
    {
        $builder->addIntX(3, $reqId, 0);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @param StringOffset
     * @return void
     */
    public static function addData(FlatBufferBuilder $builder, $data)
    {
        $builder->addOffsetX(4, $data, 0);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @param StringOffset
     * @return void
     */
    public static function addToken(FlatBufferBuilder $builder, $token)
    {
        $builder->addOffsetX(5, $token, 0);
    }

    /**
     * @param FlatBufferBuilder $builder
     * @return int table offset
     */
    public static function endData(FlatBufferBuilder $builder)
    {
        $o = $builder->endObject();
        return $o;
    }

    public static function finishDataBuffer(FlatBufferBuilder $builder, $offset)
    {
        $builder->finish($offset);
    }
}
