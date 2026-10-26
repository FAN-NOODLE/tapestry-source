<?php

namespace App\Models;

class Announcement extends Base
{
    /**
     * 类型:商城公告
     */
    const TYPE_SHOP = 1;

    /**
     * 类型:功能上线
     */
    const TYPE_FUNCTION_ONLINE = 2;

    /**
     * 类型:注意事项
     */
    const TYPE_NOTICE_THINGS = 3;

    /**
     * 表名称:公告表
     *
     * @var string
     */
    protected $table = 'announcement';

    protected $guarded = [];
}
