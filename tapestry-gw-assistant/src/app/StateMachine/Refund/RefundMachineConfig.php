<?php
/**
 * 售后状态机状态配置
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/17
 * Time: 11:13
 */

namespace App\StateMachine\Refund;


use App\StateMachine\Transition;

class RefundMachineConfig
{
    /**
     * 初始状态
     */
    const STATUS_INIT = 0;

    /**
     * 待审核（未发货，申请售后）
     */
    const STATUS_CUSTOMER_APPLY = 21;

    /**
     * 后台通过，待OPEN审核
     */
    const STATUS_CONSOLE_AGREE = 22;

    /**
     * OPEN审核通过，待用户寄回商品
     */
    const STATUS_OPEN_AGREE_RETURN = 23;

    /**
     * 用户已寄出商品，待OPEN确认收货、验货
     */
    const STATUS_CUSTOMER_DELIVERED = 24;

    /**
     * 人工中
     */
    const STATUS_MANUAL_DEALING = 25;

    /**
     * OPEN 同意退款
     */
    const STATUS_AGREE_REFUND = 26;

    /**
     * 已发货
     */
    const STATUS_SEND_OUT = 27;

    /**
     * 已取消（用户取消/过期未寄回）
     */
    const STATUS_CANCELED = 28;

    /**
     * 已拒绝
     */
    const STATUS_REJECTED = 29;

    /**
     * 事件1：用户申请售后（未发货，仅退款）
     */
    const EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND = 'customer_apply_refund_of_not_send';

    /**
     * 事件2：用户申请售后（已发货，仅退款）
     */
    const EVENT_CUSTOMER_APPLY_REFUND_OF_HAS_SENT = 'customer_apply_refund_of_has_sent';

    /**
     * 事件3：用户申请售后（已发货，退货退款）
     */
    const EVENT_CUSTOMER_APPLY_RETURN_OF_HAS_SENT = 'customer_apply_return_of_has_sent';

    /**
     * 事件4：用户申请售后（已签收，仅退款）
     */
    const EVENT_CUSTOMER_APPLY_REFUND_OF_HAS_SIGNED = 'customer_apply_refund_of_has_signed';

    /**
     * 事件5：用户申请售后（已签收，退货退款）
     */
    const EVENT_CUSTOMER_APPLY_RETURN_OF_HAS_SIGNED = 'customer_apply_return_of_has_signed';

    /**
     * 事件6：用户再次申请售后（部分退款，退货退款）
     */
    const EVENT_CUSTOMER_APPLY_RETURN_AGAIN = 'customer_apply_return_again';

    /**
     * 事件7：用户取消
     */
    const EVENT_CUSTOMER_CANCEL = 'customer_cancel';

    /**
     * 事件8：过期未寄回商品
     */
    const EVENT_AUTO_OVERDUE_DELIVER = 'customer_overdue_deliver';

    /**
     * 事件9：后台审核通过
     */
    const EVENT_CONSOLE_AGREE = 'console_agree';

    /**
     * 事件10：后台驳回
     */
    const EVENT_CONSOLE_REJECT = 'console_reject';

    /**
     * 事件11：OPEN同意退货
     */
    const EVENT_OPEN_AGREE_RETURN = 'open_agree_return';

    /**
     * 事件12：OPEN拒绝退货
     */
    const EVENT_OPEN_REJECT_RETURN = 'open_reject_return';

    /**
     * 事件13：用户寄回商品
     */
    const EVENT_CUSTOMER_DELIVER = 'customer_deliver';

    /**
     * 事件14：OPEN同意退款
     */
    const EVENT_OPEN_AGREE_REFUND = 'open_agree_refund';

    /**
     * 事件15：OPEN拒绝退款
     */
    const EVENT_OPEN_REJECT_REFUND = 'open_reject_refund';

    /**
     * 事件16：人工同意退款
     */
    const EVENT_MANUAL_AGREE_REFUND = 'manual_agree_refund';

    /**
     * 事件17：人工拒绝退款
     */
    const EVENT_MANUAL_REJECT_REFUND = 'manual_reject_refund';

    /**
     * 事件18：OPEN 通知转人工
     */
    const EVENT_OPEN_NOTICE_TO_MANUAL = 'open_notice_to_manual';

