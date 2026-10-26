<?php


namespace App\Http\Controllers\App\V1\EnterprisePayment;


use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterprisePaymentService\EnterprisePaymentService;
use Illuminate\Support\Facades\Redis;

class ListController extends BaseController
{
    public function index(EnterprisePaymentService $enterprisePaymentService)
    {

        $this->validate($this->request, [
            'enterprise_id' => 'required | integer'
        ],[
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
        ]);

        $params = $this->request->only([
            'payment_status',
            'enterprise_id',
            'start_time',
            'end_time',
            'time_type',
            'page',
            'limit'
        ]);

        $page = $params['page'] ?? 1;
        $perPage = $params['limit'] ?? 10;

        $list = $enterprisePaymentService->getList($params, $page, $perPage);
        $redisKey = redisKey('REDIS_ENTERPRISE_PAYMENT_HAS_NEW', ['enterprise_id' => $params['enterprise_id']]);
        $redis = Redis::connection();
        if ($redis->exists($redisKey)) {
            $redis->del($redisKey);
        }
        return success([
            'data' => $list['data'],
            'mate' => [
                'total' => $list['total'],
                'size' => (int) $list['per_page'],
                'page' => $list['current_page']
            ],
        ]);
    }

    public function renewList(EnterprisePaymentService $enterprisePaymentService)
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required | integer'
        ],[
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
        ]);

        $params = $this->request->only([
            'enterprise_id',
            'page',
            'limit'
        ]);

        $page = $params['page'] ?? 1;
        $perPage = $params['limit'] ?? 10;

        // 续费列表类目类型
        $params['renew'] = 1;
        $list = $enterprisePaymentService->getList($params, $page, $perPage);
        return [
            'data' => $list['data'],
            'mate' => [
                'total' => $list['total'],
                'size' => (int) $list['per_page'],
                'page' => $list['current_page']
            ],
        ];
    }
}