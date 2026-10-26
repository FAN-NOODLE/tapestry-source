<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/12/5
 * Time: 10:11
 */

namespace App\Http\Controllers\App\V1\Rank;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Services\DateService;
use App\Services\MemberService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class CustomerRankController extends BaseController
{
    /**
     * 我的排行榜-客户数榜
     *
     * @param Request $request
     * @param MemberService $memberService
     * @param DateService $dateService
     * @return array
     * @throws \Exception
     */
    public function index(Request $request, MemberService $memberService, DateService $dateService)
    {

        // 数据时间范围（1：今日，2：本周，3：本月）
        $dateType = $request->input('date_type', 1);
        // 页码：默认1
        $page = $request->input('page', 1);
        // 每页数量：默认15
        $size = $request->input('size', 15);

        // 获取人员列表
        $list = $memberService->getMemberListByUid($this->qyDomain, $this->appIdentifier, $this->memUid);

        // 获取时间起止
        $dateArray = $dateService->getDateTime($dateType, 1);

        // 分组统计各个员工的客户总数
        $allCustomerList = MemberUserRelation::query()
            ->select(DB::raw('mem_uid,count(*) as total'))
            ->where('enterprise_id', $this->enterpriseId)
            ->where([
                ['created_at', '>=', $dateArray['start_time']],
                ['created_at', '<=', $dateArray['end_time']]
            ])
            ->groupBy('mem_uid')
            ->get();

        $allCustomerData = [];
        if ($allCustomerList) {
            $allCustomerData = $allCustomerList->toArray();
            $allCustomerData = array_column($allCustomerData, null, 'mem_uid');
        }

        // 格式化列表
        $list = $this->formatData($list, $allCustomerData);

        // 获取客户数组
        $customer_nums = array_column($list, 'customer_num');

        // 根据客户数降序排列
        array_multisort($customer_nums, SORT_DESC, $list);

        $i = 1;

        //排名顺序
        $j = 1;
        // 当前用户排名数据
        $currentData = [];
        // 组装排序数据
        foreach ($list as $key => &$v) {
            // 当前导购客户数不等于上一个导购客户数，排名自增
            if ($key > 0 && $v['customer_num'] != $list[$key - 1]['customer_num']) {
                $j++;
            }
            
            $v['index_num'] = $j;

            if ($v['mem_uid'] == $this->memUid) {
                $currentData = $v;
            }

            $i++;
        }

        unset($v);

        // 数组分页
        $list = $dateService->pageArray($size, $page, $list);

        // 数据返回
        return success([
            'data' => [
                'current_data' => $currentData, // 当前用户所在排名
                'list'         => $list,
                'meta'         => [
                    'total' => ($i - 1),
                    'page'  => $page,
                    'size'  => $size
                ]
            ]
        ]);

    }

    /**
     * 格式化数据返回列表
     *
     * @param array $list 员工列表
     * @param array $allCustomerData 员工的客户统计数据列表
     * @return array
     */
    protected function formatData($list = [], $allCustomerData = [])
    {


        $result = [];
        foreach ($list as $v) {

            if ($v['memStatus'] == MemberService::MEMBER_DELETE_REMARK) {
                // 如果人员已删除，则略过
                continue;
            }

            // 循环拼装数据
            $data = [
                'mem_uid'      => $v['memUid'],
                'mem_username' => $v['memUsername'],
                'mem_face'     => $v['memFace'],
            ];

            // 组装客户数据
            if (isset($allCustomerData[$v['memUid']])) {
                $data['customer_num'] = intval($allCustomerData[$v['memUid']]['total']);
            } else {
                $data['customer_num'] = 0;
            }

            $result[] = $data;
        }

        return $result;
    }
}
