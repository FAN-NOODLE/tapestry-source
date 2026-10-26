<?php


namespace App\Services\DataStatistic;

use App\Library\Uc\Member;
use App\Models\Scrm\Departments;
use App\Models\Scrm\FissionAssistantRecords;
use App\Models\Scrm\FissionChatGroup;
use App\Models\Scrm\MemberDepartments;
use App\Models\Scrm\Members;
use App\Models\Scrm\WorkExternalUserRelation;
use App\Models\Scrm\WorkGroupChat;
use App\Models\Scrm\WorkGroupChatMember;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Services\BaseService;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\DB;

/**
 * 数据罗盘-社群数据
 * Class CommunityDataCompassService
 * @package App\Services\DataStatistic
 */
class CommunityDataCompassService extends BaseService
{

    const GROUP_DATA_FIELDS = [
        'group_count'                   => 0,
        'group_member_count'            => 0,
        'new_join_group_count'          => 0,
        'leave_group_count'             => 0,
        'join_group_rate'               => '0%',
        'leave_group_rate'              => '0%',
        'join_friend_count'             => '0%',
        'join_friend_rate'              => '0%',
        'wechat_binding_vip_count'      => '0%',
        'wechat_binding_vip_rate'       => '0%',
        'group_activity_count'          => 0,
        'join_group_activity_count'     => 0,
        'not_join_group_activity_count' => 0,
        'not_join_group_activity_rate'  => '0%',
    ];

    const PAGE_LIMIT = 10;

    /**
     * 社群数据统计
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $memberIds
     * @return array
     * @Date: 2021/6/29
     * @Time: 14:53
     */
    public function overview(
        int $enterpriseId,
        int $startStatTime,
        int $endStatTime,
        array $memberIds
    ): array
    {
        $group       = (new WorkGroupChat())->getGroup(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memberIds,
            '',
            ['id', 'owner_member_id']
        );
        $groupChatId = array_unique(array_column($group, 'id'));
        $memberIds   = array_unique(array_column($group, 'owner_member_id'));

        return $this->getGroupData(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupChatId,
            $group,
            $memberIds
        );
    }

    /**
     * 获取scrm数据库对应的导购信息
     * @param $enterpriseId
     * @param array $workerList
     * @param array $field
     * @param string $order
     * @param bool $isPage
     * @param int $size
     * @return array
     * @Date: 2021/6/24
     * @Time: 18:24
     */
    public function getScrmMember(int $enterpriseId, array $workerList = [], array $field = ['*'], string $order = '', bool $isPage = false, int $size = self::PAGE_LIMIT): array
    {
        $memUserId = array_filter(array_column($workerList, 'mem_userid'));

        $query = Members::query()
            ->select($field)
            ->where('enterprise_id', $enterpriseId)
            ->when(!empty($memUserId), function ($query) use ($memUserId) {
                return $query->whereIn('userid', $memUserId);
            })
            ->when(!empty($order), function ($query) use ($order) {
                return $query->orderBy($order, 'DESC');
            })
            ->where('deleted_at', 0);

        if ($isPage) {
            return $query->paginate($size)->toArray();
        }

        return $query->get()->toArray();
    }

