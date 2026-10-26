<?php
/**
 * Created by PhpStorm.
 * User: baizhongyang
 * Date: 2020/8/18
 * Time: 16:11
 */

namespace App\Http\Controllers\Service\V1\Order;

use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use App\Services\GuideManagerService;
use App\Http\Controllers\Service\V1\BaseController;

class OrderStatisticalRankingController extends BaseController
{

    /**
     * 订单统计排行
     *
     * @param Request $request
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/18 11:49
     */
    public function index(Request $request)
    {
        $input    = $request->input();
        $rules    = [
            'role'          => 'required',
            'start_time'    => 'integer',
            'end_time'      => 'integer',
            'type'          => 'required|integer|in:1,2,3,4',
            'qy_domain'     => 'required',
            'mobile'        => 'required',
            'enterprise_id' => 'required|integer',
            'store_id'      => 'required|integer',
            'admin_id'      => 'required',
        ];
        $messages = [
            'role.required'          => '员工角色不能为空',
            'work_ids.required'      => '员工ID不能为空',
            'type.required'          => 'type不能为空',
            'qy_domain.required'     => '企业domain不能为空',
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID必须为正整数',
            'store_id.required'      => '门店id不能为空',
            'admin_id.required'      => 'admin_id不能为空',
        ];
        $this->validate($request, $rules, $messages);

        switchAppHost(['enterprise_id' => $input['enterprise_id']]);

        // 获取对应导购列表
        $guideList = GuideManagerService::getGuideListForRank($input['mobile'], $input['qy_domain'], $input['enterprise_id'], $input['admin_id'], $input['store_id']);
        if (!empty($guideList['guide_list'])) {
            $workIds = array_column($guideList['guide_list'], 'id');
        } else {
            return success();
        }
        // 调shop获取订单排行信息
        $params      = [
            'role'          => $input['role'],
            'start_time'    => $input['start_time'] ?? 0,
            'end_time'      => $input['end_time'] ?? 0,
            'type'          => $input['type'],
            'worker_ids'    => $workIds,
            'domain'        => $input['qy_domain'],
            'enterprise_id' => $input['enterprise_id'],
        ];
        //查询导购订单
        $guideOrderList = (new Order())->orderRankingList($params);
        //导购列表
        $guideListInfo = $guideList['guide_list'];
        $data = [];
        foreach ($guideListInfo as $value){
            switch (intval($input['type'])) {
                case 1:
                    $guideOrderList = array_column($guideOrderList, null, 'store_staff_id');
                    $data[] = [
                        'id'         => $value['id'],
                        'num'        => isset($guideOrderList[$value['id']]) ? $guideOrderList[$value['id']]['order_num'] : 0,
                        'store_name' => $value['store_name'],
                        'name'       => $value['name'],
                        'avatar'     => $value['avatar']
                    ];
                    break;
                case 2:
                    $guideOrderList = array_column($guideOrderList, null, 'store_staff_id');
                    $data[] = [
                        'id'         => $value['id'],
                        'num'        => isset($guideOrderList[$value['id']]) ? $guideOrderList[$value['id']]['sale_total'] : 0,
                        'store_name' => $value['store_name'],
                        'name'       => $value['name'],
                        'avatar'     => $value['avatar'],
                    ];
                    break;
                case 3:
                    $guideOrderList = array_column($guideOrderList, null, 'service_staff_id');
                    $data[] = [
                        'id'         => $value['id'],
                        'num'        => isset($guideOrderList[$value['id']]) ? $guideOrderList[$value['id']]['service_num'] : 0,
                        'store_name' => $value['store_name'],
                        'name'       => $value['name'],
                        'avatar'     => $value['avatar'],
                    ];
                    break;
                case 4:
                    $guideOrderList = array_column($guideOrderList, null, 'service_staff_id');
                    $data[] = [
                        'id'         => $value['id'],
                        'num'        => isset($guideOrderList[$value['id']]) ? $guideOrderList[$value['id']]['sale_total'] : 0,
                        'store_name' => $value['store_name'],
                        'name'       => $value['name'],
                        'avatar'     => $value['avatar'],
                    ];
                    break;

            }
        }
        $sortField = array_column($data,'num');
        array_multisort($sortField,SORT_DESC,$data);
        count($data) > 50 ? $data = array_slice($data, 0, 50) : $data;
        return success($data);
    }
}
