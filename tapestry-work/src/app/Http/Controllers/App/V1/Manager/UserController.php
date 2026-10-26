<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Services\DateService;
use App\Services\MemberService;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class UserController extends BaseController
{
    /**
     * 客户数据列表接口
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
        // 排序类型（1：按客户数，2：按会员数）
        $sortType = $request->input('sort_type', 1);
        // 页码：默认1
        $page = $request->input('page', 1);
        // 每页数量：默认15
        $size = $request->input('size', 15);

        // 获取人员列表
        $list = $memberService->getMemberListByUid($this->qyDomain, $this->appIdentifier, $this->memUid);

        // 获取员工uid集合
        $memUids = array_column($list, 'memUid');

        // 获取时间起止
        $dateArray = $dateService->getDateTime($dateType, 1);

        // 分组统计各个员工的客户总数
        $allCustomerList = MemberUserRelation::query()
            ->select(DB::raw('mem_uid,count(*) as total'))
            ->whereIn('mem_uid', $memUids)
            ->where([
                ['created_at', '>=', $dateArray['start_time']],
                ['created_at', '<=', $dateArray['end_time']]
            ])
            ->groupBy('mem_uid')
            ->get();

        $allCustomerData = [];
        if ($allCustomerList) {
            $allCustomerData = $allCustomerList->toArray();
        }

        $allCustomerData = array_column($allCustomerData, null, 'mem_uid');

        // 调用商城会员sdk获取会员数据
        $result = (new User())->userCountByMemUids([
            'enterprise_id' => $this->enterpriseId,
            'mem_uids'      => $memUids,
            'start_time'    => $dateArray['start_time'],
            'end_time'      => $dateArray['end_time']
        ]);

        $allUserData = [];
        if ($result && isset($result['data'])) {
            $allUserData = array_column((array)$result['data'], null, 'mem_uid');
        }

        // 格式化列表
        $resultData = $this->formatData($list, $allCustomerData, $allUserData);

        $list = $resultData['list'];

        $customer_nums = array_column($list, 'customer_num');
        $users_nums    = array_column($list, 'user_num');
        // 二维数组根据相应条件排序
        if ($sortType == 1 && $customer_nums) {
            // 根据客户数降序排列
            array_multisort($customer_nums, SORT_DESC, $list);

        } elseif ($sortType == 2 && $users_nums) {
            // 根据会员数降序排列
            array_multisort($users_nums, SORT_DESC, $list);
        }

        $i = 1;
        foreach ($list as $k=>&$v) {
            // 按客户数
            if ($k > 0 && $sortType == 1) {
                if ($list[$k-1]['customer_num'] != $list[$k]['customer_num']) {
                    $i++;
                }
            }
            // 按会员数
            if ($k > 0 && $sortType == 2) {
                if ($list[$k-1]['user_num'] != $list[$k]['user_num']) {
                    $i++;
                }
            }
            $v['index_num'] = $i;
        }

        unset($v);

        // 数组分页
        $list = $dateService->pageArray($size, $page, $list);

        // 数据返回
        return success([
            'data' => [
                'customer_total' => $resultData['customer_total'],
                'user_total'     => $resultData['user_total'],
                'start_time'     => date('Y.m.d', $dateArray['start_time']),
                'end_time'       => date('Y.m.d', $dateArray['end_time']),
                'list'           => $list,
                'meta'           => [
                    'total' => count($resultData['list']),
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
     * @param array $allUserData 员工的会员统计数据列表
     * @return array
     */
    protected function formatData($list = [], $allCustomerData = [], $allUserData = [])
    {

        $customerNumTotal = 0;
        $userNumTotal     = 0;
        $result           = [];
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
                $data['customer_num'] = intval($allCustomerData[$v['memUid']]['total']);
            } else {
                $data['customer_num'] = 0;
            }

            $customerNumTotal += $data['customer_num'];

            // 组装会员数据
            if (isset($allUserData[$v['memUid']])) {
                $data['user_num'] = intval($allUserData[$v['memUid']]['total']);
            } else {
                $data['user_num'] = 0;
            }

            $userNumTotal += $data['user_num'];
            
            $result[] = $data;
        }

        return ['list' => $result, 'customer_total' => $customerNumTotal, 'user_total' => $userNumTotal];
    }

}
