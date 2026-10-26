<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/4
 * Time: 20:23
 */

namespace App\Services;

use App\Common\Constant;
use App\Library\Uc\External;
use App\Models\MemberChangeRecord;
use App\Models\MemberUserRelation;

class HandoverService extends BaseService
{
    /**
     * 离职员工外部联系人再分配
     *
     * @param int $enterpriseId 企业id
     * @param string $qyDomain 企业domain
     * @param string $memUid 离职员工id
     * @param string $newMemUid 新导购id
     * @param string $ecUserIds 交接客户ec_userid集合
     *
     * @return bool
     * @throws \Exception
     */
    public function handoverCustomer($enterpriseId, $qyDomain, $memUid, $newMemUid, $ecUserIds)
    {
        // （离职员工）外部联系人标记关系表数据
        $relationList = MemberUserRelation::where(['mem_uid' => $memUid, 'enterprise_id' => $enterpriseId])
            ->whereIn('ec_userid', $ecUserIds)
            ->get()
            ->toArray();

        // （接替员工）已存在外部联系人标记关系表数据
        $existRelationList = MemberUserRelation::where(['mem_uid' => $newMemUid, 'enterprise_id' => $enterpriseId])
            ->whereIn('ec_userid', $ecUserIds)
            ->get()
            ->toArray();

        // （接替员工）已存在外部联系人id集合
        $existEcUserIds = array_column($existRelationList, 'ec_userid');

        // （接替员工）已存在的交接记录
        $changeRecords = MemberChangeRecord::where([
            'enterprise_id' => $enterpriseId,
            'mem_uid'       => $newMemUid
        ])->whereIn('ec_userid', $ecUserIds)->get()->toArray();

        // （接替员工）已存在的交接记录外部联系人id集合
        $existChangeEcUserIds = array_column($changeRecords, 'ec_userid');

        // 待新增外部联系人再分配记录
        $newChangeRecords = [];
        // 需要更新专属导购的外部联系人标记关系集合
        $changeExclusiveRelations = [];
        // 待删除的（离职员工）的外部联系人关系集合
        $delRelations = [];
        // （离职员工）需交接外部联系人标记关系集合
        $handoverRelations = [];

        foreach ($relationList as $relation) {
            if (!in_array($relation['ec_userid'], $existChangeEcUserIds)) {
                // 外部联系人再分配记录
                $newChangeRecords[] = [
                    'enterprise_id' => $enterpriseId,
                    'old_mem_uid'   => $memUid,
                    'ec_userid'     => $relation['ec_userid'],
                    'mem_uid'       => $newMemUid,
                    'created_at'    => time()
                ];
            }

            // 专属导购
            if (Constant::MEMBER_IS_EXCLUSIVE_TRUE === intval($relation['is_exclusive'])) {
                $changeExclusiveRelations[] = $relation;
            }

            // 已存在->离职员工关系删除
            if ($existEcUserIds && in_array($relation['ec_userid'], $existEcUserIds)) {
                $delRelations[] = $relation;
            } else { // 不存在->离职员工关系交接
                $handoverRelations[] = $relation;
            }
        }

        // 事务处理
        $this->transaction(function () use (
            $enterpriseId,
            $qyDomain,
            $memUid,
            $newMemUid,
            $newChangeRecords,
            $delRelations,
            $handoverRelations,
            $changeExclusiveRelations
        ) {
            // 批量新增外部联系人再分配记录
            if ($newChangeRecords) {
                MemberChangeRecord::insert($newChangeRecords);
            }

            if ($handoverRelations) {
                // 待交接外部联系人关系id集合
                $relationIds = array_column($handoverRelations, 'id');

                // 外部联系人标记关系导购更新为新导购、非专属导购
                $data = [
                    'mem_uid'      => $newMemUid,
                    'stars'        => MemberUserRelation::STARS_FALSE, // 改为非星标用户
                    'is_exclusive' => Constant::MEMBER_IS_EXCLUSIVE_FALSE,
                    'relation_at'  => time()
                ];
                MemberUserRelation::whereIn('id', $relationIds)->update($data);

                // 待交接的外部联系人记录id集合
                $ecIds = array_column($handoverRelations, 'ec_id');
                // 循环调用uc离职成员的外部联系人再分配接口（uc不支持批量处理，支持后可优化此处代码）
                foreach ($ecIds as $ecId) {
                    $condition = [
                        'handoverMemUid' => $memUid,
                        'takeoverMemUid' => $newMemUid,
                        'ecId'           => $ecId
                    ];
                    (new External())->transfer($qyDomain, $condition);
                }
            }

            // 待删除的（离职员工）的外部联系人关系集合
            if ($delRelations) {
                // 删除（离职员工）的外部联系人关系
                $delRelationIds = array_column($delRelations, 'id');
                MemberUserRelation::whereIn('id', $delRelationIds)->delete();
            }

            if ($changeExclusiveRelations) {
                // 更新外部联系人的专属导购
                $this->changeExclusive($enterpriseId, $newMemUid, $changeExclusiveRelations);
            }
        });

        return true;
    }

