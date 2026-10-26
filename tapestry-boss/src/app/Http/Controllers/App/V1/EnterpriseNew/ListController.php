<?php


namespace App\Http\Controllers\App\V1\EnterpriseNew;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\Industry;
use App\Services\EnterpriseConfigService;
use App\Services\EnterpriseService;
use Illuminate\Http\Request;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Support\Facades\Redis;

class ListController extends BaseController
{
    public function index(Request $request, EnterpriseService $enterpriseService)
    {
        // 查询参数处理
        $params  = $request->input();
        $page    = $params['page'] ?? Constant::DEFAULT_PAGE;
        $perPage = $params['limit'] ?? Constant::DEFAULT_LIMIT;

        // 列表数据查询
        $enterprise      = $enterpriseService->enterpriseList($params, $page, $perPage);
        $enterpriseIds   = array_column($enterprise['data'], 'id');
        $industryIds     = array_column($enterprise['data'], 'industry_id');
        $warehouseConfig = EnterpriseConfigService::WarehouseList($enterpriseIds);

        $redis = Redis::connection();
        if ($enterpriseIds) {
            $names = (new Industry())->idKeyName($industryIds);
            $list  = (new EnterprisePayment())->budgetByEnterprise(['enterprise_ids' => $enterpriseIds]);
            $list  = array_column($list['data'] ?? [], null, 'enterprise_id');

            foreach ($enterprise['data'] as &$item) {

                $item['available_system'] = json_decode($item['available_system'], true);

                if (isset($warehouseConfig[$item['id']]) === false) {
                    return error(201113);
                }
                $item['warehouse_type'] = $warehouseConfig[$item['id']];
                $item['balance']        = 0;
                if (isset($list[$item['id']])) {
                    $item['balance'] = $list[$item['id']]['balance'] ?? 0;
                }
                $item['industry'] = $names[$item['industry_id']] ?? '';
                // 已过期
                if ($item['account_expire_at'] < time()) {
                    $item['expire_type'] = 1;
                } else if ($item['account_expire_at'] > time() && $item['account_expire_at'] < time() + 2592000) {
                    $item['expire_type'] = 3;
                } else {
                    $item['expire_type'] = 2;
                }

                // 是否有更新付费
                $redisKey = redisKey('REDIS_ENTERPRISE_PAYMENT_HAS_NEW', ['enterprise_id' => $item['id']]);
                $item['enterprise_payment_status'] = $redis->exists($redisKey);
            }
            unset($item);
        }
        return success([
            'data' => $enterprise['data'],
            'meta' => [
                'total' => $enterprise['total'],
                'size'  => (int)$enterprise['per_page'],
                'page'  => $enterprise['current_page']
            ]
        ]);
    }
}