    /**
     * 获取社群数据明细列表（console）(按社群)
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $memberIds
     * @param int $size
     * @param bool $isPage
     * @return array
     * @Date: 2021/6/28
     * @Time: 17:29
     */
    public function getCommunityStatisticsByGroup(
        int $enterpriseId,
        int $startStatTime,
        int $endStatTime,
        array $memberIds,
        int $size,
        bool $isPage
    ): array
    {
        $group = (new WorkGroupChat())->getGroup(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memberIds,
            'create_time',
            ['id', 'name', 'owner_member_id'],
            $size,
            $isPage
        );
        !isset($group['data']) && $group['data'] = $group;
        $groupChatId    = array_unique(array_column($group['data'], 'id'));
        $groupAndMember = array_column($group['data'], 'owner_member_id', 'id');

        if ($groupChatId) {
            // 累计群成员数相关数据(入群与退群)
            $groupMember = (new WorkGroupChatMember())->getGroupMember(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                [DB::raw('COUNT(userid) AS groupMemberCount'), 'group_chat_id'],
                'group_chat_id'
            );
            $groupMember = array_column($groupMember, 'groupMemberCount', 'group_chat_id');

            // 去重
            $groupMemberUnique = (new WorkGroupChatMember())->getGroupMember(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                [DB::raw('COUNT(DISTINCT userid) AS groupMemberCount'), 'group_chat_id'],
                'group_chat_id'
            );
            $groupMemberUnique = array_column($groupMemberUnique, 'groupMemberCount', 'group_chat_id');

            // 新入群人数相关数据
            $newJoinGroup = (new WorkGroupChatMember())->getJoinOrLeaveGroup(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                WorkGroupChatMember::JOIN_GROUP,
                [DB::raw('COUNT( userid ) AS newJoinGroupCount'), 'group_chat_id'],
                'group_chat_id'
            );
            $newJoinGroup = array_column($newJoinGroup, 'newJoinGroupCount', 'group_chat_id');

            // 退群人数相关数据
            $leaveGroup = (new WorkGroupChatMember())->getJoinOrLeaveGroup(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                WorkGroupChatMember::LEAVE_GROUP,
                [DB::raw('COUNT(userid) AS leaveGroupCount'), 'group_chat_id'],
                'group_chat_id'
            );
            $leaveGroup = array_column($leaveGroup, 'leaveGroupCount', 'group_chat_id');

            // 获取好友、微信会员相关统计
            list($workExternalUserMap, $externalUserInfo, $isVipExternalUser, $joinGroupsMap) = $this->getFriendAndWeChatVipData(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                $groupAndMember
            );

            // 获取活动相关数据
            list($groupActivityMap, $groupMemberMap, $groupActivityCountArray, $joinGroupActivity, $groupActivityCount) = $this->getActivityData(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId
            );
        }

        // 组装数据
        foreach ($group['data'] as &$value) {
            // 获取当前群的外部联系人
            $externalIds = isset($groupAndMember[$value['id']]) && isset($workExternalUserMap[$groupAndMember[$value['id']]]) ? $workExternalUserMap[$groupAndMember[$value['id']]] : [];

            // 加导购的外部联系人
            $unionIdJoinMemberArray = [];
            foreach ($externalIds as $item) {
                if (isset($externalUserInfo[$item]) && !empty($externalUserInfo[$item])) {
                    $unionIdJoinMemberArray[] = $externalUserInfo[$item];
                }
            }

            // 入群的外部联系人
            $unionIdByJoinGroupArray = isset($joinGroupsMap) && isset($joinGroupsMap[$value['id']]) ? $joinGroupsMap[$value['id']] : [];

            $groupMemberCount         = isset($groupMember) && isset($groupMember[$value['id']]) ? $groupMember[$value['id']] : 0;
            $newJoinGroupCount        = isset($newJoinGroup) && isset($newJoinGroup[$value['id']]) ? $newJoinGroup[$value['id']] : 0;
            $leaveGroupCount          = isset($leaveGroup) && isset($leaveGroup[$value['id']]) ? $leaveGroup[$value['id']] : 0;
            $groupActivityCount       = isset($groupActivityCountArray) && isset($groupActivityCountArray[$value['id']]) ? $groupActivityCountArray[$value['id']] : 0;
            $groupActivityMapExternal = isset($groupActivityMap) && isset($groupActivityMap[$value['id']]) ? $groupActivityMap[$value['id']] : [];
            $groupMemberMapExternal   = isset($groupMemberMap) && isset($groupMemberMap[$value['id']]) ? $groupMemberMap[$value['id']] : [];
            $joinFriendCount          = count(array_intersect(array_unique($unionIdJoinMemberArray), array_unique($unionIdByJoinGroupArray)));
//            $wechatBindingVipCount    = isset($isVipExternalUser) ? count(array_intersect($isVipExternalUser, $unionIdByJoinGroupArray)) : 0;

            $groupMemberUniqueCount = isset($groupMemberUnique) && isset($groupMemberUnique[$value['id']]) ? $groupMemberUnique[$value['id']] : 0;

            $value['group_member_count']       = $groupMemberCount;
            $value['new_join_group_count']     = $newJoinGroupCount;
            $value['leave_group_count']        = $leaveGroupCount;
            $value['join_group_rate']          = $this->getDataRate($newJoinGroupCount, $groupMemberCount);
            $value['leave_group_rate']         = $this->getDataRate($leaveGroupCount, $groupMemberCount);
            $value['join_friend_count']        = $joinFriendCount;
            $value['join_friend_rate']         = $this->getDataRate($joinFriendCount, $groupMemberCount);
//            $value['wechat_binding_vip_count'] = $wechatBindingVipCount;
//            $value['wechat_binding_vip_rate']  = $this->getDataRate($wechatBindingVipCount, $groupMemberCount);
            $value['group_activity_count']     = $groupActivityCount;
            // 记录表与群成员数据求交集结果就是当前群成员参与活动的数量
            $value['join_group_activity_count'] = count(array_intersect($groupActivityMapExternal, $groupMemberMapExternal));
            // 求差集
            $value['not_join_group_activity_count'] = count(array_unique(array_diff(array_unique($groupMemberMapExternal), array_unique($groupActivityMapExternal))));
            $value['not_join_group_activity_rate']  = $this->getDataRate($value['not_join_group_activity_count'], $groupMemberUniqueCount);
        }

        return $group;
    }

