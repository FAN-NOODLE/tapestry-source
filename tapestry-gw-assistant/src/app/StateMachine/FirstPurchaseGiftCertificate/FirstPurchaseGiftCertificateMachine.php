<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/13
 * Time: 16:51
 */

namespace App\StateMachine\FirstPurchaseGiftCertificate;


use App\Exceptions\Error;
use App\Services\FirstPurchaseGiftService;
use App\StateMachine\FirstPurchaseGiftCertificate\Actions\Invalid;
use App\StateMachine\FirstPurchaseGiftCertificate\Actions\ManualFailedToSuccess;
use App\StateMachine\FirstPurchaseGiftCertificate\Actions\NoticePushFailed;
use App\StateMachine\FirstPurchaseGiftCertificate\Actions\NoticePushSuccess;
use App\StateMachine\FirstPurchaseGiftCertificate\Actions\PushToOpen;
use App\StateMachine\StateMachine;
use App\StateMachine\Transition;
use Illuminate\Support\Facades\Redis;
use Illuminate\Support\Facades\Validator;

class FirstPurchaseGiftCertificateMachine extends StateMachine
{
    /**
     * 防并发，redis加锁控制
     *
     * @var
     */
    protected $redisKey;

    /**
     * 操作者
     *
     * @var
     */
    const CONSOLE_OPERATOR = 'console';
    const USER_OPERATOR = 'user';

    public function __construct()
    {
        //可以有的状态
        $allowState = array(
            FirstPurchaseGiftCertificateStateConfig::WAIT_USE     => ['desc' => "待使用"],
            FirstPurchaseGiftCertificateStateConfig::PUSHING      => ['desc' => "推送中"],
            FirstPurchaseGiftCertificateStateConfig::PUSH_SUCCESS => ['desc' => "推送成功"],
            FirstPurchaseGiftCertificateStateConfig::PUSH_FAILED  => ['desc' => "推送失败"],
            FirstPurchaseGiftCertificateStateConfig::INVALID      => ['desc' => "已失效"],
        );

        $events = array(
            FirstPurchaseGiftCertificateStateConfig::EVENT_CREATE_GIFT_CERTIFICATE => ['desc' => "创建首购礼资格"],
            FirstPurchaseGiftCertificateStateConfig::EVENT_PUSH_TO_OPEN            => ['desc' => "已推送到open"],
            FirstPurchaseGiftCertificateStateConfig::EVENT_AUTO_INVALID            => ['desc' => "自动失效"],
            FirstPurchaseGiftCertificateStateConfig::EVENT_MANUAL_INVALID          => ['desc' => "人工处理成失效"],
            FirstPurchaseGiftCertificateStateConfig::EVENT_OPEN_NOTICE_SUCCESS     => ['desc' => "open通知推送成功"],
            FirstPurchaseGiftCertificateStateConfig::EVENT_OPEN_NOTICE_FAILED      => ['desc' => "open通知推送失败"],
            FirstPurchaseGiftCertificateStateConfig::EVENT_SECOND_MOBILE_ERROR     => ['desc' => "第二单收货人手机号与会员手机号不符"],
        );

        $transitions = array(

            /*******待使用********/
            new Transition(FirstPurchaseGiftCertificateStateConfig::WAIT_USE, FirstPurchaseGiftCertificateStateConfig::EVENT_PUSH_TO_OPEN, FirstPurchaseGiftCertificateStateConfig::PUSHING, new PushToOpen()),
            new Transition(FirstPurchaseGiftCertificateStateConfig::WAIT_USE, FirstPurchaseGiftCertificateStateConfig::EVENT_MANUAL_INVALID, FirstPurchaseGiftCertificateStateConfig::INVALID, new Invalid()),
            new Transition(FirstPurchaseGiftCertificateStateConfig::WAIT_USE, FirstPurchaseGiftCertificateStateConfig::EVENT_AUTO_INVALID, FirstPurchaseGiftCertificateStateConfig::INVALID, new Invalid()),
            new Transition(FirstPurchaseGiftCertificateStateConfig::WAIT_USE, FirstPurchaseGiftCertificateStateConfig::EVENT_SECOND_MOBILE_ERROR, FirstPurchaseGiftCertificateStateConfig::INVALID, new Invalid()),
            new Transition(FirstPurchaseGiftCertificateStateConfig::WAIT_USE, FirstPurchaseGiftCertificateStateConfig::EVENT_OPEN_NOTICE_SUCCESS, FirstPurchaseGiftCertificateStateConfig::PUSH_SUCCESS, new NoticePushSuccess()),
            new Transition(FirstPurchaseGiftCertificateStateConfig::WAIT_USE, FirstPurchaseGiftCertificateStateConfig::EVENT_OPEN_NOTICE_FAILED, FirstPurchaseGiftCertificateStateConfig::PUSH_FAILED, new NoticePushFailed()),
            /*******待使用********/

            /*******推送中********/
            new Transition(FirstPurchaseGiftCertificateStateConfig::PUSHING, FirstPurchaseGiftCertificateStateConfig::EVENT_OPEN_NOTICE_SUCCESS, FirstPurchaseGiftCertificateStateConfig::PUSH_SUCCESS, new NoticePushSuccess()),
            new Transition(FirstPurchaseGiftCertificateStateConfig::PUSHING, FirstPurchaseGiftCertificateStateConfig::EVENT_OPEN_NOTICE_FAILED, FirstPurchaseGiftCertificateStateConfig::PUSH_FAILED, new NoticePushFailed()),
            /*******推送中********/

            /*******推送失败********/
            new Transition(FirstPurchaseGiftCertificateStateConfig::PUSH_FAILED, FirstPurchaseGiftCertificateStateConfig::EVENT_MANUAL_FAILED_TO_SUCCESS, FirstPurchaseGiftCertificateStateConfig::PUSH_SUCCESS, new ManualFailedToSuccess()),
            /*******推送失败********/
        );

        $this->createMachine($allowState, FirstPurchaseGiftCertificateStateConfig::START, [FirstPurchaseGiftCertificateStateConfig::INVALID, FirstPurchaseGiftCertificateStateConfig::PUSH_SUCCESS], $events, $transitions);

    }

