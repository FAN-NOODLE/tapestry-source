<?php

namespace App\Models;

class Channel extends Base
{
    /**
     * 数据来源
     */
    const SOURCE_TYPE_DEFAULT = 0; //默认正常创建
    const SOURCE_TYPE_AUTO = 1; //内购员工自建

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'channel';


    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];
}
