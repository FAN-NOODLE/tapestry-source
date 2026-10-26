<?php
/**
 * Created by PhpStorm.
 * User: wangle@vchangyi.com
 * Date: 2021/04/27
 */

namespace App\Services;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Models\Scrm\Members;
use App\Models\Scrm\WorkExternalUserTag;
use App\Models\ScrmEnterprise\EnterpriseUser;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Models\Scrm\WorkExternalUserRelation;
use App\Models\StoreWorker;
use App\Services\Crm\CrmService;

class CustomerScrmService
{

    /**
     * @desc 根据条件筛选顾客列表
     *
     * @param $param
     *
     * @return array
     * @date 2021-04-27
     * @throws \Exception
     * @author wangle@vchangyi.com
     */
    public static function getNewCustomerList($param): array
    {
        $storeWorkerInfo = StoreWorkerService::detail([
            'enterprise_id' => $param['enterprise_id'],
            'store_id'      => $param['store_id'] ?? 0,
            'mobile'        => $param['mobile'],
            'worker_status' => StoreWorker::IN_SERVICE
        ]);
        if (empty($storeWorkerInfo)) {
            throw new Error(600103, '店员信息不存在');
        }

        $memberInfo = self::checkMembers($param['enterprise_id'], $param['mobile']);

        $param['member_id'] = $memberInfo['id'];

        $userList = WorkExternalUserRelation::query()
            ->select('work_external_user_relation.id', 'work_external_user_relation.work_external_user_id',
                'work_external_user_relation.createtime', 'work_external_user_relation.add_way', 'work_external_user_relation.unionid',
                'work_external_user_relation.user_id as sa_user_id', 'work_external_user_relation.is_key_sa')
            ->where('work_external_user_relation.enterprise_id', $param['enterprise_id'])
            ->where('work_external_user_relation.member_id', $param['member_id'])
            // 外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中
            ->where('work_external_user_relation.relation_status', 0)
            // member_deleted_time 成员主动删除关系时间
//            ->where('work_external_user_relation.member_deleted_time', 0)
            // external_deleted_time 外部联系人删除关系时间
//            ->where('work_external_user_relation.external_deleted_time', 0)
            // 添加方式
            ->when(!empty($param['add_way']) && $param['add_way'] == Constant::ADD_WAY_SCAN_QR_CODE, function ($query) use ($param) {
                return $query->where('work_external_user_relation.add_way', $param['add_way']);
            })
            ->when(!empty($param['add_way']) && $param['add_way'] == Constant::ADD_WAY_OTHER, function ($query) use ($param) {
                return $query->where('work_external_user_relation.add_way', '<>', Constant::ADD_WAY_SCAN_QR_CODE);
            })
            // 导购添加外部联系人起止时间
            ->when(isset($param['join_time_start'], $param['join_time_end']), function ($query) use ($param) {

                return $query->whereBetween('work_external_user_relation.createtime', [$param['join_time_start'], $param['join_time_end']]);
            })
            // 如果根据标签筛选，子查询
            ->when(isset($param['tag_ids']) && !empty($param['tag_ids']), function ($query) use ($param) {

                return $query->whereIn('work_external_user_relation.id', function ($query) use ($param) {

                    return $query->select('work_external_user_tag.relation_id')
                        ->from('work_external_user_tag')
                        ->leftJoin('work_external_tag',
                            'work_external_user_tag.external_tag_id',
                            '=',
                            'work_external_tag.id')
                        ->when(!empty($param['tag_ids']), function ($query) use ($param) {
                            return $query->whereIn('work_external_tag.tag_id', $param['tag_ids']);
                        })
                        ->where('work_external_user_tag.deleted_at', 0)
                        ->where('work_external_tag.deleted_at', 0);
                });
            })
            ->orderBy('work_external_user_relation.createtime', 'desc')
            ->get();

        $userList = $userList->where("unionid", '!=', '')->toArray();

        $workExternalUserIds = array_column($userList, 'work_external_user_id');
        $userListArray       = array_column($userList, null, 'work_external_user_id');
        $relationIds         = array_column($userList, 'id');

        // 用户标签
        $userTags = WorkExternalUserTag::getUserTagsByRelationIds($relationIds);

        // 切库根据外部联系人user_id查询
        [$wxUserInfo, $count] = WxUserCorpExternal::getExternalUserByIds($param, $workExternalUserIds);

        // 组装数据
        $customerUserList = [];
        foreach ($wxUserInfo as $k => $v) {
//            if ($v['mobile'] == '') {
//                $v = self::updateUser($v, $param);
//            }


            if (isset($userListArray[$v['id']])) {
                $v['relation_id'] = $userListArray[$v['id']]['id'];
                $v['relation_at'] = $userListArray[$v['id']]['createtime'];
                $v['sa_user_id']  = $userListArray[$v['id']]['sa_user_id'];
                $v['is_key_sa']   = $userListArray[$v['id']]['is_key_sa'];
//                $v['add_way']     = $userListArray[$v['id']]['add_way'];
//                $v['add_way_str'] = Constant::ADD_WAY_TYPE[$v['add_way']];
                $v['add_way']     = $userListArray[$v['id']]['add_way'] == Constant::ADD_WAY_SCAN_QR_CODE ? Constant::ADD_WAY_SCAN_QR_CODE : Constant::ADD_WAY_OTHER;
                $v['add_way_str'] = $userListArray[$v['id']]['add_way'] == Constant::ADD_WAY_SCAN_QR_CODE ? '导购名片' : '其他方式';

            }

            foreach ($userTags as $kk => $vv) {
                if ($v['relation_id'] == $vv['relation_id']) {
                    $v['tags'][] = $vv;
                }
            }

            $customerUserList[] = $v;
        }

        array_multisort(array_column($customerUserList, 'relation_at'), SORT_DESC, $customerUserList);

        return [
            'customer_list'  => $customerUserList,
            'worker_qr_code' => $storeWorkerInfo['qr_code'],
            'total'          => $count,
        ];
    }


