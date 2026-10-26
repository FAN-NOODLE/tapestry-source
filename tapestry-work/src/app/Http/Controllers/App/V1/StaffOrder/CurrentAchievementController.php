<?php
/**
 * Created by PhpStorm.
 * Function: 统计当前企业员工业绩
 * User: 龚国新
 * Date: 2020/2/14
 * Time: 11:43
 */

namespace App\Http\Controllers\App\V1\StaffOrder;


use App\Http\Controllers\App\V1\BaseController;
use App\Models\StaffOrderRelation;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use Illuminate\Support\Facades\Log;

class CurrentAchievementController extends BaseController
{
//    protected $requireLogin = false;

    /**
     * 当前员工的业绩统计
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $params             = $this->validates();
//        $this->memUid       = '5647FD087F00000153FAE204CA99D21B';
//        $this->enterpriseId = 2;
//        $this->qyDomain     = '72D5430F7F000001381C7A561A35C79D';

        //查询员工的信息
        $mem_params = [
            'memUids' => [$this->memUid]
        ];
        Log::info("当前登录信息：".json_encode([$this->userId,$this->enterpriseId,$this->qyDomain,$this->memUid]));
        //查询当前员工是否有业绩，如果没有就插入一条记录
        (new StaffOrderService())->checkStaffOrderInfo($this->memUid,$this->enterpriseId);
        $result     = (new StaffManagerService())->getStaffListFromUc($this->qyDomain, $this->appIdentifier, $mem_params);

        // 格式化数据
        $memInfo = (new StaffManagerService())->formatUcMember($result);

        if(empty($memInfo)){
            Log::info("企业qyDomain".$this->qyDomain.",企业id：".$this->enterpriseId.",要查询的员工：".json_encode($mem_params));
            notice(200016);
        }

        //根据查询条件统计订单数、订单金额
        $statisticsInfo = (new StaffOrderService())->getMemStatisticsInfo($params, $this->memUid, $this->enterpriseId);

        //格式化输出数据
        $formatData = $this->formatData($memInfo, $statisticsInfo);

        return success(['data' => [
            'current_staff_info'  => $formatData['current_staff_info'],
            'current_staff_order' => $formatData['current_staff_order'],
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
        ];

        $message = [
            'query_time.required' => '查询时间区间不能为空',
            'query_time.string'   => '查询时间区间数据类型不正确',
        ];

        $this->validate(request(), $rules, $message);

        $params = request()->input();

        return $params;
    }

    private function formatData($memInfo, $statisticsInfo)
    {
        $formatData['current_staff_info']['mem_username']  = $memInfo[0]['mem_username'];
        $formatData['current_staff_info']['dp_name']       = $memInfo[0]['dp_name'];
        $formatData['current_staff_info']['mem_face']      = $memInfo[0]['mem_face'];
        $formatData['current_staff_info']['is_author']     = 1;

        //根据当前用户第一次成为内购员工的日期
        $first = StaffOrderRelation::query()->where('mem_uid',$this->memUid)->first();

        if(!empty($first)){
            $formatData['current_staff_order']['first_time'] = $first->created_at->timestamp;
        }else{
            $formatData['current_staff_order']['first_time'] = 0;
        }

        if(!empty($statisticsInfo)){
            $formatData['current_staff_order']['order_count']  = $statisticsInfo[0]->order_total_count;
            $formatData['current_staff_order']['order_amount'] = (int)$statisticsInfo[0]->order_total_money;
            $formatData['current_staff_order']['today']        = date('Y.m.d', time());
        }else{
            $formatData['current_staff_order']['order_count']  = 0;
            $formatData['current_staff_order']['order_amount'] = 0;
            $formatData['current_staff_order']['today']        = date('Y.m.d', time());
        }


        return $formatData;
    }
}