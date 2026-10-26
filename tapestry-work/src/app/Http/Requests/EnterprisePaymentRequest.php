<?php

namespace App\Http\Requests;

use App\Models\EnterprisePayment;

class EnterprisePaymentRequest
{
    /**
     * 添加企业微信红包
     *
     * @param array $params
     */
    public function add($params = [])
    {
        $rules = [
            'enterprise_id'             => 'required|integer|min:1',
            'name'                      => 'required',
            'type'                      => 'required|in:' . implode(',', [EnterprisePayment::FIXED_AMOUNT_TYPE, EnterprisePayment::RANDOM_AMOUNT_TYPE]),
            'min_amount'                => ['required_if:type,' . EnterprisePayment::RANDOM_AMOUNT_TYPE, 'required'],
            'max_amount'                => ['required_if:type,' . EnterprisePayment::RANDOM_AMOUNT_TYPE, 'required'],
            'fixed_amount'              => ['required_if:type,' . EnterprisePayment::FIXED_AMOUNT_TYPE, 'required'],
            'quantity'                  => 'required',
            'enterprise_payment_amount' => 'required|integer|min:10',
            'is_old_user_receive'       => 'required|in:1,2',
            'qr_code_url'               => 'required',
            'send_name'                 => 'required',
            'wishing'                   => 'required',
            'remark'                    => 'required',
        ];

        $message = [
            'enterprise_id.required'             => '企业ID不存在',
            'enterprise_id.integer'              => '企业ID数据类型有误',
            'enterprise_id.min'                  => '企业ID数据类型有误',
            'name.required'                      => '活动名称不能为空',
            'type.required'                      => '红包类型不能为空',
            'type.in'                            => '红包类型格式有误',
            'min_amount.required'                => '最小区间不能为空',
            'max_amount.required'                => '最大区间不能为空',
            'fixed_amount.required'              => '固定金额不能为空',
            'quantity.required'                  => '红包个数不能为空',
            'enterprise_payment_amount.required' => '红包总金额不能为空',
            'enterprise_payment_amount.integer'  => '红包总金额必须为整数',
            'enterprise_payment_amount.min'      => '红包总金额不能小于一分钱',
            'is_old_user_receive.required'       => '企业现有客户是否可领取不能为空',
            'is_old_user_receive.in'             => '企业现有客户是否可领取格式错误',
            'qr_code_url.required'               => '员工二维码不能为空',
            'send_name.required'                 => '商户名称不能为空',
            'wishing.required'                   => '红包祝福语不能为空',
            'show_off_text.required'             => '炫耀文案不能为空',
            'remark.required'                    => '红包备注不能为空',
        ];
        validator($params, $rules, $message);
    }

    /**
     * 添加企业微信红包
     *
     * @param array $params
     */
    public function commonId($params = [])
    {
        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'id'            => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',
            'id.required'            => '活动ID不能为空',
        ];
        validator($params, $rules, $message);
    }

    /**
     * 编辑企业微信红包
     *
     * @param array $params
     */
    public function save($params = [])
    {
        $rules = [
            'enterprise_id'             => 'required|integer|min:1',
            'id'                        => 'required',
            'name'                      => 'required',
            'type'                      => 'required|in:' . implode(',', [EnterprisePayment::FIXED_AMOUNT_TYPE, EnterprisePayment::RANDOM_AMOUNT_TYPE]),
            'min_amount'                => ['required_if:type,' . EnterprisePayment::RANDOM_AMOUNT_TYPE, 'required'],
            'max_amount'                => ['required_if:type,' . EnterprisePayment::RANDOM_AMOUNT_TYPE, 'required'],
            'fixed_amount'              => ['required_if:type,' . EnterprisePayment::FIXED_AMOUNT_TYPE, 'required'],
            'quantity'                  => 'required',
            'enterprise_payment_amount' => 'required|integer|min:10',
            'is_old_user_receive'       => 'required|in:1,2',
            'qr_code_url'               => 'required',
            'send_name'                 => 'required',
            'wishing'                   => 'required',
            'remark'                    => 'required',
        ];

        $message = [
            'enterprise_id.required'             => '企业ID不存在',
            'enterprise_id.integer'              => '企业ID数据类型有误',
            'enterprise_id.min'                  => '企业ID数据类型有误',
            'id.required'                        => '活动ID不能为空',
            'name.required'                      => '活动名称不能为空',
            'type.required'                      => '红包类型不能为空',
            'type.in'                            => '红包类型格式有误',
            'min_amount.required'                => '最小区间不能为空',
            'max_amount.required'                => '最大区间不能为空',
            'fixed_amount.required'              => '固定金额不能为空',
            'quantity.required'                  => '红包个数不能为空',
            'enterprise_payment_amount.required' => '红包总金额不能为空',
            'enterprise_payment_amount.integer'  => '红包总金额必须为整数',
            'enterprise_payment_amount.min'      => '红包总金额不能小于一分钱',
            'is_old_user_receive.required'       => '企业现有客户是否可领取不能为空',
            'is_old_user_receive.in'             => '企业现有客户是否可领取格式错误',
            'qr_code_url.required'               => '员工二维码不能为空',
            'send_name.required'                 => '商户名称不能为空',
            'wishing.required'                   => '红包祝福语不能为空',
            'show_off_text.required'             => '炫耀文案不能为空',
            'remark.required'                    => '红包备注不能为空',
        ];
        validator($params, $rules, $message);
    }

    /**
     * 状态更改
     *
     * @param array $params
     */
    public function status($params = [])
    {
        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'id'            => 'required',
            'status'        => 'required|in:2,3,4,5'
        ];

        $message = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',
            'id.required'            => '活动ID不能为空',
            'status.required'        => '活动状态不能为空',
            'status.in'              => '活动状态类型不正确',
        ];
        validator($params, $rules, $message);
    }
}