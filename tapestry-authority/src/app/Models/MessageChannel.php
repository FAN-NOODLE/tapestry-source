<?php

namespace App\Models;

class MessageChannel extends Base
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'message_channel';

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];

    public function settings()
    {
        return $this->hasMany(SubscribeMessageSetting::class, 'scene_type', 'scene_type');
    }
}
