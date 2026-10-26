<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 14:29
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: StatisticsForAdminController.php
 */

namespace App\Http\Controllers\Service\Brand\CommunityStatistic;


class StatisticsForAdminController extends BaseController
{
    /**
     * 社群统计
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @Date: 2021/6/30
     * @Time: 9:18
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'worker_ids'    => 'array',
            'qy_domain'     => 'required|string',
            'start_time'    => 'required|integer',
            'end_time'      => 'required|integer',
        ], [
            'enterprise_id.*'    => '企业id参数不合法',
            'store_id.*'         => '门店id参数不合法',
            'qy_domain.required' => 'qy_domain参数不合法',
            'start_time.*'       => 'start_time参数不合法',
            'end_time.*'         => 'end_time参数不合法',
        ]);

        $params['enterprise_id']  = $this->request->input('enterprise_id');
        $params['store_id']       = $this->request->input('store_id', 0);
        $params['worker_ids']     = $this->request->input('worker_ids', []);
        $params['qy_domain']      = $this->request->input('qy_domain');
        $params['start_time']     = $this->request->input('start_time');
        $params['end_time']       = $this->request->input('end_time');
        $params['admin_id']       = $this->request->input('admin_id');
        $params['department_ids'] = $this->request->input('department_ids');
        $params['store_ids']      = $this->request->input('store_ids');
        $params['mobile']         = $this->request->input('mobile');

        $params['department_ids'] = array_filter($params['department_ids']);
        $params['worker_ids'] = array_filter($params['worker_ids']);

        if (!empty($params['department_ids']) || !empty($params['worker_ids'])) {
            $workerList = $this->communityDataCompassService->getDepartmentAndMemberIntersection($params['enterprise_id'], $params['qy_domain'], $params);

            if(empty($workerList)){
                return success($this->communityDataCompassService::GROUP_DATA_FIELDS);
            }
        }

        // 获取scrm数据库中member数据表的相关导购数据id
        $memberIds = $this->communityDataCompassService->getScrmMember($params['enterprise_id'], isset($workerList) ? $workerList : [], ['id']);

        $result = $this->communityDataCompassService->overview(
            $params['enterprise_id'],
            $params['start_time'],
            $params['end_time'],
            isset($memberIds) ? $memberIds : []
        );

        return success($result);
    }
}