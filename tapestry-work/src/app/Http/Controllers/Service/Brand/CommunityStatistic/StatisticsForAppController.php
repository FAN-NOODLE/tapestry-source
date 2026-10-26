<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 14:29
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: StatisticsForAdminController.php
 */

namespace App\Http\Controllers\Service\Brand\CommunityStatistic;


use App\Services\GuideManagerService;

class StatisticsForAppController extends BaseController
{
    /**
     * 社群统计
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @Date: 2021/6/30
     * @Time: 9:19
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'worker_ids'    => 'array',
            'qy_domain'     => 'required|string',
            'admin_id'      => 'required',
            'start_time'    => 'required|integer',
            'end_time'      => 'required|integer',
        ], [
            'enterprise_id.*'    => '企业id参数不合法',
            'store_id.*'         => '门店id参数不合法',
            'mobile.*'           => '用户手机号参数不合法',
            'qy_domain.required' => 'qy_domain参数不合法',
            'admin_id.required'  => 'admin_id参数缺失',
            'start_time.*'       => 'start_time参数不合法',
            'end_time.*'         => 'end_time参数不合法',
        ]);

        $params['enterprise_id']  = $this->request->input('enterprise_id');
        $params['store_id']       = $this->request->input('store_id', 0);
        $params['mobile']         = $this->request->input('mobile');
        $params['worker_ids']     = $this->request->input('worker_ids', []);
        $params['qy_domain']      = $this->request->input('qy_domain');
        $params['admin_id']       = $this->request->input('admin_id');
        $params['start_time']     = $this->request->input('start_time');
        $params['end_time']       = $this->request->input('end_time');
        $params['department_ids'] = $this->request->input('department_ids', []);
        $params['store_ids']      = $this->request->input('store_ids');

        $workerQueryParams = [
            'worker_ids'     => $params['worker_ids'],
            'store_ids'      => $params['store_ids'],
            'department_ids' => $params['department_ids'],
        ];

        $workerList = GuideManagerService::getGuideListForStatistic($workerQueryParams, $params['mobile'], $params['qy_domain'], $params['enterprise_id'], $params['admin_id'], $params['store_id']);

        if (!empty($workerList['guide_list'])) {
            $workerList = $workerList['guide_list'];

            // 获取scrm数据库中member数据表的相关导购数据id
            $memberIds = $this->communityDataCompassService->getScrmMember($params['enterprise_id'], $workerList, ['id']);

            $result = $this->communityDataCompassService->overview(
                $params['enterprise_id'],
                $params['start_time'],
                $params['end_time'],
                $memberIds
            );

            return success($result);
        } else {
            return success($this->communityDataCompassService::GROUP_DATA_FIELDS);
        }

    }
}