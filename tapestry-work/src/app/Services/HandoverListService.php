<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/5
 * Time: 11:23
 */

namespace App\Services;

use App\Common\Constant;
use App\Library\Uc\External;
use App\Library\Uc\Member;
use App\Models\MemberChangeRecord;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\User;

class HandoverListService
{
    /**
     * 获取离职员工未分配/已分配客户列表
     *
     * @param int $enterpriseId 企业id
     * @param string $qyDomain 企业domain
     * @param string $appIdentifier 企业应用标识
     * @param string $memUid 用户id
     * @param int $handoverStatus 分配状态（1：未分配，2：已分配）
     * @param int $page 当前页码
     * @param int $size 每页数据量，当size=0时不分页
     *
     * @return array
     * @throws \Exception
     */
    public function getHandoverList(
        $enterpriseId,
        $qyDomain,
        $appIdentifier,
        $memUid,
        $handoverStatus,
        $page = Constant::DEFAULT_PAGE,
        $size = 0
    )
    {
        // 总数据量
        $total = 0;
        $list  = [];
        // 未分配（member_user_relation）
        if (Constant::CUSTOMER_HANDOVER_STATUS_FALSE == $handoverStatus) {
            // 查询条件
            $where = [
                'mem_uid'       => $memUid,
                'enterprise_id' => $enterpriseId
            ];

            // 分页
            if ($size) {
                // 总数据量
                $total = MemberUserRelation::query()->where($where)->count();
                if ($total) {
                    // 列表数据
                    $list = MemberUserRelation::query()->where($where)
                        ->forpage($page, $size)
                        ->orderBy('id', 'DESC')
                        ->get()
                        ->toArray();
                }
            } else { // 不分页
                $list = MemberUserRelation::query()->where($where)->orderBy('id', 'DESC')->get()->toArray();
            }

            // 数据格式化
            $data = $this->formatMemberUserRelation($enterpriseId, $list);
        } else { // 已分配（member_change_record）
            // 查询条件
            $where = [
                'old_mem_uid'   => $memUid,
                'enterprise_id' => $enterpriseId
            ];

            // 分页
            if ($size) {
                // 总数据量
                $total = MemberChangeRecord::query()->where($where)->count();
                if ($total) {
                    // 列表数据
                    $list = MemberChangeRecord::query()->where($where)
                        ->forpage($page, $size)
                        ->orderBy('id', 'DESC')
                        ->get()
                        ->toArray();
                }
            } else { // 不分页
                $list = MemberChangeRecord::query()->where($where)->orderBy('id', 'DESC')->get()->toArray();
            }

            // 数据格式化
            $data = $this->formatMemberChangeRecord($enterpriseId, $qyDomain, $appIdentifier, $list, $size);
        }

        return [$total, $data];
    }

    /**
     * 格式化未分配客户列表
     *
     * @param int $enterpriseId 企业id
     * @param array $relationList 待格式化数据
     *
     * @return array
     * @throws \Exception
     */
    protected function formatMemberUserRelation($enterpriseId, $relationList)
    {
        // 待返数据
        $data = [];
        if (!$relationList) {
            return $data;
        }

        // 客户union_id集合
        $unionIds = array_filter(array_unique(array_column($relationList, 'union_id')));
        sort($unionIds);

        // 用户信息列表
        $userList = [];
        if ($unionIds) {
            // 调用sdk获取用户信息列表
            $params   = [
                'enterprise_id' => $enterpriseId,
                'union_ids'     => $unionIds
            ];
            $userList = (new User())->users($params);
            $userList = array_column($userList['data'], null, 'union_id');
        }

        foreach ($relationList as $relation) {
            // 用户信息
            $user = $userList[$relation['union_id']] ?? [];

            $data[] = [
                'ec_userid'    => $relation['ec_userid'],
                'avatar'       => $user['avatar'] ?? $relation['ec_avatar'],
                'nick_name'    => $user['nick_name'] ?? $relation['ec_name'],
                'remark_name'  => $user['real_name'] ??$relation['user_remark_name'],
                'created_at'   => intval($relation['created_at']),
                'mem_uid'      => '', // 未分配，不需要导购数据
                'mem_username' => '', // 未分配，不需要导购数据
                'status'       => 0
            ];
        }

        return $data;
    }

