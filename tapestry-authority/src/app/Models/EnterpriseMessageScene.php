<?php

namespace App\Models;

class EnterpriseMessageScene extends Base
{

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_message_scene';

    protected $guarded = [];


    public function belongsToMessageScene()
    {
        return $this->belongsTo(MessageScene::class, 'message_scene_id','id');
    }

    public function channels()
    {
        return $this->hasMany(MessageChannel::class,'message_scene_id','message_scene_id');
    }

    public function activityMessage()
    {
        return $this->hasMany(EnterpriseActivityMessage::class,'message_scene_id','message_scene_id');
    }
}