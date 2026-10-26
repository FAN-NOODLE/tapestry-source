<?php

namespace App\Models;

class MessageScene extends Base
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'message_scene';

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];

    public function channels()
    {
        return $this->hasMany(MessageChannel::class,'message_scene_id','id');
    }

    public function enterpriseMessageScenes()
    {
        return $this->hasMany(EnterpriseMessageScene::class,'message_scene_id','id');
    }
}
