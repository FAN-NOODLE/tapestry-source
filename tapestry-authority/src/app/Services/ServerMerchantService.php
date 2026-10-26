<?php

namespace App\Services;

use App\Models\ServiceMerchantSetting;

class ServerMerchantService extends BaseService
{
    /**
     * 获取商户详情
     * @param $id
     * @param int $status
     * @return array
     */
    public function merchantDetail($id, $status = 0)
    {
        $query = ServiceMerchantSetting::query()
            ->where('id', '=', $id);

        if (in_array($status, [ServiceMerchantSetting::ENABLE, ServiceMerchantSetting::DISABLE])) {
            $setting = $query->where('status', '=', $status)->first();
        } else {
            $setting = $query->first();
        }

        if (empty($setting)) {
            return [];
        }

        return $setting->toArray();
    }

    public function merchantList()
    {
        $list = ServiceMerchantSetting::query()
            ->where('status', ServiceMerchantSetting::ENABLE)
            ->get(['id', 'name', 'is_default'])
            ->toArray();
        if (empty($list)) {
            return [];
        }
        return $list;
    }

    public function defaultMerchant()
    {
        $default = ServiceMerchantSetting::query()
            ->where('status', ServiceMerchantSetting::ENABLE)
            ->where('is_default', ServiceMerchantSetting::DEFAULT_YES)
            ->first()
            ->toArray();
        if (empty($default)) {
            return [];
        }
        return $default;
    }

    public function getMerchantByMerchantId($merchantId)
    {
        $detail = ServiceMerchantSetting::query()
            ->where('status', ServiceMerchantSetting::ENABLE)
            ->where('merchant_id', $merchantId)
            ->first()
            ->toArray();
        if (empty($detail)) {
            return [];
        }
        return $detail;
    }
}