    /**
     * 格式化已分配客户列表
     *
     * @param int $enterpriseId 企业id
     * @param string $qyDomain 企业domain
     * @param string $appIdentifier 企业应用标识
     * @param array $recordList 待格式化数据
     * @param int $size 每页数据量，当size=0时不分页
     *
     * @return array
     * @throws \Exception
     */
    protected function formatMemberChangeRecord($enterpriseId, $qyDomain, $appIdentifier, $recordList, $size)
    {
        // 待返数据
        $data = [];
        if (!$recordList) {
            return $data;
        }

        // 外部联系人userid集合
        $ecUserIds    = array_filter(array_column($recordList, 'ec_userid'));
        $relationList = MemberUserRelation::query()->whereIn('ec_userid', $ecUserIds)->get()->toArray();
        $relationList = array_column($relationList, null, 'ec_userid');

        // 客户union_id集合
        $unionIds = array_filter(array_unique(array_column($relationList, 'union_id')));
        sort($unionIds);

        // 用户信息列表
        $userList = [];
        if ($unionIds) {
            // 调用sdk获取用户信息列表
            $params   = [
                'enterprise_id' => $enterpriseId,
                'union_ids'     => $unionIds
            ];
            $userList = (new User())->users($params);
            $userList = array_column($userList['data'], null, 'union_id');
        }

        // 新导购id集合
        $memUids = array_unique(array_filter(array_column($recordList, 'mem_uid')));

        // 新导购员工数据
        $members = [];
        // 手机端不分页，size=0不需要新导购数据
        if ($memUids && $size) {
            sort($memUids);

            // 调用uc sdk获取员工数据
            $members = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, ['memUids' => $memUids]);
            $members = array_column($members['list'], null, 'memUid');
        }

        sort($ecUserIds);
        // 调用uc接口更新外部联系人状态
        (new External())->import($qyDomain, ['ids' => $ecUserIds]);

        // 待更新为已确认状态的记录id集合
        $updateRecordIds = [];
        foreach ($recordList as $record) {
            // 外部联系人id
            $ecUserid = $record['ec_userid'];

            // 外部联系人标记关系
            $relation = $relationList[$ecUserid] ?? [];
            $unionId  = $relation['union_id'] ?? '';
            // 用户信息
            $user = $userList[$unionId] ?? [];

            // 客户确认状态：待确认
            $status = Constant::CUSTOMER_HANDOVER_ACK_STATUS_FALSE;

            // 客户确认状态：待确认
            if (Constant::CUSTOMER_HANDOVER_ACK_STATUS_FALSE == $record['status']) {
                // 调用uc接口可查询外部联系人
                $externals = (new External())->listForMember($qyDomain, ['memUid' => $record['mem_uid']])['list'];

                if ($externals) {
                    // 外部联系人列表，ecUserid为key
                    $externals = array_column($externals, null, 'ecUserid');

                    if (isset($externals[$ecUserid]) && isset($externals[$ecUserid]['followMembers'])) {
                        // 外部联系人-员工id集合
                        $followMembers = array_column($externals[$ecUserid]['followMembers'], null, 'memUid');
                        // 客户确认状态
                        $status = $followMembers[$record['mem_uid']]['emRelationStatus'] ?? Constant::CUSTOMER_HANDOVER_ACK_STATUS_FALSE;
                    }
                }

                // 待确认->已确认，将记录id加入到待更新为已确认状态的记录id集合中
                if (Constant::CUSTOMER_HANDOVER_ACK_STATUS_TRUE == $status) {
                    $updateRecordIds[] = $record['id'];
                }
            } else {
                // 客户确认状态：已确认
                $status = Constant::CUSTOMER_HANDOVER_ACK_STATUS_TRUE;
            }

            // 修复由于该字段没有值导致的报错，by 陈朔
            if (empty($record['user_remark_name'])) {
                $record['user_remark_name'] = '';
            }

            $data[] = [
                'ec_userid'    => $ecUserid,
                'avatar'       => $user['avatar'] ?? $relation['ec_avatar'],
                'nick_name'    => $user['nick_name'] ?? $relation['ec_name'],
                'remark_name'  => $user['real_name'] ?? $record['user_remark_name'],
                'created_at'   => $relation['created_at'] ?? 0,
                'mem_uid'      => $record['mem_uid'],
                'mem_username' => $members[$record['mem_uid']]['memUsername'] ?? '',
                'status'       => $status // 客户确认状态（1：已确认，2：待确认）
            ];
        }

        // 批量更新状态为已确认
        if ($updateRecordIds) {
            $updateData = ['status' => Constant::CUSTOMER_HANDOVER_ACK_STATUS_TRUE];
            MemberChangeRecord::whereIn('id', $updateRecordIds)->update($updateData);
        }

        return $data;
    }
}
