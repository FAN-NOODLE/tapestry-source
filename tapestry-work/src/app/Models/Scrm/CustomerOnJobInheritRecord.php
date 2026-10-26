<?php

namespace App\Models\Scrm;

class CustomerOnJobInheritRecord extends ScrmInit
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'customer_onjob_inherit_record';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $guarded = [];

    /**
     * @var string[]
     */
    protected $hidden = ['deleted_at'];

    // 1-接替完毕
    const STATUS_SUCCESS = 1;
    //2-等待接替
    const STATUS_WAIT = 2;
    //3-客户拒绝
    const STATUS_REJECT = 3;
    //4-接替成员客户达到上限',
    const STATUS_UP = 4;
    //5- 当前客户正在转接中
    const STATUS_WAIT_ERROR = 5;
    //6- 24小时没有查询到该接替的信息
    const STATUS_RECORD_ERROR = 6;
}