    protected function getStateObject($data, $operator)
    {
        // 存在，不允许当前用户操作
        $redisKey       = redisKey('PREVENT_OPERATE_SAME_GIFT_CERTIFICATE', ['id' => $data['id']]);
        $this->redisKey = $redisKey;

        if ($this->redisKey) {
            if (getOrSetData($this->redisKey)) {
                throw new Error(codeNum('certificate_operate_by_other'), codeMsg('certificate_operate_by_other'), 210);
            }else{
                $this->addLock();
            }
        }
        $giftCertificate = app(FirstPurchaseGiftService::class)->getFirstPurchaseGiftCertificate($data['id']);
        if (empty($giftCertificate)) {
            throw new Error(codeNum('certificate_not_exists'), codeMsg('certificate_not_exists'), 210);
        }
        if ($operator == self::USER_OPERATOR) {
            if ($giftCertificate->user_id != $data['use_id']) {
                throw new Error(codeNum('certificate_not_belong_to_you'), codeMsg('certificate_not_belong_to_you'), 210);
            }
        }
        return $giftCertificate;
    }

    protected function addLock()
    {
        // 加锁，不允许其他用户操作
        getOrSetData($this->redisKey, 1,10);
    }

    protected function delLock()
    {
        // 解锁
        Redis::del($this->redisKey);
    }

    protected function onEnter($event, $data)
    {
        // 参数验证
        $this->validateParams($event, $data);
    }

    protected function onExit()
    {
        // TODO: Implement onExit() method.
    }


    /**
     * 参数验证
     *
     * @param $event
     * @param $post
     *
     * @return mixed
     */
    private function validateParams($event, $post)
    {
        if (!array_key_exists($event, FirstPurchaseGiftCertificateStateConfig::SPECIAL_EVENTS_RULES)) {
            $rules    = FirstPurchaseGiftCertificateStateConfig::COMMON_RULES;
            $messages = FirstPurchaseGiftCertificateStateConfig::COMMON_MESSAGES;
        } else {
            $rules = FirstPurchaseGiftCertificateStateConfig::SPECIAL_EVENTS_RULES[$event]['rules'];
            if (!empty($rules)) {
                if ($event != FirstPurchaseGiftCertificateStateConfig::EVENT_CREATE_GIFT_CERTIFICATE) {
                    $rules    = array_merge(FirstPurchaseGiftCertificateStateConfig::COMMON_RULES, $rules);
                    $messages = array_merge(FirstPurchaseGiftCertificateStateConfig::COMMON_MESSAGES, FirstPurchaseGiftCertificateStateConfig::SPECIAL_EVENTS_RULES[$event]['messages']);
                } else {
                    $messages = FirstPurchaseGiftCertificateStateConfig::SPECIAL_EVENTS_RULES[$event]['messages'];
                }
            } else {
                $rules    = FirstPurchaseGiftCertificateStateConfig::COMMON_RULES;
                $messages = FirstPurchaseGiftCertificateStateConfig::COMMON_MESSAGES;
            }
        }
        $validator = Validator::make($post, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(codeNum('system_error'), $validator->errors()->first(), 210);
        }

        return $post;
    }
}