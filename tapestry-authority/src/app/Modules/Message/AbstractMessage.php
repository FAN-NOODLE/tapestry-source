<?php

namespace App\Modules\Message;

use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\EnterpriseMessageScene;
use App\Models\MessageScene;
use App\Models\SubscribeMessageRecord;
use App\Models\SubscribeMessageSetting;
use CrsSdk\MicroUser\User;

abstract class AbstractMessage
{
    protected $channel;
    protected $enterpriseId;
    protected $enterprise;
    protected $enterpriseThirdPartAccount;
    protected $records;
    protected $enterpriseMessageScene;
    protected $messageScene;
    protected $userIds = [];
    protected $users = [];
    protected $sceneType;
    protected $message;
    protected $option;
    protected $redis;

    public function __construct()
    {
        $this->redis = redis();
    }

    public function setChannel($value)
    {
        $this->channel = $value;
    }

    public function setEnterpriseId($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
        $this->enterprise   = Enterprise::query()->where('id', '=', $enterpriseId)->first();
        // 只有小程序订阅消息查找订阅记录
        if ($this->channel == MessageFactory::MINI_PROGRAM) {
            $this->enterpriseThirdPartAccount = EnterpriseThirdPartAccount::query()
                ->select(['id', 'enterprise_id', 'app_id', 'app_secret', 'name'])
                ->where('enterprise_id', '=', $enterpriseId)
                ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
                ->where('platform', '=', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
                ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
                ->where('app_id', '!=', '')
                ->where('app_secret', '!=', '')
                ->first();
        }
    }

    public function setSceneType($value)
    {
        $this->sceneType = $value;
    }

    public function setMessage($value)
    {
        $this->message = $value;
    }

    public function setOption($value)
    {
        $this->option = $value;
    }

    public function setRecords($userIds)
    {
        // 只有小程序订阅消息查找订阅记录
        if ($this->channel == MessageFactory::MINI_PROGRAM) {
            SubscribeMessageRecord::query()
                ->where('enterprise_id', '=', $this->enterpriseId)
                ->where('object_type', '=', $this->option['object_type'])
                ->where('object_id', '=', $this->option['object_id'])
                ->where('send_num', 0)
                ->when(0 == $this->option['object_id'], function ($query) {
                    return $query->where('fail_reason', '');
                })
                ->when($userIds, function ($query) use ($userIds) {
                    return $query->whereIn('user_id', $userIds);
                })
                ->whereHas('template', function ($query) {
                    $query->where('scene_type', $this->sceneType)
                        ->when(SubscribeMessageSetting::SCENE_TYPE_GROUP_RESULT == $this->sceneType, function ($query) {
                            $settingIds = ('拼团成功' == $this->message['group_result'])
                                ? SubscribeMessageSetting::query()
                                    ->where('scene_type', SubscribeMessageSetting::SCENE_TYPE_GROUP_RESULT)
                                    ->whereIn('scene_desc', ['拼团结果通知-成功/失败', '拼团结果通知-成功'])
                                    ->pluck('id')
                                    ->all()
                                : SubscribeMessageSetting::query()
                                    ->where('scene_type', SubscribeMessageSetting::SCENE_TYPE_GROUP_RESULT)
                                    ->whereIn('scene_desc', ['拼团结果通知-成功/失败', '拼团结果通知-失败'])
                                    ->pluck('id')
                                    ->all();
                            return $query->whereIn('setting_id', $settingIds);
                        });
                })
                ->with('template.setting')
                ->chunkById(1000, function ($records) {
                    $this->records = collect($this->records)->merge($records);
                });

            $this->enterpriseMessageScene = EnterpriseMessageScene::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->where('message_scene_id', $this->option['type'])
                ->first();

            $this->messageScene = MessageScene::query()
                ->where('id', $this->option['type'])
                ->first();
        }
    }

    /**
     * setUserIds方法需放置在setRecords之后
     *
     * @param $userIds
     * @throws \Exception
     */
    public function setUserIds($userIds)
    {
        $this->userIds = $userIds;

        if (MessageFactory::MINI_PROGRAM == $this->channel && empty($userIds) && $this->records) {
            $this->userIds = $this->records->pluck('user_id')->all();
        }

        $_SERVER['enterprise_hash'] = $this->enterprise['hash'];
        $this->users                = (new User())->userList([
            'is_paginate'   => 0,
            'enterprise_id' => $this->enterpriseId,
            'conditions'    => [['id', 'in', $this->userIds]],
            'cdp'           => 1
        ])['list'];
    }

    protected function getOpenid($userId)
    {
        return array_column($this->users, 'open_id', 'id')[$userId];
    }

    protected function getUnionId($userId)
    {
        return array_column($this->users, 'union_id', 'id')[$userId];
    }

    protected function getCdpCodeId($userId)
    {
        return array_column($this->users, 'cdp_codes', 'id')[$userId];
    }

    /**
     * 获取组装好的消息内容
     *
     * @return array
     */
    abstract protected function getContent(): array;

    /**
     * 发送通知 TODO::异步
     *
     * @return mixed
     */
    abstract public function send();
}
