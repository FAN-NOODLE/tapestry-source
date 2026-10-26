<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 9:51
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: StatisticsController.php
 */

namespace App\Http\Controllers\Brand\CommunityStatistics;


use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Work\DataCompass;

class StatisticsController extends BaseController
{

    /**
     * 获取社群统计数据
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @Date: 2021/6/24
     * @Time: 14:51
     */
    public function index(Request $request)
    {
        $params                  = $request->input();
        $params['start_time']    = $request->input('start_time');
        $params['end_time']      = $request->input('end_time');
        $params['worker_ids']    = $request->input('worker_ids', []);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['admin_id']      = $this->enterpriseAdminId;
        $params['qy_domain']     = $this->qyDomain;
        $params['mobile']        = isset($this->enterpriseAdmin['mobile']) && $this->enterpriseAdmin['mobile'] ? $this->enterpriseAdmin['mobile'] : '';

        $statistics = (new DataCompass())->countForAdmin($params);

        return success($statistics);
    }

    /**
     * 获取社群统计数据明细列表（按社群）
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @Date: 2021/6/24
     * @Time: 14:51
     */
    public function communityListByGroup(Request $request)
    {
        $params                  = $request->input();
        $params['start_time']    = $request->input('start_time');
        $params['end_time']      = $request->input('end_time');
        $params['worker_ids']    = $request->input('worker_ids', []);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain']     = $this->qyDomain;

        $statistics = (new DataCompass())->getCommunityStatisticListByGroup($params);

        return success($statistics);
    }

    /**
     * 获取导购社群统计数据明细列表（按导购）
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @Date: 2021/6/24
     * @Time: 14:51
     */
    public function communityListByMember(Request $request)
    {
        $params                  = $request->input();
        $params['start_time']    = $request->input('start_time');
        $params['end_time']      = $request->input('end_time');
        $params['worker_ids']    = $request->input('worker_ids', []);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain']     = $this->qyDomain;

        $statistics = (new DataCompass())->getCommunityStatisticListByMember($params);

        return success($statistics);
    }
}