    /**
     * 获取未参加活动人员列表
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $memberIds
     * @param int $size
     * @param int $page
     * @return array
     * @Date: 2021/6/29
     * @Time: 14:51
     */
    public function getNotJoinActivityMemberList(
        int $enterpriseId,
        int $startStatTime,
        int $endStatTime,
        array $memberIds,
        int $size,
        int $page
    ): array
    {
        // 获取群组
        $group           = (new WorkGroupChat())->getGroup(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memberIds,
            '',
            ['id', 'name']
        );
        $groupChatId     = array_unique(array_column($group, 'id'));
        $groupChatIdName = array_column($group, 'name', 'id');

        $groupActivity    = (new FissionChatGroup())->getGroupActivity($enterpriseId, $startStatTime, $endStatTime, $groupChatId);
        $groupActivityIds = array_column($groupActivity, 'activity_id');

        // 记录表里查询出活动的数据
        $joinGroupActivity = (new FissionAssistantRecords())->getJoinGroupActivity(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupActivityIds,
            ['external_user_id']
        );
        $joinGroupActivity = array_column($joinGroupActivity, 'external_user_id');

        // 查出所有群的人员(仅入群)
        $groupMember          = (new WorkGroupChatMember())->getGroupMember(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupChatId,
            ['userid', 'group_chat_id', 'join_time'],
            '',
            WorkGroupChatMember::JOIN_GROUP
        );
        $groupMemberUserIds   = array_column($groupMember, 'userid');
        $groupMemberMap       = array_column($groupMember, 'group_chat_id', 'userid');
        $groupMemberJoinTimes = array_column($groupMember, 'join_time', 'userid');

        // 求上述两者之间的差集就是（未参与活动人数）
        $notJoinActivityMembers = array_unique(array_diff($groupMemberUserIds, $joinGroupActivity));

        $notJoinActivityMembersList = (new WxUserCorpExternal())->getExternalUserInfoByOpenId(
            $enterpriseId,
            $notJoinActivityMembers,
            ['id', 'nick', 'name', 'avatar', 'openid']
        );

        $sortColumn = [];
        foreach ($notJoinActivityMembersList as $key => &$value) {
            $chatId = isset($groupMemberMap[$value['openid']]) ? $groupMemberMap[$value['openid']] : '';

            $value['group_name'] = isset($groupChatIdName[$chatId]) && $groupChatIdName[$chatId] ? $groupChatIdName[$chatId] : '未命名';
            $value['join_time']  = $groupMemberJoinTimes[$value['openid']];
            $sortColumn[$key]    = $groupMemberJoinTimes[$value['openid']];
        }

        array_multisort($sortColumn, SORT_DESC, $notJoinActivityMembersList);

        $size = $size > 0 ? $size : self::PAGE_LIMIT;
        $page = $page > 0 ? $page : 1;

        $result = array_chunk($notJoinActivityMembersList, $size)[$page - 1] ?? [];

        return [
            'total' => $notJoinActivityMembersList ? count($notJoinActivityMembersList) : 0,
            'page'  => $page,
            'size'  => $size,
            'list'  => $result
        ];
    }

