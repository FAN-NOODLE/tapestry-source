<?php
/**
 * Created By PhpStorm
 * @Date: 2021/7/2
 * @Time: 17:30
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: CommunityListByMemberForConsoleController.php
 */

namespace App\Http\Controllers\Service\Brand\CommunityStatistic;


class CommunityListByMemberForConsoleController extends BaseController
{
    /**
     * 导购社群明细列表（以及导出）（按导购）
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @Date: 2021/6/29
     * @Time: 16:52
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'integer|required',
            'qy_domain'     => 'required|string',
            'start_time'    => 'required|integer',
            'end_time'      => 'required|integer',
        ], [
            'enterprise_id.*'    => '企业id参数不合法',
            'qy_domain.required' => 'qy_domain参数不合法',
            'start_time.*'       => 'start_time参数不合法',
            'end_time.*'         => 'end_time参数不合法',
        ]);

        $params['enterprise_id']  = $this->request->input('enterprise_id');
        //$params['store_id']       = $this->request->input('store_id', 0);
        //$params['worker_ids']     = $this->request->input('worker_ids', []);
        $params['qy_domain']      = $this->request->input('qy_domain');
        $params['start_time']     = $this->request->input('start_time');
        $params['end_time']       = $this->request->input('end_time');
        //$params['department_ids'] = $this->request->input('department_ids');
        //$params['store_ids']      = $this->request->input('store_ids');
        $params['size']           = $this->request->input('size', $this->communityDataCompassService::PAGE_LIMIT);
        $params['is_page']        = $this->request->input('is_page', true);

        $workerList = $this->listServicingWorker($params['enterprise_id']);

        $result = $this->communityDataCompassService->getCommunityStatisticsByMember(
            $params['enterprise_id'],
            $params['start_time'],
            $params['end_time'],
            $workerList,
            $params['size'],
            $params['is_page']
        );

        if ($params['is_page']) {
            $return = [
                'total' => $result['total'],
                'page'  => $result['current_page'],
                'size'  => $result['per_page'],
                'list'  => $result['data'],
            ];
        } else {
            $return = $result['data'];
        }

        return success($return);
    }
}