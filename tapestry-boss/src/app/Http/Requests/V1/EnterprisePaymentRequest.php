<?php

namespace App\Http\Requests\V1;

use App\Models\Enterprise;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;
use App\Models\EnterprisePayment;

class EnterprisePaymentRequest extends BaseValidateRequest
{
    /**
     * 企业付款添加验证
     *
     * @param array $data
     * @throws ValidationException
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public static function add($data = [])
    {

        $rule = [
            'enterprise_id'    => 'required | integer',
            'payment_category' => ['required', Rule::in([EnterprisePayment::PAYMENT_CATEGORY_RED_PACK, EnterprisePayment::PAYMENT_CATEGORY_SMS, EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])],
            'payment_amount'   => 'required | integer | between:0,9999999999',
            'payment_time'     => 'required | integer',
            'payment_status'   => ['required', Rule::in([EnterprisePayment::PAYMENT_STATUS_NO, EnterprisePayment::PAYMENT_STATUS_YES])],
            'payment_type'     => ['required', Rule::in([EnterprisePayment::PAYMENT_TYPE_OTHER, EnterprisePayment::PAYMENT_TYPE_CASH, EnterprisePayment::PAYMENT_TYPE_WX, EnterprisePayment::PAYMENT_TYPE_ZFB, EnterprisePayment::PAYMENT_TYPE_BANK, EnterprisePayment::PAYMENT_TYPE_TRANSFER])],
            'payment_order_sn'  => 'string | between:0,255',
            'voucher_url'       => 'string | between:0,255',
            'remark'            => 'string | between:0,100',
        ];

        $message = [
            'enterprise_id.required'    => '企业ID不能为空',
            'enterprise_id.integer'     => '企业ID类型错误',
            'payment_category.required' => '购买类目不能为空',
            'payment_amount.required'   => '支付金额不能为空',
            'payment_amount.between'    => '支付金额不能小于0且不能大于99999999.99',
            'payment_amount.integer'    => '支付金额类型错误',
            'payment_time.required'     => '支付时间不能为空',
            'payment_time.integer'      => '支付时间类型错误',
            'payment_status.required'   => '支付状态不能为空',
            'payment_status.in'         => '支付状态类型错误',
            'payment_type.required'     => '支付方式不能为空',
            'payment_type.in'           => '支付方式类型错误',
            'payment_order_sn.string'   => '支付流水号类型错误',
            'payment_order_sn.between'  => '流水号长度在0-255个字符',
            'voucher_url.string'        => '凭证图片类型错误',
            'voucher_url.between'       => '凭证图片连接长度错误',
            'remark.string'             => '备注类型错误',
            'remark.between'            => '备注长度在0-100个字符',
        ];


        // 0：红包增加相关字段验证
        if ($data['payment_category'] == EnterprisePayment::PAYMENT_CATEGORY_RED_PACK) {
            $rule = array_merge($rule, [
                'available_amount' => 'required | integer | between:0,9999999999',
            ]);

            $message = array_merge($message, [
                'available_amount.required' => '可用金额不能为空',
                'available_amount.integer'  => '可用金额类型错误',
                'available_amount.between'  => '可用金额不能小于0且不能大于99999999.99',
            ]);
        }
        // 1：短信增加相关字段验证

        if ($data['payment_category'] == EnterprisePayment::PAYMENT_CATEGORY_SMS) {
            $rule = array_merge($rule, [
                'buy_number' => 'required | integer | between:0,99999999',
            ]);

            $message = array_merge($message, [
                'buy_number.required' => '购买数量不能为空',
                'buy_number.integer'  => '购买数量类型错误',
                'buy_number.between'  => '购买数量不能小于0且不能大于99999999',
            ]);

        }

        // 2：智慧导购，3：SCRM，4：商城，5：综合系统 增加相关字段验证
        if (in_array($data['payment_category'], [EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])) {
            $rule = array_merge($rule, [
                'collection_amount' => 'integer | between:0,9999999999',
                'server_charge'     => 'integer | between:0,9999999999',
                'buy_number'        => 'required | integer | between:0,9999999999',
                'buy_day'           => 'required | integer | between:0,9999999999',
                'use_time'          => 'required | integer',
            ]);

            $message = array_merge($message, [
                'collection_amount.integer' => '可用金额类型错误',
                'collection_amount.between' => '可用金额不能小于0且不能大于99999999.99',
                'server_charge.integer'     => '服务费类型错误',
                'server_charge.between'     => '服务费不能小于0且不能大于99999999.99',
                'buy_number.required'       => '购买账号数量不能为空',
                'buy_number.integer'        => '购买账号数量类型错误',
                'buy_number.between'        => '购买账号数量不能小于0且不能大于99999999.99',
                'buy_day.required'          => '购买时长不能为空',
                'buy_day.integer'           => '购买时长类型错误',
                'buy_day.between'           => '购买时长不能小于1天且不能大于99999999天',
                'use_time.required'         => '生效时间不能为空',
                'use_time.integer'          => '生效时间类型错误',
            ]);
        }
        self::validator($data, $rule, $message);
    }

    /**
     * 企业付款详情验证
     * @param $data
     * @throws ValidationException
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public static function detail($data = [])
    {
        $rule    = [
            'id' => 'required | integer',
        ];
        $message = [
            'id.required' => '付费id不能为空',
            'id.integer'  => '付费id类型错误',
        ];
        self::validator($data, $rule, $message);

    }

    /**
     * 企业付款确认收款验证
     * @param array $data
     * @throws ValidationException
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public static function receivePayment($data = [])
    {
        $rule = [
            'id'                => 'required | integer',
            'collection_amount' => 'required | integer | between:0,9999999999',
            'server_charge'     => 'required | integer | between:0,9999999999'
        ];

        $message = [
            'id.required'                => '付费id不能为空',
            'id.integer'                 => '付费id类型错误',
            'collection_amount.required' => '收款金额不能为空',
            'collection_amount.integer'  => '收款金额类型错误',
            'collection_amount.between'  => '收款金额不能小于0且不能大于99999999.99',
            'server_charge.required'     => '服务费不能为空',
            'server_charge.integer'      => '服务费类型错误',
            'server_charge.between'      => '服务费不能小于0且不能大于99999999.99',
        ];
        self::validator($data, $rule, $message);
    }

    /**
     * 企业付款续费验证
     * @param array $data
     * @throws ValidationException
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public static function renew($data = [])
    {
        $rule = [
            'id'                => 'required | integer',
            'payment_time'      => 'required | integer',
            'payment_amount'    => 'required | integer | between:0,9999999999',
            'buy_day'           => 'required | integer | between:0,99999999',
            'payment_status'    => ['required', Rule::in([EnterprisePayment::PAYMENT_STATUS_NO, EnterprisePayment::PAYMENT_STATUS_YES])],
            'payment_type'      => ['required', Rule::in([EnterprisePayment::PAYMENT_TYPE_OTHER, EnterprisePayment::PAYMENT_TYPE_CASH, EnterprisePayment::PAYMENT_TYPE_WX, EnterprisePayment::PAYMENT_TYPE_ZFB, EnterprisePayment::PAYMENT_TYPE_BANK, EnterprisePayment::PAYMENT_TYPE_TRANSFER])],
            'collection_amount' => 'integer | between:0,9999999999',
            'server_charge'     => 'integer | between:0,9999999999',
            'payment_order_sn'  => 'string | between:0,255',
            'voucher_url'       => 'string | between:0,255',
            'remark'            => 'string | between:0,100',

        ];

        $message = [
            'id.required'               => '付费id不能为空',
            'id.integer'                => '付费id类型错误',
            'payment_amount.required'   => '支付金额不能为空',
            'payment_amount.between'    => '支付金额不能小于0且不能大于99999999.99',
            'payment_amount.integer'    => '支付金额类型错误',
            'buy_day.required'          => '购买时长不能为空',
            'buy_day.integer'           => '购买时长类型错误',
            'buy_day.between'           => '购买时长不能小于0天且不能大于99999999天',
            'payment_time.required'     => '支付时间不能为空',
            'payment_time.integer'      => '支付时间类型错误',
            'payment_status.required'   => '支付状态不能为空',
            'payment_status.in'         => '支付状态类型错误',
            'payment_type.required'     => '支付方式不能为空',
            'payment_type.in'           => '支付方式类型错误',
            'collection_amount.integer' => '收款金额类型错误',
            'collection_amount.between' => '收款金额不能小于0且不能大于99999999.99',
            'server_charge.integer'     => '服务费类型错误',
            'server_charge.between'     => '服务费不能小于0且不能大于99999999.99',
            'payment_order_sn.string'   => '支付流水号类型错误',
            'payment_order_sn.between'  => '流水号长度在0-255个字符',
            'voucher_url.string'        => '凭证图片类型错误',
            'voucher_url.between'       => '凭证图片连接长度错误',
            'remark.string'             => '备注类型错误',
            'remark.between'            => '备注长度在0-100个字符',
        ];
        self::validator($data, $rule, $message);
    }
}
