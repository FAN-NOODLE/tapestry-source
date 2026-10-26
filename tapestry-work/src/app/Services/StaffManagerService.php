<?php


namespace App\Services;


use App\Exceptions\Error;
use App\Library\Uc\Department;
use App\Library\Uc\Member;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use App\Models\MemberGuide;
use App\Models\StaffRight;
use App\Models\StaffSetting;
use App\Models\SyncOrgLog;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\StaffRight as StaffRightSDK;
use Illuminate\Auth\Events\Login;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use tests\TransportUdpTest;

class StaffManagerService extends BaseService
{
    /**
     * 默认返回的uc员工数量
     */
    const PAGESIZE = 1000;

    /**
     * 插入允许参加内购的员工
     *
     * @param $input
     *
     * @return bool
     */
    public function createForStaffRight($input)
    {
        $enterpriseId   = $input['enterprise_id'] ?? 0;
        $staffRightList = $input['staff_buy_list'];

        if (empty($enterpriseId)) {
            notice(600000);
        }

        // 插入数据
        foreach ($staffRightList as &$value) {
            $value['object_extra_id'] = $value['obj_id'];
            if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                $insertDepartmentMember   = (new Member())->ignorePermissionsList($input['qy_domain'], $input['app_identifier'], ['dpId' => $value['obj_id'], 'departmentChildrenFlag' => 1, 'memSubscribeStatus' => 1]);
                $value['object_extra_id'] = implode(',', array_column($insertDepartmentMember['list'], 'memUid'));
            }
            $value['enterprise_id'] = $enterpriseId;
            $value['created_at']    = time();
            $value['updated_at']    = time();
        }

