<?php
/**
 * 首购礼资格状态机状态配置
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/08/13
 * Time: 18:13
 */

namespace App\StateMachine\FirstPurchaseGiftCertificate;


class FirstPurchaseGiftCertificateStateConfig
{
    const START = 0;
    const WAIT_USE = 1;   //待使用
    const PUSHING = 2;    //推送中
    const PUSH_SUCCESS = 3; //推送成功
    const PUSH_FAILED = 4; //推送失败
    const INVALID = 5;  //已失效

    //操作者
    const CONSOLE_OPERATOR = 'console';
    const USER_OPERATOR = 'user';

    //事件
    const EVENT_CREATE_GIFT_CERTIFICATE = 'create_gift_certificate'; //创建首购礼资格
    const EVENT_PUSH_TO_OPEN = 'push_to_open'; //已推送到open
    const EVENT_AUTO_INVALID = 'auto_invalid'; //自动失效
    const EVENT_MANUAL_INVALID = 'manual_invalid'; //人工处理成失效
    const EVENT_OPEN_NOTICE_SUCCESS = 'push_success';  //open通知推送成功
    const EVENT_OPEN_NOTICE_FAILED = 'push_failed'; //open通知推送失败
    const EVENT_MANUAL_FAILED_TO_SUCCESS = 'manual_failed_to_success'; //人工处理成功
    const EVENT_SECOND_MOBILE_ERROR = 'second_delivery_mobile_not_same'; //第二单收货人手机号与会员手机号不符

    //公共的校验规则
    const COMMON_RULES = [
        'id'    => 'required | integer',
        'event' => 'required',
    ];

    //公共的校验返回信息
    const COMMON_MESSAGES = [
        'id.required'    => 'Sorry,首购礼资格id不能为空',
        'id.integer'     => 'Sorry,首购礼资格id必须为正整数',
        'event.required' => 'Sorry,事件event不能为空',
    ];

    const SPECIAL_EVENTS_RULES = [
        self::EVENT_OPEN_NOTICE_FAILED      => [
            'rules'    => [
                'failed_message' => 'required',
            ],
            'messages' => [
                'failed_message.required' => '请填写失败原因',
            ]
        ],
        self::EVENT_PUSH_TO_OPEN      => [
            'rules'    => [
                'second_order_no' => 'required',
            ],
            'messages' => [
                'second_order_no.required' => '请填写复购订单号',
            ]
        ],
    ];
}