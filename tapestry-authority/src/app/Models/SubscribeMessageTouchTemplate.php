<?php

namespace App\Models;

class SubscribeMessageTouchTemplate extends Base
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'subscribe_message_touch_template';

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];

    public function setting()
    {
        return $this->belongsTo(SubscribeMessageSetting::class, 'setting_id', 'id');
    }
}