    /**
     * 获取导购社群数据明细列表（console）(按导购)
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $workerList
     * @param int $size
     * @param bool $isPage
     * @return array
     * @Date: 2021/7/5
     * @Time: 15:47
     */
    public function getCommunityStatisticsByMember(
        int $enterpriseId,
        int $startStatTime,
        int $endStatTime,
        array $workerList,
        int $size,
        bool $isPage
    ): array
    {
        // 获取scrm数据库中member数据表的相关导购数据id
        $member = $this->getScrmMember(
            $enterpriseId,
            $workerList,
            ['id', 'userid'],
            'id',
            $isPage,
            $size
        );
        !isset($member['data']) && $member['data'] = $member;
        $memberIds = $member['data'];

        // 导购ID、头像、昵称
        $workerLists = array_column($workerList, null, 'mem_userid');

        // 导购所在门店
        $memberDepartments = MemberDepartments::query()
            ->whereIn('member_id', array_column($memberIds, 'id'))
            ->where('enterprise_id', $enterpriseId)
            ->get(['member_id', 'department_id'])
            ->toArray();
        $memberDepartments = array_column($memberDepartments, 'department_id', 'member_id');
        $allDepartments    = array_unique(array_values($memberDepartments));

        $departments = Departments::query()
            ->whereIn('id', $allDepartments)
            ->get(['id', 'name'])
            ->toArray();
        $departments = array_column($departments, 'name', 'id');

        // 获取群
        $group    = (new WorkGroupChat())->getGroup(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memberIds,
            'create_time',
            ['id', 'owner_member_id']
        );
        $groupMap = $this->getAssemblyData($group, 'owner_member_id', 'id');

        $groupDataByMember = [];
        if ($groupMap) {
            foreach ($groupMap as $memberId => $groupIds) {
                $groupDataByMember[$memberId] = $this->getGroupData(
                    $enterpriseId,
                    $startStatTime,
                    $endStatTime,
                    $groupIds,
                    $group,
                    [$memberId]
                );
            }
        }

        // 组装数据
        foreach ($member['data'] as &$value) {
            $groupData = isset($groupDataByMember[$value['id']]) ? $groupDataByMember[$value['id']] : [];

            $value['worker_no']                     = isset($workerLists[$value['userid']]) ? $workerLists[$value['userid']]['worker_no'] : '';
            $value['nickname']                      = isset($workerLists[$value['userid']]) ? $workerLists[$value['userid']]['nickname'] : '';
            $value['avatar']                        = isset($workerLists[$value['userid']]) ? $workerLists[$value['userid']]['avatar'] : '';
            $value['department_name']               = isset($memberDepartments[$value['id']]) && isset($departments[$memberDepartments[$value['id']]]) ? $departments[$memberDepartments[$value['id']]] : '';
            $value['group_member_count']            = isset($groupData['group_member_count']) ? $groupData['group_member_count'] : 0;
            $value['new_join_group_count']          = isset($groupData['new_join_group_count']) ? $groupData['new_join_group_count'] : 0;
            $value['leave_group_count']             = isset($groupData['leave_group_count']) ? $groupData['leave_group_count'] : 0;
            $value['join_group_rate']               = isset($groupData['join_group_rate']) ? $groupData['join_group_rate'] : '0%';
            $value['leave_group_rate']              = isset($groupData['leave_group_rate']) ? $groupData['leave_group_rate'] : '0%';
            $value['join_friend_count']             = isset($groupData['join_friend_count']) ? $groupData['join_friend_count'] : 0;
            $value['join_friend_rate']              = isset($groupData['join_friend_rate']) ? $groupData['join_friend_rate'] : '0%';
            $value['wechat_binding_vip_count']      = isset($groupData['wechat_binding_vip_count']) ? $groupData['wechat_binding_vip_count'] : 0;
            $value['wechat_binding_vip_rate']       = isset($groupData['wechat_binding_vip_rate']) ? $groupData['wechat_binding_vip_rate'] : '0%';
            $value['group_activity_count']          = isset($groupData['group_activity_count']) ? $groupData['group_activity_count'] : 0;
            $value['join_group_activity_count']     = isset($groupData['join_group_activity_count']) ? $groupData['join_group_activity_count'] : 0;
            $value['not_join_group_activity_count'] = isset($groupData['not_join_group_activity_count']) ? $groupData['not_join_group_activity_count'] : 0;
            $value['not_join_group_activity_rate']  = isset($groupData['not_join_group_activity_rate']) ? $groupData['not_join_group_activity_rate'] : '0%';

            unset($value['id'], $value['userid']);
        }

        return $member;
    }

