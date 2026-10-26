<?php
/**
 * Created by PhpStorm.
 * Function: 企业员工业绩排行
 * User: 龚国新
 * Date: 2020/2/14
 * Time: 11:43
 */

namespace App\Http\Controllers\App\V1\StaffOrder;


use App\Http\Controllers\App\V1\BaseController;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;

class AchievementRankController extends BaseController
{
//    protected $requireLogin = false;

    /**
     * 企业员工的业绩排名
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $params = $this->validates();
//        $this->memUid       = '5647FD087F00000153FAE204CA99D21B';
//        $this->enterpriseId = 2;
//        $this->qyDomain     = '72D5430F7F000001381C7A561A35C79D';
        $total                    = (new StaffOrderService())->getRankTotal($params, $this->enterpriseId);
        $noOrderStaff             = (new StaffOrderService())->getNoOrderStaff($params, $this->enterpriseId);
        $current_staff_order_info = (new StaffOrderService())->getCurrentMemRank($params, $this->enterpriseId, $this->memUid, $noOrderStaff);
        $allStatisticsInfo        = (new StaffOrderService())->getMemStatisticsRank($params, $this->enterpriseId, $noOrderStaff);

        $memUids    = [];
        $mem_params = [];

        //开始时就将当前员工存入到要查的数组里
        $memUids[] = $this->memUid;

        //查询员工的信息
        foreach ($allStatisticsInfo as $key => $value) {
            if ($this->memUid == $value->mem_uid) {

            } else {
                $memUids[] = $value->mem_uid;
            }
        }

        if (!empty($memUids)) {
            $mem_params = [
                'memUids' => $memUids
            ];
        }

        $result = (new StaffManagerService())->getStaffListFromUc($this->qyDomain, $this->appIdentifier, $mem_params);

        // 格式化数据
        $memAllInfo = (new StaffManagerService())->formatUcMember($result);
        $memAllInfo = array_column($memAllInfo, null, 'mem_uid');

        //格式化输出数据
        $formatData = $this->formatAllData($memAllInfo, $allStatisticsInfo, $current_staff_order_info);

        return success(['data' => [
            'current_staff_order' => $formatData['current_staff_order'],
            'all_staff_order'     => $formatData['all_staff_order'],
            'meta'                => [
                'page'  => $params['page'],
                'size'  => $params['size'],
                'total' => $total
            ]
        ]]);
    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'query_time' => 'required|string',
            'page'       => 'required|int',
            'size'       => 'required|int',
        ];

        $message = [
            'query_time.required' => '查询时间区间不能为空',
            'query_time.string'   => '查询时间区间数据类型不正确',
            'page.required'       => '当前页码不能为空',
            'page.int'            => '当前页码数据类型不正确',
            'size.required'       => '每页条数不能为空',
            'size.int'            => '每页条数数据类型不正确',
        ];

        $this->validate(request(), $rules, $message);

        $params = request()->input();

        return $params;
    }

    private function formatAllData($memAllInfo, $statisticsInfo, $current_staff_order_info)
    {
        $formatAllData['all_staff_order']                = [];
        $formatAllData['current_staff_order']['mem_uid'] = $this->memUid;
        if (!empty($memAllInfo[$this->memUid])) {
            $formatAllData['current_staff_order']['name']     = $memAllInfo[$this->memUid]['mem_username'];
            $formatAllData['current_staff_order']['mem_face'] = $memAllInfo[$this->memUid]['mem_face'];
        } else {
            $formatAllData['current_staff_order']['name']     = '';
            $formatAllData['current_staff_order']['mem_face'] = '';
        }

        if (!empty($current_staff_order_info)) {
            $formatAllData['current_staff_order']['order_count']  = $current_staff_order_info[0]->order_total_count;
            $formatAllData['current_staff_order']['order_amount'] = $current_staff_order_info[0]->order_total_money;
            $formatAllData['current_staff_order']['rownum']       = $current_staff_order_info[0]->rownum;
        } else {
            $formatAllData['current_staff_order']['order_count']  = 0;
            $formatAllData['current_staff_order']['order_amount'] = 0;
            $formatAllData['current_staff_order']['rownum']       = 0;
        }

        foreach ($statisticsInfo as $item) {
            if (!empty($memAllInfo[$item->mem_uid])) {
                $temp['name']     = $memAllInfo[$item->mem_uid]['mem_username'];
                $temp['mem_face'] = $memAllInfo[$item->mem_uid]['mem_face'];
            } else {
                $temp['name']     = '';
                $temp['mem_face'] = '';
            }
            $temp['mem_uid']                    = $item->mem_uid;
            $temp['order_count']                = $item->order_total_count;
            $temp['order_amount']               = (int)$item->order_total_money;
            $temp['rownum']                     = $item->rownum;
            $formatAllData['all_staff_order'][] = $temp;
        }

        return $formatAllData;
    }
}