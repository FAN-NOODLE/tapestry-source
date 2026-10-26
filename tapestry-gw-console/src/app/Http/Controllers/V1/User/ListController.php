<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use App\Common\Constant;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use CrsSdk\Work\MemberExtend;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 员工列表
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', 0);
        // 查询条件
        $search = $request->get('search', []);

        if (empty($this->qyDomain)) {
            return success([
                'list'  => [],
                'total' => 0,
                'page'  => intval($page),
                'size'  => intval($size)
            ]);
        }

        // 在职状态，默认在职
        $memActive = Constant::MEMBER_ACTIVE_TRUE;
        if (isset($search['mem_active'])) {
            $memActive = intval($search['mem_active']);
        }

        // 不分页，size重新赋值（根据部门id取员工列表选人时不需要分页）
        if (!$size) {
            $size = 1;
        }

        // 参数列表
        $params = [
            'pageNum'  => $page,
            'pageSize' => $size
        ];

        // 离职员工（已删除员工）
        if (Constant::MEMBER_ACTIVE_FALSE === $memActive) {
            $params['eqMemStatus'] = Constant::MEMBER_STATUS_DELETED;
            // 按删除时间倒序
            $params['orderList'][] = ['column' => 'memDeleted', 'orderType' => 'DESC'];
        } else {
            // 在职员工，按创建时间倒序
            $params['orderList'][] = ['column' => 'memCreated', 'orderType' => 'DESC'];
        }

        // 姓名
        if (isset($search['mem_username'])) {
            $params['memUsername'] = $search['mem_username'];
        }

        // 手机号码
        if (isset($search['mem_mobile'])) {
            $params['memMobile'] = $search['mem_mobile'];
        }

        // 门店id
        if (isset($search['dp_id']) && $search['dp_id']) {
            $params['dpId'] = $search['dp_id'];
        }

        // 角色
        if (isset($search['mem_role'])) {
            $params['memRole'] = $search['mem_role'];
        }

        // 调用sdk
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => $params,
            'enterprise_id'  => $this->enterpriseId
        ];

        $list = (new Member())->memberList($condition);

        // 总数
        $total = intval($list['total']);
        // 不分页 && 员工多余一个，查询列表
        if (1 == $size && $total > 1) {
            // size重新赋值
            $condition['params']['pageSize'] = $total;
            // 列表重新赋值
            $list = (new Member())->memberList($condition);
        }

        // 用户列表数据格式化
        $data = $this->format($memActive, $list['list']);

        // 有分页
        if ($size != 1) {
            return success([
                'list'  => $data,
                'total' => $total,
                'page'  => intval($page),
                'size'  => intval($size)
            ]);
        }

        return success(['list' => $data]);
    }

    /**
     * 格式化员工列表
     *
     * @param int $memActive 员工在职状态
     * @param array $list 待格式化员工列表
     *
     * @return array
     */
    protected function format($memActive, $list)
    {
        // 初始化待返数据
        $data = [];
        if (!$list) {
            return $data;
        }

        // 员工id集合
        $memUids = array_column($list, 'memUid');

        // 离职员工 - 客户分配状态列表
        $handoverStatusList = [];
        // 离职员工列表
        if (Constant::MEMBER_ACTIVE_FALSE === $memActive) {
            // 调用sdk获取员工客户分配状态
            $handoverStatusList = (new Customer())->handoverStatus(['mem_uids' => $memUids])['data'];
            $handoverStatusList = array_column($handoverStatusList, null, 'mem_uid');
            // 获取离职员工 - 待分配/待确认客户数
            $handoverUserCount = (new Customer())->handoverCount(['mem_uids' => $memUids, 'enterprise_id' => $this->enterpriseId])['data'];
        }
        // 批量获取员工头像列表
        $avatarList = (new MemberExtend())->extendList([
            'enterprise_id' => $this->enterpriseId,
            'mem_uids'      => $memUids
        ]);
        $avatarList = array_column($avatarList['data'], null, 'mem_uid');

        // 获取离职员工mem_uid集合
        foreach ($list as $member) {
            // 员工id
            $memUid = $member['memUid'];

            // 离职员工客户是否分配完成状态：未完成，当mem_status=2时有效
            $handoverStatus = Constant::CUSTOMER_HANDOVER_COMPLETED_FALSE;

            // 离职状态 - 格式化数据
            if (Constant::MEMBER_ACTIVE_FALSE === $memActive && isset($handoverStatusList[$memUid])) {
                $handoverStatus = $handoverStatusList[$memUid]['handover_status'];
                $data[]         = $this->formatMemberFalse($memUid, $member, $avatarList, $memActive, $handoverStatus, $handoverUserCount);
            }

            // 在职状态 - 格式化数据
            if (Constant::MEMBER_ACTIVE_TRUE === $memActive) {
                $data[] = $this->formatMemberTrue($memUid, $member, $avatarList, $memActive, $handoverStatus);
            }
        }

        return $data;
    }

    /**
     * 格式化 在职员工数据
     * @param $memUid
     * @param $member
     * @param $avatarList
     * @param $memActive
     * @param $handoverStatus
     * @return array
     */
    public function formatMemberTrue($memUid, $member, $avatarList, $memActive, $handoverStatus)
    {
        return [
            'mem_uid'         => $memUid,
            'mem_username'    => $member['memUsername'],
            'mem_face'        => $avatarList[$memUid]['avatar'] ?? $member['memFace'],
            'dp_id'           => $member['dpName'][0]['dpId'] ?? '',
            'dp_name'         => $member['dpName'][0]['dpName'] ?? '',
            'mem_role'        => $member['memRole'] ?? '',
            'mem_gender'      => $member['memGender'],
            'mem_mobile'      => $member['memMobile'],
            'mem_active'      => $memActive,
            'handover_status' => $handoverStatus,
            'created_at'      => !is_null($member['memSubscribeTime']) ? date('Y-m-d H:i:s', $member['memSubscribeTime'] / 1000) : '未加入',
            'user_total'      => $member['user_count'] ?? 0
        ];
    }

    /**
     * 格式化 离职员工数据
     * @param $memUid
     * @param $member
     * @param $avatarList
     * @param $memActive
     * @param $handoverStatus
     * @param $handoverUserCount
     * @return array
     */
    public function formatMemberFalse($memUid, $member, $avatarList, $memActive, $handoverStatus, $handoverUserCount)
    {
        return [
            'mem_uid'                 => $memUid,
            'mem_username'            => $member['memUsername'],
            'mem_face'                => $avatarList[$memUid]['avatar'] ?? $member['memFace'],
            'dp_id'                   => $member['dpName'][0]['dpId'] ?? '',
            'dp_name'                 => $member['dpName'][0]['dpName'] ?? '',
            'mem_role'                => $member['memRole'] ?? '',
            'mem_gender'              => $member['memGender'],
            'mem_mobile'              => $member['memMobile'],
            'mem_active'              => $memActive,
            'handover_status'         => $handoverStatus,
            'handover_no_distributed' => $handoverUserCount['mem_uid'][$memUid]['no_distributed'] ?? 0,
            'handover_wait_confirm'   => $handoverUserCount['mem_uid'][$memUid]['wait_confirm'] ?? 0,
            'deleted_at'              => $member['memDeleted'] === 0 ? '0' : date('Y-m-d', $member['memDeleted'] / 1000),
        ];
    }
}