    /**
     * 外部联系人详情和标签
     *
     * @param $enterpriseId
     * @param $mobile
     * @param $relationId
     *
     * @return array
     */
    public static function userInfo($enterpriseId, $mobile, $relationId): array
    {

        $memberInfo = self::checkMembers($enterpriseId, $mobile);

        $memberId = $memberInfo['id'];

        $workExternalUserInfo = WorkExternalUserRelation::query()
            ->select('work_external_user_relation.id', 'work_external_user_relation.work_external_user_id', 'work_external_user_relation.createtime', 'work_external_user_relation.add_way')
            ->where('work_external_user_relation.enterprise_id', $enterpriseId)
            ->where('work_external_user_relation.member_id', $memberId)
            ->where('work_external_user_relation.id', $relationId)
            ->first()
            ->toArray();

        // 用户标签
        $userTags = WorkExternalUserTag::getUserTagsByRelationIds([$relationId]);

        // 切库获取外部联系人详情
        $wxUserInfo = WxUserCorpExternal::getExternalUserInfoById($enterpriseId, $workExternalUserInfo['work_external_user_id']);

        $wxUserInfo['relation_at'] = $workExternalUserInfo['createtime'];
//        $wxUserInfo['add_way']     = $workExternalUserInfo['add_way'];
//        $wxUserInfo['add_way_str'] = Constant::ADD_WAY_TYPE[$workExternalUserInfo['add_way']];
        $wxUserInfo['add_way']     = $workExternalUserInfo['add_way'] == Constant::ADD_WAY_SCAN_QR_CODE ? Constant::ADD_WAY_SCAN_QR_CODE : Constant::ADD_WAY_OTHER;
        $wxUserInfo['add_way_str'] = $workExternalUserInfo['add_way'] == Constant::ADD_WAY_SCAN_QR_CODE ? '导购名片' : '其他方式';

//        $wxUserInfo['mobile'] = '';

        $wxUserInfo['tags'] = $userTags;

        return $wxUserInfo;

    }


    /**
     * 根据导购手机号查询导购信息
     *
     * @param $enterpriseId
     * @param $mobile
     *
     * @return array
     */
    public static function checkMembers($enterpriseId, $mobile): array
    {
        $memberInfo = Members::query()
            ->where([
                ['enterprise_id', $enterpriseId],
                ['status', Members::MEMBERS_STATUS_ACTIVATED],
                ['mobile', $mobile]
            ])
            ->first();

        if (empty($memberInfo)) {
            throw new Error(600103, '店员信息不存在');
        }

        return $memberInfo->toArray();
    }

    public static function updateUser($v, $param)
    {
        $userObj = WxUserCorpExternal::getExternalUserInfo($param['enterprise_id'], $v['unionid']);

        if (!empty($userObj->union_code)) {
            return $v;
        }
        // 调用SDK方法
        $crmParams = [
            'enterprise_id' => $param['enterprise_id'],
            'union_id'      => $v['unionid'],
        ];

        $result = (new CrmService($param['enterprise_id']))->getMemberInfo($crmParams);

        if ($result['status'] == true) {
            $updateData = [
                'union_code'   => $result['membership'] ?? '',
                'birthday'     => $result['birthday'],
                'birth_year'   => $result['birth_year'],
                'birth_month'  => $result['birth_month'],
                'birth_day'    => $result['birth_day'],
                'card_level'   => $result['customer_level'],
                'clt_type'     => $result['client_type'],
                'amount_total' => $result['last_12m_purchase_amount'] * 100,
            ];

            WxUserCorpExternal::updateExternalUserInfoByUnionId($param['enterprise_id'], $v['unionid'], $updateData);

            $updateUserData = [
                'union_code' => $result['membership'] ?? '',
                'birthday'   => $result['birthday'],
            ];
            EnterpriseUser::updateEnterpriseUserById($param['enterprise_id'], $v['eu_id'], $updateUserData);

            $v = array_merge($v, $updateData);
        }

        return $v;
    }

}