    /**
     * 获取活动相关数据
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $groupChatId
     * @return array
     * @Date: 2021/6/29
     * @Time: 14:52
     */
    private function getActivityData(
        int $enterpriseId,
        int $startStatTime,
        int $endStatTime,
        array $groupChatId
    ): array
    {
        $groupActivity    = (new FissionChatGroup())->getGroupActivity(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupChatId
        );
        $groupActivityIds = array_unique(array_column($groupActivity, 'activity_id'));

        $groupActivity = collect($groupActivity)->groupBy('chat_id')->toArray();

        // 群活动数相关数据(一个群对应一个活动)
        $groupActivityCountArray = [];
        $groupActivityMap        = [];
        foreach ($groupActivity as $chatId => $value) {
            $groupActivityCountArray[$chatId] = count(array_column($value, 'activity_id'));
            $groupActivityMap[$chatId]        = $value[0]['activity_id'];
        }

        // 测试数据
        /*$groupActivityMap = [
            3326625290047360 => 3322942304384896,
            3325400319231872 => 3322942304384896,
            3324187247298432 => 3322942304384896,
            3322250348507008 => 3326799511363456,
            3322250142791552 => 3326799511363456
        ];
        $groupActivityIds = [3326799511363456, 3322942304384896];*/

        // 获取活动相关记录数据
        $groupActivityIds && $joinGroupActivity = (new FissionAssistantRecords())->getJoinGroupActivity(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupActivityIds,
            ['external_user_id', 'activity_id']
        );

        // 活动id与外部联系人关系数据
        isset($joinGroupActivity) && $joinGroupActivityRecords = $this->getAssemblyData($joinGroupActivity, 'activity_id', 'external_user_id');

        // 获取群id对应的外部联系人相关数据（记录表里的）
        foreach ($groupActivityMap as $key => $value) {
            $groupActivityMap[$key] = isset($joinGroupActivityRecords) && isset($joinGroupActivityRecords[$value]) ? $joinGroupActivityRecords[$value] : [];
        }

        // 获取群成员数据（仅入群）
        $groupMember = (new WorkGroupChatMember())->getGroupMember(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupChatId,
            ['group_chat_id', 'userid'],
            '',
            WorkGroupChatMember::JOIN_GROUP
        );

        // 获取群id与外部联系人关联数据
        $groupMemberMap = $this->getAssemblyData($groupMember, 'group_chat_id', 'userid');

        return [
            $groupActivityMap,
            $groupMemberMap,
            $groupActivityCountArray,
            isset($joinGroupActivity) ? $joinGroupActivity : [],
            count($groupActivityIds)
        ];
    }