    /**
     * 事件19：OPEN 通知已发货
     */
    const EVENT_OPEN_NOTICE_SEND_OUT = 'open_notice_send_out';

    /**
     * 事件20：OPEN 通知仅退款
     */
    const EVENT_OPEN_NOTICE_REFUND = 'open_notice_refund';

    /**
     * 事件21：人工处理为已发货
     */
    const EVENT_MANUAL_SEND_OUT = 'manual_send_out';

    /**
     * 事件22：open通知人工中的售后单为已发货
     */
    const EVENT_OPEN_NOTICE_MANUAL_SEND_OUT = 'open_notice_manual_send_out';

    /**
     * 事件23：用户申请售后（未发货，仅退款）,特殊企业处理
     */
    const EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND_SPECIAL = 'customer_apply_refund_of_not_send_special';

    /**
     * 事件24：已发货用户申请仅退款
     */
    const EVENT_SEND_OUT_APPLY_ONLY_REFUND = 'send_out_apply_only_refund';

    /**
     * 事件25：已发货用户申请退货退款
     */
    const EVENT_SEND_OUT_APPLY_RETURN_GOODS = 'send_out_apply_return_goods';

    /**
     * 事件26：客服一键转人工
     */
    const EVENT_CUSTOMER_SERVICE_TO_MANUAL = 'customer_service_to_manual';

    /**
     * 事件27：客服创建手工售后单
     */
    const EVENT_MANUAL_APPLY_RETURN_GOODS = 'manual_apply_return_goods';

    /**
     * 申请售后（仅退款或退货退款）
     */
    const EVENT_APPLY_REFUND_OR_RETURN = [
        self::EVENT_CUSTOMER_APPLY_REFUND_OF_HAS_SIGNED,
        self::EVENT_CUSTOMER_APPLY_REFUND_OF_HAS_SENT,
        self::EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND,
        self::EVENT_CUSTOMER_APPLY_RETURN_AGAIN,
        self::EVENT_CUSTOMER_APPLY_RETURN_OF_HAS_SIGNED,
        self::EVENT_CUSTOMER_APPLY_RETURN_OF_HAS_SENT,
        self::EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND_SPECIAL,
    ];

    /**
     * 已发货下的售后
     */
    const EVENT_SEND_OUT_REFUND_APPLY = [
        self::EVENT_SEND_OUT_APPLY_ONLY_REFUND,
        self::EVENT_SEND_OUT_APPLY_RETURN_GOODS,
        self::EVENT_MANUAL_APPLY_RETURN_GOODS,
    ];

    /**
     *  所有状态
     */
    const ALLOW_STATES = [
        self::STATUS_CUSTOMER_APPLY     => '待审核（未发货，申请售后）',
        self::STATUS_CONSOLE_AGREE      => '后台通过，待OPEN审核',
        self::STATUS_OPEN_AGREE_RETURN  => 'OPEN审核通过，待用户寄回商品',
        self::STATUS_CUSTOMER_DELIVERED => '用户已寄出商品，待OPEN确认收货、验货',
        self::STATUS_MANUAL_DEALING     => '人工中',
        self::STATUS_AGREE_REFUND       => 'OPEN 同意退款',
        self::STATUS_SEND_OUT           => '已发货',
        self::STATUS_CANCELED           => '已取消（用户取消/过期未寄回）',
        self::STATUS_REJECTED           => '已拒绝',
    ];

    /**
     *  售后有限状态机终态
     *      同意退款
     *      已发货
     *      已取消
     *      已拒绝
     */
    const FINAL_STATES = [
        self::STATUS_AGREE_REFUND,
        self::STATUS_SEND_OUT,
        self::STATUS_CANCELED,
        self::STATUS_REJECTED
    ];

