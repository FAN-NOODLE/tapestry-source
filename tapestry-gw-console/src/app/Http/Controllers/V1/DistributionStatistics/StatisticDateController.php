<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-02
 * Time: 16:15
 */

namespace App\Http\Controllers\V1\DistributionStatistics;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionStatistics;
use Illuminate\Http\Request;

class StatisticDateController extends BaseController
{
    /**
     * 分销统计按时间筛选
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $params                  = $request->only(['worker_ids', 'start_time', 'end_time', 'sort_type', 'statistical_sort','page','size']);
        $params['enterprise_id'] = $this->enterpriseId;
        $data                    = (new DistributionStatistics())->customerStatisticDate($params);
        return success($data);
    }
}