    /**
     * 根据群获取统计数据
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $groupChatId
     * @param array $group
     * @param array $memberIds
     * @return array
     * @Date: 2021/7/5
     * @Time: 15:50
     */
    private function getGroupData(
        int $enterpriseId,
        int $startStatTime,
        int $endStatTime,
        array $groupChatId,
        array $group,
        array $memberIds
    ): array
    {
        if ($groupChatId) {
            // 累计群成员数(入群与退群)(不去重)
            $groupMember      = (new WorkGroupChatMember())->getGroupMember(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                ['userid', 'status', 'group_chat_id']
            );
            $groupMemberCount = count($groupMember);

            // 去重后群成员累计数
            $groupUniqueMemberCount = count(array_unique(array_column($groupMember, 'userid')));

            // 新入群人数
            $newJoinGroup = (new WorkGroupChatMember())->getJoinOrLeaveGroup(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                WorkGroupChatMember::JOIN_GROUP,
                ['id', 'group_chat_id', 'unionid']
            );

            $newJoinGroupCount = count($newJoinGroup);

            // 退群人数
            $leaveGroup      = (new WorkGroupChatMember())->getJoinOrLeaveGroup(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId,
                WorkGroupChatMember::LEAVE_GROUP,
                ['id', 'group_chat_id', 'userid']
            );
            $leaveGroupCount = count($leaveGroup);

            // 客户入群率
            $joinGroupRate = $this->getDataRate($newJoinGroupCount, $groupMemberCount);

            // 客户流失率
            $leaveGroupRate = $this->getDataRate($leaveGroupCount, $groupMemberCount);

            // 获取好友、微信会员相关统计-A
            $allGroupJoinMembers = array_unique(array_column($newJoinGroup, 'unionid'));

            // A-1、获取当前所有导购对应的所有外部联系人关系数据
            $workExternalUser = (new WorkExternalUserRelation())->getWorkExternalUserRelation(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $memberIds,
                ['work_external_user_id', 'member_id']
            );

            // A-2、获取外部联系人信息
            $workExternalUser && $externalUserInfo = (new WxUserCorpExternal())->getExternalUserInfoByIds(
                $enterpriseId,
                array_unique(array_column($workExternalUser, 'work_external_user_id')),
                ['unionid', 'id']
            );
            $externalUserInfo = isset($externalUserInfo) ? array_column($externalUserInfo, 'unionid') : [];

            // A-3、求交集获取已入群且添加导购好友数
            $joinFriendCount = count(array_intersect($externalUserInfo, $allGroupJoinMembers));

            // A-4、获取微信会员数
            $allGroupJoinMembers && $isVipExternalUser = (new WxUserCorpExternal())->getExternalUserInfoByUnionId(
                $enterpriseId,
                $allGroupJoinMembers,
                ['unionid']
            );
            $isVipExternalUser = isset($isVipExternalUser) ? array_column($isVipExternalUser, 'unionid') : [];

            // 获取活动相关数据
            list($groupActivityMap, $groupMemberMap, $groupActivityCountArray, $joinGroupActivity, $groupActivityCount) = $this->getActivityData(
                $enterpriseId,
                $startStatTime,
                $endStatTime,
                $groupChatId
            );
        }

        // 去重后群成员累计数
        $groupUniqueMemberCount = isset($groupUniqueMemberCount) ? $groupUniqueMemberCount : 0;

        // 群活动参与人数
        $joinGroupActivityCount = 0;
        foreach ($groupChatId as $chatId) {
            $groupActivityMemberIds = isset($groupActivityMap) && isset($groupActivityMap[$chatId]) ? $groupActivityMap[$chatId] : [];
            $groupMemberIds         = isset($groupMemberMap) && isset($groupMemberMap[$chatId]) ? $groupMemberMap[$chatId] : [];

            $joinGroupActivityCount += count(array_intersect($groupActivityMemberIds, $groupMemberIds));
        }

        $joinGroupActivity = isset($joinGroupActivity) ? array_column($joinGroupActivity, 'external_user_id') : [];
        $groupJoinMember   = isset($groupMember) ? array_filter($groupMember, function ($value) {
            return $value['status'] == WorkGroupChatMember::JOIN_GROUP;
        }) : [];
        $groupJoinMember   = isset($groupJoinMember) ? array_column($groupJoinMember, 'userid') : [];

        // 未参加群活动人数
        $notJoinGroupActivityCount = count(array_unique(array_diff(array_unique($groupJoinMember), array_unique($joinGroupActivity))));

        // 未参加群活动率
        $notJoinGroupActivityRate = $this->getDataRate($notJoinGroupActivityCount, $groupUniqueMemberCount);

        $joinFriendCount        = isset($joinFriendCount) ? $joinFriendCount : 0;
        $isVipExternalUserCount = isset($isVipExternalUser) ? count($isVipExternalUser) : 0;

        return [
                'group_count'                   => count($group),
                'group_member_count'            => isset($groupMemberCount) ? $groupMemberCount : 0,
                'new_join_group_count'          => isset($newJoinGroupCount) ? $newJoinGroupCount : 0,
                'leave_group_count'             => isset($leaveGroupCount) ? $leaveGroupCount : 0,
                'join_group_rate'               => isset($joinGroupRate) ? $joinGroupRate : '0%',
                'leave_group_rate'              => isset($leaveGroupRate) ? $leaveGroupRate : '0%',
                'join_friend_count'             => $joinFriendCount,
                'join_friend_rate'              => isset($groupMemberCount) ? $this->getDataRate($joinFriendCount, $groupUniqueMemberCount) : '0%',
                'wechat_binding_vip_count'      => $isVipExternalUserCount,
                'wechat_binding_vip_rate'       => isset($groupMemberCount) ? $this->getDataRate($isVipExternalUserCount, $groupMemberCount) : '0%',
                'group_activity_count'          => isset($groupActivityCount) ? $groupActivityCount : 0,
                'join_group_activity_count'     => $joinGroupActivityCount,
                'not_join_group_activity_count' => $notJoinGroupActivityCount,
                'not_join_group_activity_rate'  => $notJoinGroupActivityRate,
            ] + self::GROUP_DATA_FIELDS;
    }

