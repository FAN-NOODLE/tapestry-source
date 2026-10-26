<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\CommissionRecord;
use App\Models\MemberExtend;
use App\Services\DateService;
use App\Services\MemberService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class CommissionController extends BaseController
{
    /**
     * 佣金数据接口
     *
     * @param Request $request
     * @param MemberService $memberService
     * @param DateService $dateService
     * @return array
     * @throws \Exception
     */
    public function index(Request $request, MemberService $memberService, DateService $dateService)
    {
        // 获取日期类型
        $dateType = $request->input('date_type', 1);

        $page = $request->input('page', 1);

        $size = $request->input('size', 15);

        // 获取人员列表
        $list = $memberService->getMemberListByUid($this->qyDomain, $this->appIdentifier, $this->memUid);

        // 获取员工uid集合
        $memUids = array_column($list, 'memUid');

        // 获取时间起止
        $dateArray = $dateService->getDateTime($dateType, 1);

        // 分组统计各个员工的佣金总数
        $allCommissionList = CommissionRecord::query()
            ->select(DB::raw('mem_uid,sum(charge_price) as total'))
            ->whereIn('mem_uid', $memUids)
            ->where([
                ['created_at', '>=', $dateArray['start_time']],
                ['created_at', '<=', $dateArray['end_time']]
            ])
            ->groupBy('mem_uid')
            ->get();

        $allCommissionData = [];
        if ($allCommissionList) {
            // 将查询结果转换为数组
            $allCommissionData = $allCommissionList->toArray();
        }

        // 将数据转换成以员工uid为主键的二维数组
        $allCustomerData = array_column($allCommissionData, null, 'mem_uid');

        $result     = [];
        $priceTotal = 0;
        foreach ($list as $v) {

            // 循环拼装数据
            $data = [
                'mem_uid'      => $v['memUid'],
                'mem_username' => $v['memUsername'],
                'mem_face'     => $v['memFace'],
                'mem_status'   => $v['memStatus'] != MemberService::MEMBER_DELETE_REMARK ? 0: 1,
            ];

            // 组装客户数据
            if (isset($allCustomerData[$v['memUid']])) {
                // 将佣金转换为元
                $data['price'] = sprintf('%.2f', intval($allCustomerData[$v['memUid']]['total']) / 100);
            } else {
                $data['price'] = sprintf('%.2f', 0);
            }

            $priceTotal += $data['price'];

            $result[] = $data;
        }

        // 取出佣金金额列
        $prices = array_column($result, 'price');

        if (!empty($prices)) {
            // 根据佣金金额列降序排列
            array_multisort($prices, SORT_DESC, $result);
        }

        // 拼装排序数据
        $i = 1;
        foreach ($result as $k=>&$v) {
            if ($k > 0) {
                if ($result[$k-1]['price'] != $result[$k]['price']) {
                    $i++;
                }
            }
            $v['index_num'] = $i;
        }

        unset($v);

        // 数组分页
        $list = $dateService->pageArray($size, $page, $result);

        // 数据返回
        return success([
            'data' => [
                'price_total' => sprintf('%.2f', $priceTotal), // 计算佣金总额
                'start_time'  => date('Y.m.d', $dateArray['start_time']),
                'end_time'    => date('Y.m.d', $dateArray['end_time']),
                'list'        => $list,
                'meta'        => [
                    'total' => count($result),
                    'page'  => $page,
                    'size'  => $size
                ]
            ]
        ]);

    }

}
