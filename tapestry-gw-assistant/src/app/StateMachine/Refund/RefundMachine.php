<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/13
 * Time: 16:51
 */

namespace App\StateMachine\Refund;

use App\Exceptions\Error;
use App\StateMachine\Refund\Actions\SendOutRefundApply;
use App\StateMachine\StateMachine;
use App\StateMachine\Refund\Actions\MoveRefundApply;
use App\StateMachine\Refund\Actions\CreateRefundApply;
use Illuminate\Support\Facades\Validator;

class RefundMachine extends StateMachine
{
    /**
     * 防并发，redis加锁控制
     *
     * @var
     */
    protected $redisKey;

    protected $event;

    protected $post;

    protected $refundApplyUpdate = [];

    protected $request;

    protected $createRefundApply;

    protected $moveRefundApply;

    /**
     * 实例化状态机
     */
    public function __construct()
    {
        $refundMachineConfig = new RefundMachineConfig();

        //可以有的状态
        $allowState = $refundMachineConfig::ALLOW_STATES;

        $events = RefundMachineConfig::ALLOW_EVENTS;

        $transitions = $refundMachineConfig->getTransitions($events);

        $initialState = RefundMachineConfig::STATUS_INIT;

        $finalStates = RefundMachineConfig::FINAL_STATES;

        $this->createMachine($allowState, $initialState, $finalStates, $events, $transitions);
    }

    public function run($data, $operator = 0)
    {
        if (empty($data['event'])) {
            throw new Error(codeNum('event_required'), codeMsg('event_required'), 210);
        }

        $event = $data['event'];

        $this->event = $event;

        // 参数验证
        $this->validateParams($data);

        unset($data['event']);

        if (in_array($event, RefundMachineConfig::EVENT_APPLY_REFUND_OR_RETURN)) {
            // 申请售后
            (new CreateRefundApply())->run($event, $data);
        } elseif(in_array($event,RefundMachineConfig::EVENT_SEND_OUT_REFUND_APPLY)) {
            // 售后流转
            (new SendOutRefundApply())->run($event, $data);
        }else{
            // 售后流转
            (new MoveRefundApply())->run($event, $data);
        }


        return true;
    }

    /**
     * 验证事件合法
     */
    protected function validateEvent()
    {
        $event = $this->event;

        if (!in_array($event, array_keys(RefundMachineConfig::ALLOW_EVENTS))) {
            throw new Error(codeNum('refund_apply_event_not_exists'), codeMsg('refund_apply_event_not_exists'), 210);
        }
    }

    /**
     * 表单提交参数验证
     */
    protected function validateParams($post)
    {
        $event = $this->event;

        // 事件验证
        $this->validateEvent();

        $rules = RefundMachineConfig::ALLOW_EVENTS[$event]['rules'];

        $messages = RefundMachineConfig::ALLOW_EVENTS[$event]['messages'];

        $validator = Validator::make($post, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(100000, $validator->errors()->first(), 210);
        }
    }

    protected function getStateObject($data, $operator)
    {

    }

    protected function addLock()
    {

    }

    protected function delLock()
    {

    }

    protected function onEnter($event, $data)
    {
        // TODO: Implement onEnter() method.
    }

    protected function onExit()
    {
        // TODO: Implement onExit() method.
    }
}