    /**
     * 获取加好友、微信会员相关统计
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $groupChatId
     * @param array $groupAndMember
     * @return array
     * @Date: 2021/7/6
     * @Time: 11:15
     */
    private function getFriendAndWeChatVipData(int $enterpriseId, int $startStatTime, int $endStatTime, array $groupChatId, array $groupAndMember): array
    {
        // 加群数据
        $joinGroups = (new WorkGroupChatMember())->getJoinOrLeaveGroup(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $groupChatId,
            WorkGroupChatMember::JOIN_GROUP,
            ['group_chat_id', 'unionid']
        );

        // 获取群id与外部联系人关联数据
        $joinGroupsMap = $this->getAssemblyData($joinGroups, 'group_chat_id', 'unionid');

        // 获取所有创建群组的导购的外部联系人关系数据
        $groupAndMember && $workExternalUser = (new WorkExternalUserRelation())->getWorkExternalUserRelation(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            array_values($groupAndMember),
            ['work_external_user_id', 'member_id']
        );
        $workExternalUser    = isset($workExternalUser) ? $workExternalUser : [];
        $workExternalUserMap = $this->getAssemblyData($workExternalUser, 'member_id', 'work_external_user_id');

        // 从外部联系人数据表获取数据
        $workExternalUser && $externalUserInfo = (new WxUserCorpExternal())->getExternalUserInfoByIds(
            $enterpriseId,
            array_unique(array_column($workExternalUser, 'work_external_user_id')),
            ['unionid', 'id']
        );
        $externalUserInfo = isset($externalUserInfo) ? array_column($externalUserInfo, 'unionid', 'id') : [];

        // 根据unionid查询是会员信息数据
        $joinGroups && $isVipExternalUser = (new WxUserCorpExternal())->getExternalUserInfoByUnionId(
            $enterpriseId,
            array_column($joinGroups, 'unionid'),
            ['unionid']
        );
        $isVipExternalUser = isset($isVipExternalUser) ? $isVipExternalUser : [];
        $isVipExternalUser = array_column($isVipExternalUser, 'unionid');

        return [
            $workExternalUserMap,
            $externalUserInfo,
            $isVipExternalUser,
            $joinGroupsMap
        ];
    }

