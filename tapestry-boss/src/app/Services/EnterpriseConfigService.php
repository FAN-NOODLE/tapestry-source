<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:20
 */

namespace App\Services;

use App\Models\EnterpriseConfig;
use App\Models\EnterpriseThirdPartAccount;

class EnterpriseConfigService extends BaseService
{
    const WAREHOUSE_KEY = 'warehouse_type';

    /**
     * 查询当前企业是否配置
     * @param int $enterpriseId 企业ID
     *
     * @return bool
     */
    public function isSettingConfig($enterpriseId = 0)
    {
        $total = EnterpriseThirdPartAccount::query()->where('enterprise_id', $enterpriseId)
            ->where('app_id', '<>', '')
            ->where('app_secret', '<>', '')
            ->count();

        if ($total > 0) {
            return true;
        }
        return false;
    }

    /**
     * 获取企业仓库
     * @param $eid
     * @return array
     */
    public static function getWarehouseConfig($eid)
    {
        $eConfig = EnterpriseConfig::query()
            ->where('enterprise_id', '=', $eid)
            ->where('key', '=', self::WAREHOUSE_KEY)
            ->first();

        if ($eConfig) {
            return $eConfig->toArray();
        }

        return [];
    }

    /**
     * 批量企业仓库配置
     * @param $ids
     * @return array
     */
    public static function WarehouseList($ids)
    {
        return EnterpriseConfig::query()
            ->whereIn('enterprise_id', $ids)
            ->where('key', '=', self::WAREHOUSE_KEY)
            ->pluck('value', 'enterprise_id')
            ->toArray();
    }

    /**
     * 设置企业仓库
     * @param $eid
     * @param $type
     * @return EnterpriseConfig|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public static function setWarehouseConfig($eid, $type)
    {
        $eConfig = EnterpriseConfig::query()
            ->where('enterprise_id', '=', $eid)
            ->where('key', '=', self::WAREHOUSE_KEY)
            ->first();

        if ($eConfig) {
            $eConfig->update(['value' => $type]);
            return $eConfig;
        }

        return EnterpriseConfig::query()->create([
            'key'           => self::WAREHOUSE_KEY,
            'type'          => 0,
            'value'         => $type,
            'enterprise_id' => $eid,
            'desc'          => '商品是多仓还是单仓(1:单仓,2:多仓)',
        ]);
    }
}