    /**
     * 更新外部联系人专属导购
     * （客户分配给新的导购，该客户的专属导购根据标记最早时间顺延给下一个导购，如果没有顺延导购将新导购设置为专属导购）
     *
     * @param int $enterpriseId 企业id
     * @param string $newMemUid 新导购id
     * @param array $changeExclusiveRelations 需要更新专属导购的外部联系人标记关系集合
     *
     * @return bool
     */
    protected function changeExclusive($enterpriseId, $newMemUid, $changeExclusiveRelations)
    {
        // 需要更新专属导购的外部联系人id集合
        $changeExclusiveEcUserIds = array_column($changeExclusiveRelations, 'ec_userid');
        // 外部联系人id为键
        $changeExclusiveRelations = array_column($changeExclusiveRelations, null, 'ec_userid');

        // 除去新导购的外部联系人关系列表（顺延导购）
        $nextRelationList = MemberUserRelation::where(['enterprise_id' => $enterpriseId])
            ->where('mem_uid', '<>', $newMemUid)
            ->whereIn('ec_userid', $changeExclusiveEcUserIds)
            ->orderBy('relation_at', 'ASC')
            ->get()
            ->toArray();

        // 有顺延专属导购的外部联系人关系列表
        $nextGroupRelations = [];
        array_walk($nextRelationList, function ($value) use (&$nextGroupRelations) {
            // 以联系人id分组，取顺延导购关系数据
            if (!isset($nextGroupRelations[$value['ec_userid']])) {
                $nextGroupRelations[$value['ec_userid']] = $value;
            }
        });

        // 有顺延专属导购的外部联系人标记关系id集合
        $nextRelationIds = array_column($nextGroupRelations, 'id');
        if ($nextRelationIds) {
            // 更新顺延导购为专属导购
            MemberUserRelation::whereIn('id', $nextRelationIds)
                ->update(['is_exclusive' => Constant::MEMBER_IS_EXCLUSIVE_TRUE]);
        }

        // 有顺延专属导购的外部联系人id集合
        $nextEcUserIds = array_column($nextGroupRelations, 'ec_userid');
        // 无顺延专属导购的外部联系人id集合
        $newEcUserIds = array_diff($changeExclusiveEcUserIds, $nextEcUserIds);

        // 新导购更新为专属导购并更新创建时间为离职员工的外部联系人关系创建时间
        if ($newEcUserIds) {
            // 新导购为专属导购的外部联系人关系列表
            $newRelationList = MemberUserRelation::where(['enterprise_id' => $enterpriseId, 'mem_uid' => $newMemUid])
                ->whereIn('ec_userid', $newEcUserIds)
                ->get()
                ->toArray();

            $changeExclusiveRelations = array_column($changeExclusiveRelations, null, 'ec_userid');
            foreach ($newRelationList as $newRelation) {
                // 外部联系人id
                $updateEcUserId = $newRelation['ec_userid'];
                // 待更新数据
                $updateData = [
                    'is_exclusive' => Constant::MEMBER_IS_EXCLUSIVE_TRUE,
                    'relation_at'  => $changeExclusiveRelations[$updateEcUserId]['relation_at'] ?? 0
                ];

                // 更新新导购（接替员工）为新导购
                MemberUserRelation::where([
                    'mem_uid'   => $newMemUid,
                    'ec_userid' => $updateEcUserId
                ])->update($updateData);
            }
        }

        return true;
    }
}