    /**
     * 根据部门（所谓的门店）与导购求交集(console)
     * @param int $enterpriseId
     * @param string $qyDomain
     * @param array $params
     * @return array
     * @throws \Exception
     * @Date: 2021/7/12
     * @Time: 14:49
     */
    public function getDepartmentAndMemberIntersection(int $enterpriseId, string $qyDomain, array $params): array
    {
        $guideList = [];

        if (isset($params['department_ids']) && !empty($params['department_ids'])) {
            $result = (new Member())->ignorePermissionsList($qyDomain, '', [
                'dpIdList' => $params['department_ids'],
                'pageNum'  => 1,
                'pageSize' => 9999,
            ]);
            $memUidList               = array_column($result['list'], 'memUid');
            $guideListByDepartmentIds = StoreWorkerService::list(['enterprise_id' => $enterpriseId, 'mem_uids' => $memUidList], false);
            $guideList                = $guideListByDepartmentIds['list'];
        }

        if (isset($params['worker_ids']) && !empty($params['worker_ids'])) {
            $guideListByWorkerIds = StoreWorkerService::list(['ids' => $params['worker_ids'], 'enterprise_id' => $enterpriseId], false);
            $guideList            = $guideListByWorkerIds['list'];
        }

        if (isset($params['department_ids']) && !empty($params['department_ids']) && isset($params['worker_ids']) && !empty($params['worker_ids'])) {
            $guideListByDepartmentIds = isset($guideListByDepartmentIds) && isset($guideListByDepartmentIds['list']) ? array_column($guideListByDepartmentIds['list'], null, 'id') : [];
            $guideListByWorkerIds     = isset($guideListByWorkerIds) && isset($guideListByWorkerIds['list']) ? array_column($guideListByWorkerIds['list'], null, 'id') : [];
            $guideList                = array_intersect_key($guideListByDepartmentIds, $guideListByWorkerIds);
        }
        return array_values($guideList);
    }

    /**
     * 计算比率
     * @param int $molecule
     * @param int $denominator
     * @param bool $isPercentage
     * @return string
     * @Date: 2021/7/7
     * @Time: 11:22
     */
    public function getDataRate(int $molecule, int $denominator, bool $isPercentage = false)
    {
        $rate = 0;
        if ($denominator > 0) {
            $rate = $molecule / $denominator * 100;
        }

        if ($isPercentage) {
            return $rate;
        }

        if ($rate > 0) {
            return round($rate, 2) . '%';
        }

        return '0%';
    }

    /**
     * 数据组装
     * @param array $data
     * @param string $kKey
     * @param string $vKey
     * @return array
     * @Date: 2021/6/29
     * @Time: 14:52
     */
    private function getAssemblyData(array $data, string $kKey, string $vKey): array
    {
        $result = [];

        foreach ($data as $key => $value) {
            $result[$value[$kKey]][] = $value[$vKey];
        }

        return $result;
    }
}