        return StaffRight::query()->insert($staffRightList);
    }

    /**
     * @param $input
     *
     * @return array
     */
    public function updateForStaffRight($input)
    {
        $enterpriseId   = $input['enterprise_id'] ?? 0;
        $staffRightList = $input['staff_buy_list'];

        if (empty($enterpriseId)) {
            notice(600000);
        }

        if ($staffRightList) {
            // 员工内购权限部门
            $staffRightDepartment = StaffRight::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('obj_type', StaffRight::TYPE_DEP)
                ->get()
                ->toArray();

            // 员工内购权限员工
            $staffRightMember = StaffRight::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('obj_type', StaffRight::TYPE_STAFF)
                ->get()
                ->toArray();

            $staffRightDepartmentIds = array_column($staffRightDepartment, 'obj_id');
            $staffRightMemberIds     = array_column($staffRightMember, 'obj_id');
            $staffRightMember        = array_column($staffRightMember, null, 'obj_id');

            $notDeleteDepartmentIds = [];
            $notDeleteMemberIds     = [];

            foreach ($staffRightList as $key => $value) {
                $staffRightList[$key]['enterprise_id'] = $enterpriseId;
                $staffRightList[$key]['created_at']    = time();
                $staffRightList[$key]['updated_at']    = time();
                if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                    if (in_array($value['obj_id'], $staffRightDepartmentIds)) {
                        $notDeleteDepartmentIds[] = $value['obj_id'];
                        unset($staffRightList[$key]);
                    }
                }
                if ($value['obj_type'] == StaffRight::TYPE_STAFF) {
                    if (in_array($value['obj_id'], $staffRightMemberIds) && $value['obj_parent_id'] == $staffRightMember[$value['obj_id']]['obj_parent_id']) {
                        $notDeleteMemberIds[] = $value['obj_id'];
                        unset($staffRightList[$key]);
                    }
                }
            }

            if ($notDeleteDepartmentIds) {
                StaffRight::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('obj_type', StaffRight::TYPE_DEP)
                    ->whereNotIn('obj_id', $notDeleteDepartmentIds)
                    ->delete();
            }

            if ($notDeleteMemberIds) {
                StaffRight::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('obj_type', StaffRight::TYPE_STAFF)
                    ->whereNotIn('obj_id', $notDeleteMemberIds)
                    ->delete();
            }

            if ($staffRightList) {
                StaffRight::query()->insert($staffRightList);
            }
        } else {
            StaffRight::query()
                ->where('enterprise_id', $enterpriseId)
                ->delete();
        }

        return [];
    }

    /**
     * 插入或更新内购员工设置
     *
     * @param $input
     *
     * @return \Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|int
     */
    public function createOrUpdateForStaffSetting($input)
    {
        $enterpriseId = $input['enterprise_id'] ?? 0;
        $setting      = $input['setting'];
        unset($setting['enterprise_id']);

        if ($setting['binding_type'] == StaffSetting::FOREVER) {
            if ($setting['change_owner'] == StaffSetting::ALLOW) {
                notice(720010);
            }
        }

        if ($setting['binding_type'] == StaffSetting::EXPIRE) {
            if (!isset($setting['duration_time']) || !is_int($setting['duration_time'])
                || $setting['duration_time'] < 0 || $setting['duration_time'] > 365) {
                notice(720008);
            }
        }

        if ($setting['change_owner'] == StaffSetting::ALLOW) {
            if (!isset($setting['protected_time']) || !is_int($setting['protected_time']) || $setting['protected_time'] < 0) {
                notice(720009);
            }
            if ($setting['binding_type'] == StaffSetting::EXPIRE && $setting['protected_time'] > $setting['duration_time']) {
                notice(720009);
            }
        }

        if ((empty(
                $setting['duration_time']) && $setting['duration_time'] != 0) || $setting['duration_time'] == 'null'
            || empty($setting['duration_time'])
        ) {
            $setting['duration_time'] = 0;
        }

        if (
            (empty($setting['protected_time']) && $setting['protected_time'] != 0) || $setting['protected_time'] == 'null'
            || empty($setting['protected_time'])
        ) {
            $setting['protected_time'] = 0;
        }

        $staffSetting = StaffSetting::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'type'          => StaffSetting::STAFF_TYPE,
            ])
            ->first();

        if ($staffSetting) {
            return StaffSetting::query()->where('enterprise_id', $enterpriseId)->update($setting);
        } else {
            $setting['enterprise_id'] = $enterpriseId;
            $setting['type']          = StaffSetting::STAFF_TYPE;
            return StaffSetting::query()->create($setting);
        }
    }

    /**
     * 保存员工内购时处理员工内购分享绑定关系确定或解除
     *
     * @param $input
     *
     * @return array
     * @throws \Exception
     */
    public function staffRightUnbindStaffShare($input)
    {
        $enterpriseId   = $input['enterprise_id'] ?? 0;
        $staffRightList = $input['staff_buy_list'];

        $memberIds     = [];
        $departmentIds = [];

        foreach ($staffRightList as $value) {
            if ($value['obj_type'] == StaffRight::TYPE_STAFF) {
                $memberIds[] = $value['obj_id'];
            }
            if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                $departmentIds[] = $value['obj_id'];
            }
        }

        Log::info('保存的员工内购权限集合：' . json_encode($staffRightList));

        $staffRight = StaffRight::query()->where('enterprise_id', $enterpriseId)->get()->toArray();

        Log::info('员工内购权限集合：' . json_encode($staffRight));

        $insertDepartmentMemberIds = [];
        $deleteDepartmentMemberIds = [];

        if ($staffRight) {
            $staffRightMemberIds     = [];
            $staffRightDepartmentIds = [];
            foreach ($staffRight as $value) {
                if ($value['obj_type'] == StaffRight::TYPE_STAFF) {
                    $staffRightMemberIds[] = $value['obj_id'];
                }
                if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                    $staffRightDepartmentIds[] = $value['obj_id'];
                }
            }
            // 移除员工
            $deleteMemberIds = array_values(array_unique(array_diff($staffRightMemberIds, $memberIds)));
            // 新增员工
            $insertMemberIds = array_values(array_unique(array_diff($memberIds, $staffRightMemberIds)));
            // 移除部门
            $deleteDepartmentIds = array_values(array_unique(array_diff($staffRightDepartmentIds, $departmentIds)));
            // 新增部门
            $insertDepartmentIds = array_values(array_unique(array_diff($departmentIds, $staffRightDepartmentIds)));
            Log::info('移除部门ID集合：' . json_encode($deleteDepartmentIds));
            Log::info('新增部门ID集合：' . json_encode($insertDepartmentIds));
            if ($deleteDepartmentIds) {
                $deleteDepartmentMember    = (new Member())->ignorePermissionsList($input['qy_domain'], $input['app_identifier'], ['dpIdList' => $deleteDepartmentIds, 'departmentChildrenFlag' => 1, 'memSubscribeStatus' => 1]);
                $deleteDepartmentMemberIds = array_column($deleteDepartmentMember['list'], 'memUid');
            }
            if ($insertDepartmentIds) {
                $insertDepartmentMember    = (new Member())->ignorePermissionsList($input['qy_domain'], $input['app_identifier'], ['dpIdList' => $insertDepartmentIds, 'departmentChildrenFlag' => 1, 'memSubscribeStatus' => 1]);
                $insertDepartmentMemberIds = array_column($insertDepartmentMember['list'], 'memUid');
            }
            $insertAllMemberIds = array_values(array_filter(array_unique(array_merge($insertDepartmentMemberIds, $insertMemberIds))));
            $deleteAllMemberIds = array_values(array_filter(array_unique(array_merge($deleteDepartmentMemberIds, $deleteMemberIds))));
        } else {
            if ($departmentIds) {
                $departmentMember          = (new Member())->ignorePermissionsList($input['qy_domain'], $input['app_identifier'], ['dpIdList' => $departmentIds, 'departmentChildrenFlag' => 1, 'memSubscribeStatus' => 1]);
                $insertDepartmentMemberIds = array_column($departmentMember['list'], 'memUid');
            }
            $insertMemberIds    = $memberIds;
            $insertAllMemberIds = array_values(array_filter(array_unique(array_merge($insertDepartmentMemberIds, $insertMemberIds))));
            $deleteAllMemberIds = [];
        }

        Log::info('移除员工内购集合：' . json_encode($deleteAllMemberIds));

        Log::info('新增员工内购集合：' . json_encode($insertAllMemberIds));

        if ($deleteAllMemberIds) {
            // 移除员工内购时
            (new StaffSettingService())->unBindingStaffShareBymemUid($deleteAllMemberIds);
        }

        if ($insertAllMemberIds) {
            // 添加员工内购时
            (new StaffRightSDK())->unbindStaffShareByOwn(['mem_uids' => $insertAllMemberIds]);
        }

        return [];
    }

    /**
     * 更新员工理员工内购分享绑定关系确定或解除
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param $enterpriseId
     * @param $dpId
     * @param $params
     *
     * @throws \Exception
     */
    public function editMemberUnbindStaffShare($qyDomain, $appIdentifier, $enterpriseId, $dpId, $params)
    {
        // 员工ID
        $memUid = $params['memUid'];

        // 内购员工权限
        $staffRight = StaffRight::query()->where('enterprise_id', $enterpriseId)->get()->toArray();

        if ($staffRight) {
            // 内购员工权限员工ID集合
            $staffRightMemberIds = [];
            // 内购员工权限部门ID集合
            $staffRightDepartmentIds = [];
            foreach ($staffRight as $value) {
                if ($value['obj_type'] == StaffRight::TYPE_STAFF) {
                    $staffRightMemberIds[] = $value['obj_id'];
                }
                if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                    $staffRightDepartmentIds[] = $value['obj_id'];
                }
            }
            // 如果存在内购员工权限员工ID并且员工ID不在内购员工权限员工ID集合里并且存在内购员工权限部门ID，或者不存在内购员工权限员工ID并且存在内购员工权限部门ID
            if (($staffRightMemberIds && !in_array($memUid, $staffRightMemberIds) && $staffRightDepartmentIds) || (!$staffRightMemberIds && $staffRightDepartmentIds)) {
                $memberDetail       = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $memUid]);
                $memberDpIds        = array_column($memberDetail['dpName'], 'dpId');
                $allDepartmentDpIds = $staffRightDepartmentIds;
                foreach ($staffRightDepartmentIds as $departmentId) {
                    $departmentResult   = (new Department())->list($qyDomain, $appIdentifier, ['dpId' => $departmentId, 'departmentChildrenFlag' => 1], false);
                    $departmentDpIds    = array_column($departmentResult['list'], 'dpId');
                    $allDepartmentDpIds = array_merge($allDepartmentDpIds, $departmentDpIds);
                }
                $intersectDpIds = array_intersect($allDepartmentDpIds, $memberDpIds);
                Log::info('员工部门ID集合：' . json_encode($memberDpIds));
                Log::info('内购员工部门ID(递归)集合：' . json_encode($allDepartmentDpIds));
                Log::info('更改部门ID：' . $dpId);
                if (!$intersectDpIds && in_array($dpId, $allDepartmentDpIds)) {
                    // 编辑员工时部门在员工内购时
                    (new StaffRightSDK())->unbindStaffShareByOwn(['mem_uids' => [$memUid]]);
                }
                if ($intersectDpIds && !in_array($dpId, $allDepartmentDpIds)) {
                    // 编辑员工时部门不在员工内购时
                    (new StaffSettingService())->unBindingStaffShareBymemUid([$memUid]);
                }
            }
        }
    }

    /**
     * 更新员工员工内购部门
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param $enterpriseId
     * @param $dpId
     * @param $params
     *
     * @throws \Exception
     */
    public function editMemberStaffRight($qyDomain, $appIdentifier, $enterpriseId, $dpId, $params)
    {
        // 员工ID
        $memUid = $params['memUid'];

        $memberStaffRight = StaffRight::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('obj_type', StaffRight::TYPE_STAFF)
            ->where('obj_id', $memUid)
            ->first();

        if ($memberStaffRight) {
            $memberDetail = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $memUid]);
            $memberDpIds  = array_column($memberDetail['dpName'], 'dpId');

            if (!in_array($memberStaffRight->obj_parent_id, $memberDpIds)) {
                StaffRight::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('obj_type', StaffRight::TYPE_STAFF)
                    ->where('obj_id', $memUid)
                    ->update([
                        'obj_parent_id' => $dpId,
                    ]);
            }
        }
    }

    /**
     * 删除企业下允许参与内购的员工
     *
     * @param $enterpriseId
     *
     * @return mixed
     */
    public function deleteStaffRight($enterpriseId)
    {
        return StaffRight::query()->where('enterprise_id', $enterpriseId)->delete();
    }

    /**
     * 查询返回允许参与内购员工的列表
     *
     * @param $enterpriseId
     *
     * @return array
     */
    public function getStaffRightList($enterpriseId)
    {
        return StaffRight::query()
            ->select(['obj_id', 'obj_parent_id', 'obj_type', 'label', 'value'])
            ->where('enterprise_id', $enterpriseId)
            ->get()
            ->toArray();
    }

    /**
     * 查询返回内购员工设置
     *
     * @param $enterpriseId
     * @param $type
     * @return array
     */
    public function getStaffSetting($enterpriseId, $type)
    {
        if (is_array($enterpriseId)) {
            $staffSetting = StaffSetting::query()
                ->select(['enterprise_id', 'binding_type', 'duration_time', 'change_owner', 'protected_time', 'commission_type', 'settle_type', 'achievement_desc', 'fission_type'])
                ->whereIn('enterprise_id', $enterpriseId)
                ->where('type', '=', $type)
                ->get();
        } else {
            $staffSetting = StaffSetting::query()
                ->select(['enterprise_id', 'binding_type', 'duration_time', 'change_owner', 'protected_time', 'commission_type', 'settle_type', 'achievement_desc', 'fission_type'])
                ->where('enterprise_id', $enterpriseId)
                ->where('type', '=', $type)
                ->first();
        }

        return $staffSetting ? $staffSetting->toArray() : [];
    }

    /**
     * 查询返回允许参与内购员工的列表
     *
     * @param $enterpriseId
     * @param $qyDomain
     * @param $appIdentifier
     *
     * @return array
     */
    public function getStaffRight($enterpriseId, $qyDomain, $appIdentifier)
    {
        $staffList = StaffRight::query()
            ->select(['obj_id', 'obj_type'])
            ->where('enterprise_id', $enterpriseId)
            ->get();

        if (empty($staffList)) {
            return [];
        }

        // 分别获取员工和部门的id集合
        $memUids = [];
        $depIds  = [];
        foreach ($staffList as $value) {
            if ($value->obj_type == StaffRight::TYPE_STAFF) {
                $memUids[] = $value->obj_id;
            } else {
                $depIds[] = $value->obj_id;
            }
        }

        // 获取员工对应的数据列表
        $memList = $this->getStaffListFromUc($qyDomain, $appIdentifier, ['memUids' => $memUids]);
        $memList = array_column($this->formatUcMember($memList), null, 'mem_uid');

        // 获取部门对应数据列表
        $depList = $this->getDepartmentList($qyDomain, $appIdentifier);
        $depList = array_column($depList, null, 'dp_id');

        // 组装数据
        foreach ($staffList as &$item) {
            if ($item->obj_type == StaffRight::TYPE_STAFF) {
                $item->obj_name = isset($memList[$item->obj_id]) ? $memList[$item->obj_id]['mem_username'] : '';
            } else {
                $item->obj_name = isset($depList[$item->obj_id]) ? $depList[$item->obj_id]['dp_name'] : '';
            }
        }

        unset($item);
        return $staffList ? $staffList->toArray() : [];
    }

    /**
     * 判断用户是否具有内购权限
     *
     * @param        $enterpriseId
     * @param        $memUid
     * @param        $qyDomain
     * @param string $appIdentifier
     *
     * @return int
     * @throws \Exception
     */
    public function staffIsBuy($enterpriseId, $memUid, $qyDomain, $appIdentifier = 'contact')
    {
        $isStaffBuy = 0; // 默认不是内购员工

        // 查询数据库看是否存在此员工
        $dataStaff = (new StaffRight())->getStaffRightFirst($enterpriseId, $memUid);
        if (!empty($dataStaff)) {
            $isStaffBuy = 1;
        }

        return $isStaffBuy;
    }

    /**
     * 获取部门列表信息
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param $params
     *
     * @return array
     */
    public function getDepartmentList($qyDomain, $appIdentifier, $params = [])
    {
        // 获取部门对应的名称
        try {
            $params['pageSize'] = self::PAGESIZE;
            $depList            = (new Department())->list($qyDomain, $appIdentifier, $params);

        } catch (\Exception $e) {

            notice(720002);

        }

        $data = [];
        foreach ($depList['list'] as $item) {
            $data[] = [
                'dp_id'   => $item['dpId'],
                'dp_name' => $item['dpName']
            ];
        }

        return $data;
    }

    /**
     * 从uc获取员工列表数据
     *
     * @param $qyDomain
     * @param $appIdentifier
     * @param $params
     *
     * @return array|mixed
     */
    public function getStaffListFromUc($qyDomain, $appIdentifier, $params)
    {
        $result             = [];
        $params['pageSize'] = self::PAGESIZE;
        $params['memAll']   = 1;

        try {

            $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $params);

        } catch (\Exception $e) {

            notice(720002);
        }

        return $result;
    }

    /**
     * 格式化员工列表数据
     *
     * @param $result
     *
     * @return array
     */
    public function formatUcMember($result)
    {
        if (empty($result)) {
            return [];
        }

        $res = [];

        foreach ($result['list'] as $value) {
            $res[] = [
                'mem_uid'      => $value['memUid'],
                'mem_mobile'   => $value['memMobile'],
                'mem_username' => $value['memUsername'],
                'mem_face'     => $value['memFace'],
                'mem_qr_code'  => $value['memQrCode'] ?? '',
                'dp_name'      => isset($value['dpName']) ? $value['dpName'][0]['dpName'] : "",
                'dp_id'        => isset($value['dpName']) ? $value['dpName'][0]['dpId'] : 0,
                'dp_ids'       => array_column($value['dpName'] ?? [], 'dpId') ?? []
            ];
        }

        return $res;
    }

    /**
     * 查询内购员工对应部门
     *
     * @param $enterpriseId
     * @param $memUid
     * @param $qyDomain
     * @param $appIdentifier
     *
     * @return mixed|string
     * @throws \Exception
     */
    public function getStaffRightDepartmentRelation($enterpriseId, $memUid, $qyDomain, $appIdentifier)
    {
        $staffRightList = StaffRight::query()
            ->where('enterprise_id', $enterpriseId)
            ->get()
            ->toArray();

        $dpId   = '';
        $dpName = '';

        Log::info('员工内购权限集合：' . json_encode($staffRightList));

        if ($staffRightList) {
            $staffDetail = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $memUid]);
            $staffDpIds  = array_column($staffDetail['dpName'], 'dpId');
            Log::info('员工ID：' . $memUid);
            Log::info('员工部门ID集合：' . json_encode($staffDpIds));

            foreach ($staffRightList as $staffRight) {
                if ($staffRight['obj_type'] == StaffRight::TYPE_STAFF) {
                    if ($staffRight['obj_id'] == $memUid) {
                        if (in_array($staffRight['obj_parent_id'], $staffDpIds)) {
                            $dpId = $staffRight['obj_parent_id'];
                        } else {
                            $dpId = current($staffDpIds);
                        }
                        break;
                    }
                }
                if ($staffRight['obj_type'] == StaffRight::TYPE_DEP) {
                    if (in_array($staffRight['obj_id'], $staffDpIds)) {
                        $dpId = $staffRight['obj_id'];
                        break;
                    } else {
                        $departmentResult = (new Department())->list($qyDomain, $appIdentifier, ['dpId' => $staffRight['obj_id'], 'departmentChildrenFlag' => 1], false);
                        $departmentDpIds  = array_column($departmentResult['list'], 'dpId');
                        $intersectDpIds   = array_intersect($departmentDpIds, $staffDpIds);
                        if ($intersectDpIds) {
                            $dpId = current($intersectDpIds);
                            break;
                        }
                    }
                }
            }
        }

        if ($dpId) {
            $department = (new Department())->detail($qyDomain, $appIdentifier, $dpId);
            $dpName     = $department['dpName'] ?? '';
        }

        Log::info('员工部门ID：' . $dpId);

        return [$dpId, $dpName];
    }

    public function getStaffRightDep($enterpriseId)
    {

        //现有所有具有内购权限的员工和部门集合
        $staffRight = StaffRight::query()->where('enterprise_id', $enterpriseId)->get()->toArray();
        //现有所有具有内购权限的员工集合
        $staffRightMemberIds = [];
        //现有所有具有内购权限的部门集合
        $staffRightDepartmentIds = [];
        foreach ($staffRight as $value) {
            if ($value['obj_type'] == StaffRight::TYPE_STAFF) {
                $staffRightMemberIds[] = $value['obj_id'];
            }
            if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                $staffRightDepartmentIds[] = $value['obj_id'];
            }
        }

        $departmentMember = (new Member())->ignorePermissionsList(
            $input['qy_domain'],
            $input['app_identifier'],
            ['dpIdList' => $staffRightDepartmentIds, 'departmentChildrenFlag' => 1, 'memSubscribeStatus' => 1, 'pageSize' => 5000]
        );
    }

    /**
     * 内购设置为所有内购员工自建联系人二维码
     *
     * @param $input
     *
     * @return bool
     * @author: DLong
     */
    public function syncContactCodeByMem($input)
    {
        Log::info('后台内购设置自建联系人二维码--队列开始', ['$input' => $input]);

        $enterpriseId   = $input['enterprise_id'];
        $staffRightList = $input['staff_buy_list'];

        //本次新增员工集合
        $memberIds = [];
        //本次新增部门集合
        $departmentIds = [];

        foreach ($staffRightList as $value) {
            if ($value['obj_type'] == StaffRight::TYPE_STAFF) {
                $memberIds[] = $value['obj_id'];
            }
            if ($value['obj_type'] == StaffRight::TYPE_DEP) {
                $departmentIds[] = $value['obj_id'];
            }
        }
        //具有内购权限下的所有部门
        $departmentMember = [];
        if (!empty($departmentIds)) {
            $departmentMember = (new Member())->ignorePermissionsList(
                $input['qy_domain'],
                $input['app_identifier'],
                ['dpIdList' => $departmentIds, 'departmentChildrenFlag' => 1, 'memSubscribeStatus' => 1, 'pageSize' => 500]
            );
        }
        //具有内购权限的所有部门下的所有员工
        $insertDepartmentMemberIds = array_column($departmentMember['list'] ?? [], 'memUid');
        //具有内购权限的所有部门下的所有员工+本次新增员工=所有具有内购权限的员工
        $insertAllMemberIds = array_values(array_filter(array_unique(array_merge($insertDepartmentMemberIds, $memberIds))));

        //所有已经有联系人二维码的内购员工
        $allHasContactMem = ContactQrCode::query()->where([
            ['type', ContactQrCodeParams::MEMBER],
            ['source_type', ContactQrCodeParams::SOURCE_TYPE_AUTO],
            ['enterprise_id', $enterpriseId],
        ])->pluck('id')->toArray();
        $allHasContactMem = ContactQrCodeParams::query()->whereIn('contact_qr_code_id', $allHasContactMem)
            ->pluck('object')->toArray();

        $insertAllMemberIds = array_values(array_unique(array_diff($insertAllMemberIds, $allHasContactMem)));

        // 查询所有员工
        $normalMembers = MemberService::memberList(
            $input['qy_domain'], $input['app_identifier'], ['memSubscribeStatus' => 1]
        );
        //匹配userId
        $ucAllUsers = array_column($normalMembers, 'memUserid', 'memUid');
        //匹配员工姓名
        $ucAllMember = array_column($normalMembers, 'memUsername', 'memUid');
        //每个员工自建二维码
        foreach ($insertAllMemberIds as $value) {
            if (empty($ucAllUsers[$value]) || empty($ucAllMember[$value])) {
                continue;
            }
            try {
                $this->syncContactQrCode(
                    $enterpriseId,
                    $input['qy_domain'],
                    $ucAllUsers[$value],
                    $value,
                    $ucAllMember[$value],
                    ContactQrCode::SOURCE_TYPE_AUTO
                );
            } catch (\Exception $e) {
                Log::info(' 自建二维码', ['Error' => $e->getMessage()]);
                continue;
            }
        }
        return true;
    }


    /**
     * 自动同步所有企业具有内购权限的部门下的所有员工的联系人二维码
     *
     * @param $input
     *
     * @author: DLong
     */
    public function autoSyncContactCode($enterpriseId, $domain, $identifier)
    {
        try {
            $this->setLeaveMem($enterpriseId, $domain, $identifier);
        } catch (\Exception $e) {
            Log::info('获取所有需要创建二维码的员工数据异常', ['Error' => $e->getMessage()]);
        }

        return true;
    }

    /**
     * 为单个员工创建联系人二维码
     *
     * @param        $enterpriseId
     * @param        $domain
     * @param        $ucMemId
     * @param        $ucMemName
     * @param int $adminId
     * @param string $adminName
     *
     * @throws \Error
     * @author: DLong
     */
    public function syncContactQrCode($enterpriseId, $domain, $ucUserId, $ucMemId, $ucMemName, $sourceType)
    {
        try {
            DB::beginTransaction();
            //根据渠道来源定义名称
            $channelName = '';
            //联系人二维码名称前缀
            $qyCodePrefix = '';
            switch ($sourceType) {
                case ContactQrCode::SOURCE_TYPE_AUTO:
                    $channelName  = '员工亲友内购';
                    $qyCodePrefix = '内购';
                    break;
                case ContactQrCode::SOURCE_TYPE_AUTO_GUIDE:
                    $channelName  = '支付即导购';
                    $qyCodePrefix = '支付即导购';
                    break;
            }
            //对应渠道是否存在
            $channelCollect = Channel::query()->where('enterprise_id', $enterpriseId)
                ->where('source_type', $sourceType)
                ->get(['id', 'channel']);
            $channelObj     = $channelCollect->where('channel', $channelName)->first();
            if (!$channelObj) {
                //同步渠道
                $channelData = [
                    'enterprise_id' => $enterpriseId,
                    'channel'       => $channelName,
                    'source_type'   => $sourceType,
                    'admin_id'      => 0,
                    'admin_name'    => '系统'
                ];
                $channelObj  = Channel::query()->create($channelData);
                if (empty($channelObj)) {
                    Log::info('同步联系人二维码渠道异常！');
                }
            }
            //联系人二维码名称
            $qrCodeName = $this->getLastId($sourceType, $ucMemName);

            // 本地存储二维码
            $contactQr = ContactQrCode::query()->create([
                'enterprise_id' => $enterpriseId,
                'name'          => $qyCodePrefix . $qrCodeName,
                'is_verify'     => ContactQrCode::VERIFY_NO,
                'channel_id'    => $channelObj['id'],
                'channel'       => $channelName,
                'type'          => ContactQrCode::CODE_TYPE_ONE,
                'source_type'   => $sourceType,
            ]);

            ContactQrCodeParams::query()->create([
                'enterprise_id'      => $enterpriseId,
                'contact_qr_code_id' => $contactQr['id'],
                'type'               => ContactQrCodeParams::MEMBER,
                'object'             => $ucMemId,
                'object_extra'       => $ucMemId,
            ]);

            $ucResult = (new Member())->getContactQrCode($domain, [
                "ecwScene"      => 2,//场景，1-在小程序中联系，2-通过二维码联系
                'ecwType'       => ContactQrCode::CODE_TYPE_ONE,
                'ecwRemark'     => $channelObj->channel,
                "ecwSkipVerify" => ContactQrCode::VERIFY_NO,
                "ecwState"      => $contactQr['id'],
                'memUserid'     => $ucUserId,
            ]);
            if (!is_array($ucResult) || empty($ucResult['qrCode'])) {
                notice(720002);
            }

            // 企业微信二维码链接
            $officialUrl = $ucResult['qrCode'];

            // 企业微信二维码存储到腾讯云
            $transformUrl = CommonService::uploadImage($officialUrl);
            // 更新本地数据
            $contactQr['uc_id']         = $ucResult['id'];
            $contactQr['config_id']     = $ucResult['configId'];
            $contactQr['official_url']  = $officialUrl;
            $contactQr['transform_url'] = $transformUrl;

            $contactQr->save();
            DB::commit();

            return $contactQr->toArray();
        } catch (\Exception $e) {
            DB::rollback();
            Log::info('事务创建二维码！', ['ErrorInfo---' => $e->getMessage()]);
        }

        return [];
    }

    /**
     * 更新联系我二维码及导购信息
     *
     * @param $enterpriseId
     * @param $domain
     * @param $identifier
     *
     * @return array
     * @author: DLong
     */
    public function getStaffMem($enterpriseId, $domain, $identifier)
    {
        try {
            // 查询所有员工
            $normalMembers = MemberService::memberList($domain, $identifier, ['memSubscribeStatus' => 1]);
            $normalMemUids = array_column($normalMembers, 'memUid');
            //处理离职员工
            $this->setLeaveMem($normalMemUids, $enterpriseId);

        } catch (\Exception $e) {
            Log::info('获取所有需要创建二维码的内购员工数组异常', ['Error' => $e->getMessage()]);
        }

        return true;
    }

    /**
     * 获取联系人二维码最后一个名称编号
     *
     * @param $sourceType_联系人二维码来源类型（1 =内购，3=支付既导购）
     * @param $userName_员工姓名
     *
     * @return string
     * @author: DLong
     */
    public function getLastId($sourceType, $userName)
    {
        $lastName = ContactQrCode::query()
            ->where('source_type', $sourceType)
            ->orderBy('id', 'desc')
            ->first(['name']);
        $num      = 1;
        if (!empty($lastName->name) && preg_match('/\d+/', $lastName->name, $arr)) {
            $num      = $arr[0] + 1;
            $userName = $userName . $num;
        } else {
            return $userName . $num;
        }
        return $userName;
    }


    /**
     * 处理离职员工二维码和移除导购
     * @param $enterpriseId
     * @param $domain
     * @param $identifier
     * @return bool
     */
    public function setLeaveMem($enterpriseId, $domain, $identifier)
    {
        $qrCodeMem = ContactQrCodeParams::query()
            ->where('enterprise_id', $enterpriseId)
            ->get(['id', 'object_extra'])
            ->toArray();

        if (!$qrCodeMem) {
            return true;
        }

        // 查询所有员工
        $normalMembers = MemberService::memberList($domain, $identifier, ['memSubscribeStatus' => 1]);
        $normalMemUids = array_column($normalMembers, 'memUid');

        //同步减少对应二维码包含人数
        foreach ($qrCodeMem as $val) {
            $array          = explode(',', $val['object_extra']);
            $arrByintersect = array_intersect($array, $normalMemUids);
            if ($arrByintersect != $array) {
                $str = implode(',', $arrByintersect);
                ContactQrCodeParams::query()
                    ->where('id', $val['id'])
                    ->update(['object_extra' => $str]);
            }
        }

        return true;
    }

    /**
     * 自动同步所有企业具有内购权限的部门下的所有员工的联系人二维码
     *
     * @param $input
     *
     * @author: DLong
     */
    public function autoStaffData($enterpriseId, $domain, $identifier)
    {
        try {
            list($insertAllMemberIds, $ucAllUsers, $ucAllMember) = $this->getStaffDataMem($enterpriseId, $domain, $identifier);
        } catch (\Exception $e) {
            Log::info('获取所有需要创建二维码的内购员工数据异常', ['Error' => $e->getMessage()]);
        }
        //队列为每个员工自建二维码
        $num = 0;
        foreach ($insertAllMemberIds ?? [] as $value) {
            if (empty($ucAllUsers[$value]) || empty($ucAllMember[$value])) {
                continue;
            }
            $contactQrCodeData = $this->syncContactQrCode(
                $enterpriseId,
                $domain,
                $ucAllUsers[$value] ?? '',
                $value,
                $ucAllMember[$value] ?? '',
                ContactQrCode::SOURCE_TYPE_AUTO
            );

            // 创建二维码连续失败超过十次break;
            if ($num == 10) {
                break;
            }
            if (!empty($contactQrCodeData) && !empty($contactQrCodeData['official_url'])) {
                $num = 0;
            } else {
                $num++;
            }
        }
        return true;
    }

    /**
     * 获取所有需要创建联系人二维码的内购员工和姓名数组
     *
     * @param $enterpriseId
     * @param $domain
     * @param $identifier
     *
     * @return array
     * @author: DLong
     */
    public function getStaffDataMem($enterpriseId, $domain, $identifier)
    {
        try {
            // 同步之前所有具有内购权限的员工
            $staffRightMemberIds = StaffRight::query()->where('obj_type', StaffRight::TYPE_STAFF)
                ->where('enterprise_id', $enterpriseId)->pluck('obj_id')->toArray();
            //同步之前所有具有内购权限的部门
            $staffRightDepartment = StaffRight::query()->where('obj_type', StaffRight::TYPE_DEP)
                ->where('enterprise_id', $enterpriseId)->get()->toArray();
            if (!$staffRightMemberIds && !$staffRightDepartment) {
                return [[], [], []];
            }

            //同步更新后的所有部门
            $newDpList = MemberService::allDepList($domain, $identifier);
            $newDpIds  = array_column($newDpList, 'dpId');
            $oldDpIds  = array_column($staffRightDepartment, 'obj_id');
            //更新后消失的部门
            $decreaseDpIds = array_diff($oldDpIds, $newDpIds);
            if ($decreaseDpIds) {
                //减少的部门在内购权限里同步删除
                StaffRight::query()->where('enterprise_id', $enterpriseId)
                    ->where('obj_type', StaffRight::TYPE_DEP)
                    ->whereIn('obj_id', $decreaseDpIds)->delete();
            }

            //同步通讯录之前所有有内购权限部门下的所有员工
            $staffByDp = [];
            //同步之后所有内购部门下所有员工
            $allNewMemUidByDp = [];
            //消失的部门附带的所有员工
            $allMemUidByDecreaseDp = [];
            foreach ($staffRightDepartment as $value) {
                $objectIds = explode(',', $value['object_extra_id']);
                $staffByDp = array_merge($staffByDp, $objectIds);
                //如果此部门消失，合并所有附带员工
                if ($value['obj_type'] == StaffRight::TYPE_DEP && in_array($value['obj_id'], $decreaseDpIds)) {
                    $allMemUidByDecreaseDp = array_merge(
                        $allMemUidByDecreaseDp, explode(',', $value['object_extra_id'])
                    );
                    continue;
                }
                //同步通讯录以后对应部门下的新员工
                $membersByDp    = MemberService::memberList($domain, $identifier, [
                    'dpId'                   => $value['obj_id'],
                    'departmentChildrenFlag' => 1,
                    'memSubscribeStatus'     => 1,
                ]);
                $newMemIdsByDep = array_column($membersByDp, 'memUid');

                if ($newMemIdsByDep != $objectIds) {
                    //更新对应部门现有员工为本次同步之后的员工id集合
                    StaffRight::query()->where('id', $value['id'])
                        ->update(['object_extra_id' => implode(',', $newMemIdsByDep)]);
                }

                $allNewMemUidByDp = array_merge($allNewMemUidByDp, $newMemIdsByDep);
            }
            //同步之前所有内购员工
            $currentAllStaffMem = array_unique(array_merge($staffRightMemberIds, $staffByDp));

            // 查询所有员工
            $normalMembers = MemberService::memberList($domain, $identifier, ['memSubscribeStatus' => 1]);
            $normalMemUids = array_column($normalMembers, 'memUid');
            //处理离职员工
//            $this->setLeaveMem($normalMemUids, $enterpriseId);
            //本次同步所有减少的员工
            $leaveAllMem = array_diff($currentAllStaffMem, $normalMemUids);
            if ($leaveAllMem) {
                try {
                    // 解绑链路关系
                    (new StaffSettingService())->unBindingStaffShareBymemUid($leaveAllMem);
                } catch (\Exception $e) {
                    Log::info('解绑链路关系异常', ['Error' => $e->getMessage()]);
                }
                //减少的员工在内购权限里同步删除
                StaffRight::query()->where('enterprise_id', $enterpriseId)
                    ->where('obj_type', StaffRight::TYPE_STAFF)
                    ->whereIn('obj_id', $leaveAllMem)->delete();
            }

            $increasedAllMem = array_diff(array_unique($allNewMemUidByDp), $currentAllStaffMem);
            if ($increasedAllMem) {
                //确定链路关系
                try {
                    (new StaffRightSDK())->unbindStaffShareByOwn(['mem_uids' => $increasedAllMem]);
                } catch (\Exception $e) {
                    Log::info('确定链路关系异常', ['Error' => $e->getMessage()]);
                }
            }

            //所有已经有联系人二维码的内购员工
            $allHasContactMem = ContactQrCode::query()->where([
                ['type', ContactQrCodeParams::MEMBER],
                ['source_type', ContactQrCodeParams::SOURCE_TYPE_AUTO],
                ['enterprise_id', $enterpriseId],
            ])->pluck('id')->toArray();
            $allHasContactMem = ContactQrCodeParams::query()
                ->whereIn('contact_qr_code_id', $allHasContactMem)
                ->pluck('object')->toArray();

            //现有所有具有内购权限的员工id数组
            $staffMemberIds = StaffRight::query()->where('obj_type', StaffRight::TYPE_STAFF)
                ->where('enterprise_id', $enterpriseId)
                ->pluck('obj_id')->toArray();
            //同步之后所有的内购员工
            $newAllMemIds = array_merge($allNewMemUidByDp, $staffMemberIds);

            //消失的部门附带的所有员工在同步之后是否有内购权限
            $noStaffMem = array_diff(array_unique($allMemUidByDecreaseDp), $newAllMemIds);
            if ($noStaffMem) {
                try {
                    // 解绑链路关系
                    (new StaffSettingService())->unBindingStaffShareBymemUid($noStaffMem);
                } catch (\Exception $e) {
                    Log::info('解绑链路关系异常', ['Error' => $e->getMessage()]);
                }
            }
            $insertAllMemberIds = array_values(array_unique(array_diff($newAllMemIds, $allHasContactMem)));
            //匹配userId
            $ucAllUsers = array_column($normalMembers, 'memUserid', 'memUid');
            //匹配员工姓名
            $ucAllMember = array_column($normalMembers, 'memUsername', 'memUid');

            return [$insertAllMemberIds, $ucAllUsers, $ucAllMember];
        } catch (\Exception $e) {
            Log::info('获取所有需要创建二维码的内购员工数组异常', ['Error' => $e->getMessage()]);
            return [[], [], []];
        }
    }


    /**
     * getStoreWorkerByIds
     * @desc 根据workerid从gw_shop查询导购信息
     * @param $enterpriseId
     * @param $workerIdList
     *
     * @return mixed
     * @date 2020-07-20
     * @throws \Exception
     */
    public static function getStoreWorkerByIds($enterpriseId, $workerIdList) {
        return StoreWorkerService::list([
            'enterprise_id' => $enterpriseId,
            'ids' => $workerIdList], false);
    }

    /**
     * 判断用户是否有内购权限
     *
     * @param $userMemUid
     * @param $enterpriseHash
     * @return bool
     * @throws \Exception
     */
    public function staffIsRight($userMemUid, $enterpriseHash)
    {
        // 获取企业信息
        $enterpriseInfo = (new Enterprise())->detail(['hash' => $enterpriseHash]);
        if (empty($enterpriseInfo['data'])) {
            throw new Error(6006036, '企业信息不存在');
        }

        // 判断是否有内购权限
        if (!empty($userMemUid)) {
            $isRight = $this->staffIsBuy($enterpriseInfo['data']['id'], $userMemUid, $enterpriseInfo['data']['domain'], 'contact');
            if ($isRight['data']['is_right']) {
                return true;
            }
        }

        return false;
    }
}
