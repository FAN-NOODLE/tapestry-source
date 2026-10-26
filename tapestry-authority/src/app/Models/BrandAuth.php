<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\DB;
use App\Services\EnterpriseAdminService;

class BrandAuth extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'brand_auth';

    protected $guarded = [];

    /**
     * 获取品牌权限相关配置
     * @param int $enterpriseId
     * @param string $business
     * @param array $params
     * @param array $fields
     * @return array
     * @Date: 2021/6/30
     * @Time: 16:06
     */
    public function getBrandAuth(int $enterpriseId, string $business = '', array $params = [], array $fields = ['*']): array
    {
        return self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->when(!empty($business), function ($query) use ($business) {
                return $query->where('business', $business);
            })
            ->when(!empty($params), function ($query) use ($params) {
                if (isset($params['method']) && !empty($params['method'])) {
                    return $query->where('method', $params['method']);
                }

                if (isset($params['type']) && !empty($params['type'])) {
                    return $query->where('type', $params['type']);
                }
            })
            ->get()
            ->toArray();
    }
}