    const ALLOW_EVENTS = [
        self::EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND => [
            'desc'        => '未发货，申请仅退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CONSOLE_AGREE,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND_SPECIAL => [
            'desc'        => '未发货，申请仅退款，特殊处理',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_MANUAL_DEALING,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_APPLY_REFUND_OF_HAS_SENT => [
            'desc'        => '已发货，申请仅退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_APPLY_REFUND_OF_HAS_SIGNED => [
            'desc'        => '已签收，申请仅退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_APPLY_RETURN_OF_HAS_SENT => [
            'desc'        => '已发货，申请退货退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_APPLY_RETURN_OF_HAS_SIGNED => [
            'desc'        => '已签收，申请退货退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_APPLY_RETURN_AGAIN => [
            'desc'        => '部分退款，申请退货退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],

        self::EVENT_CUSTOMER_CANCEL => [
            'desc'        => '用户取消售后',
            'allowStatus' => [
                self::STATUS_CUSTOMER_APPLY,
                self::STATUS_CONSOLE_AGREE,
                self::STATUS_OPEN_AGREE_RETURN,
            ],
            'toStatus'    => self::STATUS_CANCELED,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_AUTO_OVERDUE_DELIVER => [
            'desc'        => '超期未寄回商品',
            'allowStatus' => self::STATUS_OPEN_AGREE_RETURN,
            'toStatus'    => self::STATUS_CANCELED,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_CONSOLE_AGREE => [
            'desc'        => '后台审核通过',
            'allowStatus' => self::STATUS_CUSTOMER_APPLY,
            'toStatus'    => self::STATUS_CONSOLE_AGREE,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_CONSOLE_REJECT => [
            'desc'        => '后台驳回申请',
            'allowStatus' => self::STATUS_CUSTOMER_APPLY,
            'toStatus'    => self::STATUS_REJECTED,
            'rules'       => [
                'id'            => 'required',
                'reject_reason' => 'required',
            ],
            'messages'    => [
                'id.required'            => '售后主键不能为空',
                'reject_reason.required' => '请填写拒绝原因',
            ]
        ],

        self::EVENT_OPEN_AGREE_RETURN => [
            'desc'        => 'OPEN同意退货',
            'allowStatus' => self::STATUS_CONSOLE_AGREE,
            'toStatus'    => self::STATUS_OPEN_AGREE_RETURN,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_OPEN_REJECT_RETURN => [
            'desc'        => 'OPEN拒绝退货',
            'allowStatus' => self::STATUS_CONSOLE_AGREE,
            'toStatus'    => self::STATUS_REJECTED,
            'rules'       => [
                'id'            => 'required',
                'reject_reason' => 'required',
            ],
            'messages'    => [
                'id.required'            => '售后主键不能为空',
                'reject_reason.required' => '请填写拒绝原因',
            ]
        ],

        self::EVENT_CUSTOMER_DELIVER => [
            'desc'        => '用户寄回商品',
            'allowStatus' => self::STATUS_OPEN_AGREE_RETURN,
            'toStatus'    => self::STATUS_CUSTOMER_DELIVERED,
            'rules'       => [
                'id'                => 'required',
                'logistic_identity' => 'required',
                'logistics_no'      => 'required',
            ],
            'messages'    => [
                'id.required'                => '售后主键不能为空',
                'logistic_identity.required' => '请填写物流公司标识',
                'logistics_no.required'      => '请填写物流单号',
            ]
        ],

        self::EVENT_OPEN_AGREE_REFUND => [
            'desc'        => 'OPEN同意退款',
            'allowStatus' => self::STATUS_CUSTOMER_DELIVERED,
            'toStatus'    => self::STATUS_AGREE_REFUND,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_OPEN_REJECT_REFUND => [
            'desc'        => 'OPEN拒绝退款',
            'allowStatus' => self::STATUS_CUSTOMER_DELIVERED,
            'toStatus'    => self::STATUS_MANUAL_DEALING,
            'rules'       => [
                'id'            => 'required',
                'reject_reason' => 'required',
            ],
            'messages'    => [
                'id.required'            => '售后主键不能为空',
                'reject_reason.required' => '请填写拒绝原因',
            ]
        ],

        self::EVENT_MANUAL_AGREE_REFUND => [
            'desc'        => '人工同意退款',
            'allowStatus' => self::STATUS_MANUAL_DEALING,
            'toStatus'    => self::STATUS_AGREE_REFUND,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_MANUAL_REJECT_REFUND => [
            'desc'        => '人工拒绝退款',
            'allowStatus' => self::STATUS_MANUAL_DEALING,
            'toStatus'    => self::STATUS_REJECTED,
            'rules'       => [
                'id'            => 'required',
                'reject_reason' => 'required',
            ],
            'messages'    => [
                'id.required'            => '售后主键不能为空',
                'reject_reason.required' => '请填写拒绝原因',
            ]
        ],

        self::EVENT_OPEN_NOTICE_TO_MANUAL => [
            'desc'        => 'OPEN通知转人工',
            'allowStatus' => self::STATUS_CONSOLE_AGREE,
            'toStatus'    => self::STATUS_MANUAL_DEALING,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_OPEN_NOTICE_SEND_OUT => [
            'desc'        => 'OPEN通知已发货',
            'allowStatus' => self::STATUS_CONSOLE_AGREE,
            'toStatus'    => self::STATUS_SEND_OUT,
            'rules'       => [
                'id'                => 'required',
            ],
            'messages'    => [
                'id.required'                => '售后主键不能为空',
            ]
        ],

        self::EVENT_OPEN_NOTICE_REFUND => [
            'desc'        => 'OPEN通知仅退款',
            'allowStatus' => self::STATUS_CONSOLE_AGREE,
            'toStatus'    => self::STATUS_AGREE_REFUND,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],

        self::EVENT_MANUAL_SEND_OUT => [
            'desc'        => '人工处理为已发货',
            'allowStatus' => self::STATUS_MANUAL_DEALING,
            'toStatus'    => self::STATUS_SEND_OUT,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],
        self::EVENT_OPEN_NOTICE_MANUAL_SEND_OUT => [
            'desc'        => 'open通知人工中的售后单为已发货',
            'allowStatus' => self::STATUS_MANUAL_DEALING,
            'toStatus'    => self::STATUS_SEND_OUT,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],
        self::EVENT_SEND_OUT_APPLY_ONLY_REFUND => [
            'desc'        => '已发货申请仅退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],
        self::EVENT_SEND_OUT_APPLY_RETURN_GOODS => [
            'desc'        => '已发货申请退货退款',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_id'         => 'required',
                'refund_reason_id' => 'required',
                'reason'           => 'required',

            ],
            'messages'    => [
                'order_id.required' => '订单主键不能为空',
                'refund_reason_id'  => '申请售后理由id不能为空',
                'reason.required'   => '请选择申请售后的理由',
            ]
        ],
        self::EVENT_CUSTOMER_SERVICE_TO_MANUAL => [
            'desc'        => '客服一键转人工',
            'allowStatus' => [
                self::STATUS_CONSOLE_AGREE,
                self::STATUS_OPEN_AGREE_RETURN,
                self::STATUS_CUSTOMER_DELIVERED,
            ],
            'toStatus'    => self::STATUS_MANUAL_DEALING,
            'rules'       => [
                'id' => 'required',
            ],
            'messages'    => [
                'id.required' => '售后主键不能为空',
            ]
        ],
        self::EVENT_MANUAL_APPLY_RETURN_GOODS => [
            'desc'        => '客服创建手工售后单',
            'allowStatus' => self::STATUS_INIT,
            'toStatus'    => self::STATUS_CUSTOMER_APPLY,
            'rules'       => [
                'order_no'           => 'required',
                'refund_apply_type'  => 'required',
                'order_goods_status' => 'required',
                'user_problem'       => 'required',
                'return_goods_info'  => 'required'
            ],
            'messages'    => [
                'order_no.required'           => '订单编号为空',
                'refund_apply_type.required'  => '售后方式不能为空',
                'order_goods_status.required' => '货物状态不能为空',
                'user_problem.required'       => '用户问题不能为空',
                'return_goods_info.required'  => '售后商品数据不能为空',
            ]
        ],
    ];

    /**
     * 状态流转线路配置（同一事件，状态不同，对应多条线路）
     *
     * @param $events
     *
     * @return array
     */
    public function getTransitions($events)
    {
        $transitions = [];
        foreach ($events as $event => $eventConfig) {
            $allowStatus = $eventConfig['allowStatus'];
            $toStatus    = $eventConfig['toStatus'];

            if (is_array($allowStatus)) {
                foreach ($allowStatus as $allow) {
                    $transitions[] = new Transition($allow, $event, $toStatus);
                }
            } else {
                $transitions[] = new Transition($allowStatus, $event, $toStatus);
            }
        }

        return $transitions;
    }
}
