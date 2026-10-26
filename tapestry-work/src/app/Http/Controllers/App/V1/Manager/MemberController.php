<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\MemberChangeRecord;
use App\Models\MemberUserRelation;
use App\Services\MemberService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class MemberController extends BaseController
{
    /**
     * UC查询删除人员状态标识
     */
    const MEMBER_DELETE_REMARK = 3;

    /**
     * 获取当前部门员工列表接口
     *
     * @param Request $request
     * @param MemberService $memberService
     * @return array|mixed
     * @throws \Exception
     */
    public function index(Request $request, MemberService $memberService)
    {

        // 获取数据类型：1：全部在职，2：离职未分配，3：离职已分配
        $dataType = $request->input('data_type', 0);

        $memUsername = $request->input('mem_username', '');

        // 获取当前登录用户的详情信息
        $ucResult = (new Member())->detail($this->qyDomain, $this->appIdentifier, ['memUid' => $this->memUid]);

        if (!isset($ucResult)) {
            error(60103);
        }

        $dpId = '';
        if (isset($ucResult['dpName'])) {
            // 取出所属部门id
            $dpId = $ucResult['dpName'][0]['dpId'];
        }

        if (empty($dpId)) {
            // 部门信息不存在
            error(600102);
        }

        $where['dpId'] = $dpId;
        if ($memUsername) {
            $where['memUsername'] = $memUsername;
        }

        if ($dataType != 1) {
            // 如果需要查询删除人员列表
            $where['eqMemStatus'] = self::MEMBER_DELETE_REMARK;
        }

        // 获取人员列表
        $list = $memberService->listAllMember($this->qyDomain, $this->appIdentifier, $where);

        // 格式化返回数据
        $data = $this->formatList($list, $dataType);

        // 数据返回
        return success(['data' => $data]);
    }

    /**
     * 格式化数据
     *
     * @param array $list
     * @param int $dataType 数据类型
     * @return array
     */
    protected function formatList($list = [], $dataType = 1)
    {
        $newList = [];

        // 获取员工uid集合
        $memUids = array_column($list, 'memUid');

        // 分组统计各个员工的客户总数
        $allCustomerList = MemberUserRelation::query()
            ->select(DB::raw('mem_uid,count(*) as total'))
            ->whereIn('mem_uid', $memUids)
            ->groupBy('mem_uid')
            ->get();

        $allCustomerData = [];
        if ($allCustomerList) {
            $allCustomerData = $allCustomerList->toArray();
        }

        $allCustomerData = array_column($allCustomerData, null, 'mem_uid');

        $distributioData = [];
        // 如果是已离职员工列表
        if ($dataType != 1) {

            // 分组统计各个离职员工的已分配客户总数
            $distributionCustomerList = MemberChangeRecord::query()
                ->select(DB::raw('old_mem_uid,count(*) as total'))
                ->whereIn('old_mem_uid', $memUids)
                ->groupBy('old_mem_uid')
                ->get();

            if ($distributionCustomerList) {
                $distributioData = $distributionCustomerList->toArray();
            }

            $distributioData = array_column($distributioData, null, 'old_mem_uid');
        }

        // 循环格式化
        foreach ($list as $v) {

            // 已分配
            $distributioNum = isset($distributioData[$v['memUid']]) ? $distributioData[$v['memUid']]['total'] : 0;

            $nowCustomerNum = isset($allCustomerData[$v['memUid']]) ? $allCustomerData[$v['memUid']]['total'] : 0;

            // 总客户数=现有客户数+已分配客户数
            $data = [
                'mem_uid'      => $v['memUid'],
                'mem_username' => $v['memUsername'],
                'mem_face'     =>  $v['memFace'],
                'customer_num' => $nowCustomerNum + $distributioNum,
            ];

            if ($dataType == 1) {
                // 如果是在职员工，则待分配客户数为0
                $data['customer_no_change'] = 0;
            } else {
                // 待分配客户数=当前客户数
                $data['customer_no_change'] = $nowCustomerNum;
            }

            if ($dataType == 3 && $nowCustomerNum != 0 ) {
                // 如果传入的状态是离职已分配，并且待分配客户数不等于0，则不符合条件
                continue;
            }

            if ($dataType == 2 && $nowCustomerNum == 0) {
                // 如果传入的状态是离职待分配，并且待分配客户数等于0，则不符合条件
                continue;
            }

            $newList[] = $data;
        }

        return $newList;
    }
}
