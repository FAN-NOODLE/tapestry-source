<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-07-16
 * Time: 09:32
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Traits\EnterpriseTools;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Work\Customer;
use Illuminate\Support\Facades\Log;

class DataCompassController extends BaseController
{
    use EnterpriseTools;

    private $yDataMap
        = [
            '顾客总数',
            '净增顾客数',
            '流失顾客总数',
            '新增顾客数',
            '删除人数',
            '被删除人数'
        ];

    /**
     * @var bool
     */
    private $isGoOn = true;

    /**
     * 顾客数据罗盘统计
     *
     * @param Customer $customer
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Customer $customer)
    {
        $input = $this->checkParams();
        if (!$this->isGoOn) {
            return success($this->defaultResult($input));
        }
        $input['qy_domain']      = $this->loginAdminEnterprise['domain'];
        $input['app_identifier'] = $this->loginAdminEnterprise['identifier'];
        $input['enterprise_id']  = $this->enterpriseId;
        try {
            $data = $this->format($customer->customerCalc($input));
            return success($data);
        } catch (Error $error) {
            Log::error("调用Customer->customerCalc接口出错，错误信息:{$error->getMessage()}", $input);
            notice('search_customer_compass_wrong');
        }
    }

    /**
     * 验证字段
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function checkParams(): array
    {
        $searchCalc = [];
        $rules      = [
            'start_time' => 'integer',
            'end_time'   => 'integer',
            'guide_type' => 'integer',
            'store_ids'  => 'array',
            'store_id'   => 'integer',
            'data_type'  => 'integer',
            'page'       => 'integer',
            'size'       => 'integer'
        ];

        $input = $this->validate($this->request, $rules);
        $searchCalc['page']       = $input['page'] ?? Constant::PAGE;
        $searchCalc['size']       = $input['size'] ?? Constant::LIMIT;
        $input['store_ids']       = $input['store_ids'] ?? [];
        $input['guide_type']      = $input['guide_type'] ?? 0;
        $searchCalc['data_type']  = $input['data_type'];
        $searchCalc['start_time'] = $input['start_time'];
        $searchCalc['end_time']   = strtotime(date('Y-m-d 23:59:59', $input['end_time']));

        if (!empty($input['store_id'])) {
            $input['store_ids'][] = $input['store_id'];
        }
        $enterpriseId     = $this->enterpriseId;
        $searchWorkParams = [
            'is_page'       => false,
            'worker_status' => 1,
            'enterprise_id' => $enterpriseId,
            'store_ids'     => $input['store_ids'],
            'type'          => $input['guide_type']
        ];

        $workers = (new StoreWorker())->getList($searchWorkParams);
        if (empty($workers['list'])) {
            $this->isGoOn = false;
            return $searchCalc;
        }
        $searchCalc['mem_uid'] = array_unique(array_column($workers['list'], 'mem_uid'));
        return $searchCalc;
    }

    /**
     * 数据
     *
     * @param array $input
     * @return array
     */
    private function format(array $input): array
    {
        $chart = [];
        foreach ($input['data']['chart']['x_data'] as $i => $day) {
            foreach ($this->yDataMap as $key => $item) {
                $chart[] = [
                    'value' => $input['data']['chart']['y_data'][$i][$key],
                    'name'  => $item,
                    'date'  => date('Y-m-d', $day)
                ];
            }
        }
        unset($input['data']['chart']);
        $input['data']['chart']                    = $chart;
        $paginate['list']                          = $input['data']['paginate']['data'];
        $paginate['total']                         = $input['data']['paginate']['meta']['total'];
        $paginate['page']                          = $input['data']['paginate']['meta']['page'];
        $paginate['size']                          = $input['data']['paginate']['meta']['size'];
        $input['data']['paginate']                 = $paginate;
        $input['data']['data_overview']['outflow'] = $input['data']['data_overview']['customer_total'] - $input['data']['data_overview']['pure_customer_total'];
        return $input['data'];
    }

    /**
     * 返回默认数据
     *
     * @param array $input
     * @return array
     */
    private function defaultResult(array $input): array
    {
        $result = [
            'chart'         => [],
            'paginate'      => [
                'type'  => $input['data_type'],
                'list'  => [],
                'total' => 0,
                'page'  => 1,
                'size'  => $input['size']
            ],
            'data_overview' => [
                'customer_total'      => 0,
                'pure_customer_total' => 0,
                'outflow'             => 0
            ]
        ];

        $startTime = strtotime(date('Y-m-d 00:00:00', $input['start_time']));
        $endTime   = strtotime(date('Y-m-d 23:59:59', $input['end_time']));
        for ($day = $startTime; $day < $endTime; $day += 86400) {
            foreach ($this->yDataMap as $item) {
                array_push($result['chart'], [
                    'value' => 0,
                    'name'  => $item,
                    'date'  => date('Y-m-d', $day)
                ]);
            }
        }

        return $result;
    }
}
