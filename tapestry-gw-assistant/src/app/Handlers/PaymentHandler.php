<?php

namespace App\Handlers;

use App\Models\PaymentSetting;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Facades\Log;

class PaymentHandler
{
    /**
     * 商城小程序授权类型为：1
     */
    const AUTH_TYPE = 1;

    /**
     * 获取微信支付配置（服务商模式）
     *
     * @param $enterpriseId
     * @param $notifyUrl
     *
     * @return array
     */
    public static function getConfig($enterpriseId, $notifyUrl = '')
    {
        $setting = PaymentSetting::where([
            'enterprise_id'   => $enterpriseId,
            'payment_channel' => PaymentSetting::PAYMENT_CHANNEL_WECHAT,
            'status'          => PaymentSetting::PAYMENT_STATUS_ENABLE,
        ])->first();

        // 企业授权后的账号信息
        $authAccount = (new Auth())->authDetail([
            'enterprise_id' => $enterpriseId,
            'auth_type'     => self::AUTH_TYPE
        ]);

        if (empty($notifyUrl)) {
            // 获取微信支付回调配置
            $wechatCallbackUrl = config('callback.wechat');

            // 查询企业信息
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];

            // 回调地址拼接企业hash
            $hash      = !empty($enterprise) ? $enterprise['hash'] : 'default';
            $notifyUrl = $wechatCallbackUrl['notify'] . '/' . $hash;

            PaymentHandler::getConfig($enterpriseId, $notifyUrl);
        }

        // 是否是普通支付模式
        $ordinaryPayEnterpriseIds = config('enterprise.ordinary_merchant_pay');
        if (in_array($enterpriseId, $ordinaryPayEnterpriseIds)) {
            // 普通支付模式
            return [
                'miniapp_id'  => $authAccount['data']['app_id'], // 小程序 APPID
                'mch_id'      => $setting->merchant_id,
                'key'         => $setting->merchant_key,
                'serial_no'   => $setting->serial_no,
                'notify_url'  => $notifyUrl, // route函数第3个参数代表是否是https路由
                'cert_client' => storage_path($setting->public_key_path), // optional，退款等情况时用到
                'cert_key'    => storage_path($setting->private_key_path),// optional，退款等情况时用到
            ];
        } else {
            // 服务商模式
            $merchant_id     = config('micro_service.service_mch_id');
            $service_mch_key = config('micro_service.service_mch_key');

            return [
                'miniapp_id'     => config('micro_service.service_appid'),
                'mch_id'         => $merchant_id,
                'key'            => $service_mch_key,
                'notify_url'     => $notifyUrl,
                'cert_client'    => storage_path('cert/apiclient_cert.pem'), // optional，退款等情况时用到
                'cert_key'       => storage_path('cert/apiclient_key.pem'), // optional，退款等情况时用到
                'sub_miniapp_id' => $authAccount['data']['app_id'],
                'sub_mch_id'     => $setting->merchant_id,
                'mode'           => 'service',
            ];
        }
    }

    /**
     * 获取微信支付配置（普通模式）
     *
     * @param $enterpriseId
     *
     * @return array
     */
    public static function oldGetConfig($enterpriseId)
    {
        $setting = PaymentSetting::where([
            'enterprise_id'   => $enterpriseId,
            'payment_channel' => PaymentSetting::PAYMENT_CHANNEL_WECHAT,
            'status'          => PaymentSetting::PAYMENT_STATUS_ENABLE,
        ])->first();

        // 企业授权后的账号信息
        $authAccount = (new Auth())->authDetail([
            'enterprise_id' => $enterpriseId,
            'auth_type'     => self::AUTH_TYPE
        ]);

        return [
            'miniapp_id'  => $authAccount['data']['app_id'], // 小程序 APPID
            'mch_id'      => $setting->merchant_id,
            'key'         => $setting->merchant_key,
            'serial_no'   => $setting->serial_no,
            'notify_url'  => route('payNotify', [], true), // route函数第3个参数代表是否是https路由
            'cert_client' => storage_path($setting->public_key_path), // optional，退款等情况时用到
            'cert_key'    => storage_path($setting->private_key_path),// optional，退款等情况时用到
        ];
    }
}
