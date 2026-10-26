<?php
/**
 * Created by PhpStorm.
 * User: wpf
 * Date: 2020/4/12
 * Time: 16:52
 */

namespace App\Services;

use App\Models\Enterprise;
use App\Models\EnterpriseWorkPlugin;
use CrsSdk\Shop\Payment;
use App\Models\EnterpriseThirdPartAccount;

class EnterpriseService extends BaseService
{
    /**
     * 版本 1=标准版，2=门店商城版，3=中心商城版
     */
    const VERSION_SHOP = 1;
    const VERSION_STORES = 2;
    const VERSION_CENTER = 3;

    /**
     * 通过企业ID获取企业详情
     * @param $enterpriseId
     * @return array
     */
    public function enterpriseDetail($enterpriseId)
    {
        $detail = Enterprise::query()
            ->where('id', $enterpriseId)
            ->first();

        return $detail ? $detail->toArray() : [];
    }

    /**
     * 通过企业hash获取企业详情
     * @param $hash
     * @return array
     */
    public function detailByHash($hash)
    {
        $detail = Enterprise::query()
            ->where('hash', '=', $hash)
            ->first();

        return $detail ? $detail->toArray() : [];
    }

    public static function baseConfigCheck($eId)
    {
        $shopAuth = 1;
        $eInfo = Enterprise::getDetail($eId);
        if (empty($eInfo)) {
            notice(701005);
        }
        // 授权相关
        $thirdAccount = EnterpriseThirdPartAccount::existsMiniProgram($eId);
        if (empty($thirdAccount)) {
            $shopAuth = 2;
        }

        // 企微相关
        $qwConfig = 1;
        $plugin = EnterpriseWorkPlugin::identitySearch($eId, 'crs-manu-customer');
        if (empty($plugin)) {
            $qwConfig = 2;
        }
        if (empty($eInfo->external_secret)) {
            $qwConfig = 2;
        }

        // 支付相关
        $paymentConfig = 1;
        $payment = (new Payment())->detail(['enterprise_id' => $eId])['data'];
        if (empty($payment)) {
            $paymentConfig = 2;
        }

        return [
            'shop_auth'      => $shopAuth,
            'qw_config'      => $qwConfig,
            'payment_config' => $paymentConfig,
        ];
    }
}
