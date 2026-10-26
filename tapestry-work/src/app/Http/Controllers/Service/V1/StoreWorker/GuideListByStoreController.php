<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\Scrm\Departments;
use App\Models\Scrm\Members;
use App\Models\StoreWorker;
use Exception;

class GuideListByStoreController extends BaseController
{
    /**
     * 订单统计导购列表
     * type (1：推广订单数 2：推广订单实收金额 3：服务订单数 4：服务订单实收金额)
     * role (1：店员 2：店长 3：企业导购 4：社群导购 5：企业管理者)
     *
     * @return array
     * @throws Exception
     */
    public function index()
    {
        $input    = $this->request->input();
        $rules    = [
            'store_ids'     => 'array',
            'dp_ids'        => 'array',
            'third_ids'     => 'array',
            'domain'        => 'required',
            'enterprise_id' => 'required|integer',
            'worker_name'   => 'string',
            'worker_status' => 'integer',

        ];
        $messages = [
            'store_ids.array'        => '门店id必须是数组',
            'dp_ids.array'           => '部门id必须是数组',
            'third_ids.array'        => 'third_ids必须是数组',
            'domain.required'        => '企业domain不能为空',
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID必须为正整数',
            'worker_name.string'     => '导购名称必须为字符串',
            'worker_status.integer'  => '导购状态必须为int类型',
        ];
        $this->validate($this->request, $rules, $messages);
        $guideList = [];
        if (!empty($input['store_ids'])) {
            // 获取门店下的导购
            $guideList = $this->getGuideByStore($input);
        }
        if (!empty($input['dp_ids'])) {
            // 获取部门下的导购
            $guideList = $this->getGuideByDepartment($input);
        }
        if (!empty($input['worker_name'])) {
            // 依据导购编码/名称获取导购
            $guideList = $this->getGuideByNameOrNo($input);
        }

        return success($guideList);
    }

    /**
     * 根据门店获取导购
     *
     * @param $input
     *
     * @return array
     * @throws
     */
    protected function getGuideByStore($input)
    {
        $guideInfo = StoreWorker::query()
            ->whereIn('store_id', $input['store_ids'])
            ->where('enterprise_id', $input['enterprise_id'])
            ->when(!empty($input['worker_status']), function ($query) use ($input) {
                return $query->where('worker_status', $input['worker_status']);
            })
            ->select(['id', 'store_id', 'worker_no', 'name', 'mobile', 'type', 'worker_status', 'mem_uid', "mem_userid"])
            ->get()
            ->toArray();
        if ($guideInfo) {
            $userIds       = array_column($guideInfo, "mem_userid");
            $members       = Members::query()->whereIn("userid", $userIds)->where('enterprise_id', $input['enterprise_id'])->select(['id', 'userid', 'status'])->get()->toArray();
            $membersUserid = array_column($members, null, "userid");
            foreach ($guideInfo as &$value) {
                $value['store_name']    = '';
                $value['worker_status'] = $membersUserid[$value['mem_userid']]['status'] == Members::MEMBERS_STATUS_ACTIVATED ? 1 : 2;
                $value['member_id']     = $membersUserid[$value['mem_userid']]['id'];
            }
        }
        return $guideInfo;
    }

    /**
     * 获取部门下的导购
     *
     * @param $input
     *
     * @return array
     * @throws
     */
    protected function getGuideByDepartment($input)
    {

        if (count($input['dp_ids']) == 1 && is_string($input['dp_ids'][0])) { //兼容如果是sw或者ks的选择多个门店的时候，dp_ids传的数据是一个长度为1的多个门店id字符串
            $input['dp_ids'] = explode(",", $input['dp_ids'][0]);
        }
        $dpList    = Departments::query()->with('membersRel.members')
            ->where('enterprise_id', $input['enterprise_id'])
            ->whereIn('third_id', $input['dp_ids'])
            ->orWhereIn('id', $input['dp_ids'])
            ->get()->toArray();
        $memberArr = [];
        foreach ($dpList as $val) {
            foreach ($val['members_rel'] as $rel) {
                $member = $rel['members'];
                if (!$member) continue;
                $memberArr[$member['userid']] = [
                    "store_id"=>$val["third_id"],
                    'store_name'    => $val['name'],
                    'name'          => $member['name'],
                    'mobile'        => $member['mobile'],
                    'type'          => $member['type'],
                    'worker_status' => $member['status'] == Members::MEMBERS_STATUS_ACTIVATED ? 1 : 2,
                    'mem_userid'    => $member['userid'],
                    'member_id'     => $member['id']
                ];
            }
        }

        $userIdToMemUid = StoreWorker::query()->where('enterprise_id', $input['enterprise_id'])
            ->whereIn('mem_userid', array_keys($memberArr))->get(['mem_uid', 'mem_userid'])
            ->keyBy('mem_userid')->toArray();

        foreach ($memberArr as &$val) {
            $mem_uid = $userIdToMemUid[$val['mem_userid']]['mem_uid'] ?? '';
            if ($mem_uid) {
                $val['mem_uid'] = $mem_uid;
            }
        }
        return array_values($memberArr);
    }

    /**
     * 根据导购名称获取导购列表
     *
     * @param $input
     *
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/2 10:14
     */
    public function getGuideByNameOrNo($input)
    {
        $guideInfo = StoreWorker::query()
            ->where('name', 'like', '%' . $input['worker_name'] . '%')
            ->where('enterprise_id', $input['enterprise_id'])
            ->when(!empty($input['worker_status']), function ($query) use ($input) {
                return $query->where('worker_status', $input['worker_status']);
            })
            ->select(['id', 'store_id', 'worker_no', 'name', 'mobile', 'type', 'worker_status', 'mem_uid', "mem_userid"])
            ->get()
            ->toArray();
        if ($guideInfo) {
            $userIds       = array_column($guideInfo, "mem_userid");
            $members       = Members::query()->whereIn("userid", $userIds)->where('enterprise_id', $input['enterprise_id'])->select(['id', 'userid', 'status'])->get()->toArray();
            $membersUserid = array_column($members, null, "userid");
            foreach ($guideInfo as &$value) {
                $value['store_name']    = '';
                $value['worker_status'] = $membersUserid[$value['mem_userid']]['status'] == Members::MEMBERS_STATUS_ACTIVATED ? 1 : 2;
                $value['member_id']     = $membersUserid[$value['mem_userid']]['id'];
            }
        }

        return $guideInfo;
    }

    /**
     * 查询人员列表
     *
     * @param $condition
     *
     * @return array
     */
    protected function memberList($condition)
    {
        $page       = 1;
        $maxPages   = 0;
        $memberList = [];
        do {
            $condition['params']['pageNum']  = $page;
            $condition['params']['pageSize'] = 200;
            // 结果处理
            $result = (new Member())->ignorePermissionsList($condition['qy_domain'], '', $condition['params']);
//            $result     = (new Member())->memberList($condition);
            $memberList = array_merge($memberList, $result['list']);
            $maxPages   = $maxPages ?: ceil($result['total'] / 200);
            // 步程自增
            $page++;
        } while ($page <= $maxPages);
        return $memberList;
    }

}