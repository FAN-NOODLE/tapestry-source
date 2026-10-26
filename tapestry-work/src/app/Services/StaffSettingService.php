<?php


namespace App\Services;


use App\Models\StaffShare;
use App\Models\User;
use Illuminate\Support\Facades\DB;

class StaffSettingService extends BaseService
{
    public function unBindingStaffShareBymemUid(array $memUids)
    {
        StaffShare::query()
            ->whereIn('mem_uid', $memUids)
            ->where('binding_type', StaffShare::BINDING_TYPE_YES)
            ->update(['binding_type' => StaffShare::BINDING_TYPE_NO]);

        return true;
    }

    public function unBindingStaffShareByUserId(array $userIds)
    {
        StaffShare::query()
            ->whereIn('user_id' ,$userIds)
            ->where('binding_type', StaffShare::BINDING_TYPE_YES)
            ->update(['binding_type' => StaffShare::BINDING_TYPE_NO]);

        return true;
    }

    public function unBindingStaffShareById(array $Ids)
    {
        StaffShare::query()
            ->whereIn('id' ,$Ids)
            ->where('binding_type', StaffShare::BINDING_TYPE_YES)
            ->update(['binding_type' => StaffShare::BINDING_TYPE_NO]);

        return true;
    }

    public function changeOwner($parentUserId, $memUid, $shareUserInfo, $shareParentUserInfo, $isStaff)
    {
        // 上级未绑定
        if (!empty($shareParentUserInfo) && $shareParentUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO) {
            return false;
        }

        // 开始绑定时间
        $startTime = $shareUserInfo['created_at'] ?? 0;

        // 获取给则数据
        $settingInfo =  (new StaffManagerService())->getStaffSetting($shareUserInfo['enterprise_id'],1);

        $setting = $settingInfo ?? [];

        if (empty($setting)) {
            return false;
        }

        // 永久有效
        if ($setting['binding_type'] == StaffShare::BINDING_TYPE_ONE) {
            return false;
        }

        // 不允许抢客
        if ($setting['change_owner'] == StaffShare::CHANGE_OWNER_ONE) {
            return false;
        }

        // 保护期
        if (time() <= ($startTime + $setting['protected_time'] * 60 * 60 * 24)) {
            return false;
        }

        // 开始抢客
        if ($shareUserInfo['mem_uid'] != $memUid) {

            try {
                DB::beginTransaction();
                // 解绑之前关系
                $this->unBindingStaffShareByUserId([$shareUserInfo['user_id']]);

                // 绑定新的关系
                (new StaffShareService())
                    ->createStaffShareData($memUid, $parentUserId, $shareUserInfo['user_id'], $shareUserInfo['enterprise_id'], $shareParentUserInfo, $shareUserInfo, $isStaff);

                DB::commit();

                return true;

            } catch (\Exception $e) {
                DB::rollBack();

                return false;
            }
        }

        return false;
    }

    public function isRight($userId)
    {
        $userInfo = (new StaffShare())->getStaffShareInfoByCondition(['user_id' => $userId], false);

        return empty($userInfo) ? false : true;
    }

    public function handleBindUser($userId, $isStaff, $enterpriseId)
    {
        $staffShareService = new StaffShareService();
        $userMemUid = $staffShareService->getUserInfoById($enterpriseId, $userId)['mem_uid'];
        $bindUser = (new StaffShare())->getStaffShareInfoByCondition(['user_id' => $userId]);
        // 是内购员工解绑关系(管理后台异常保险起见)
        if ($isStaff && !empty($bindUser) && $bindUser['mem_uid'] != $userMemUid) {
            $this->unBindingStaffShareByUserId([$userId]);
            $staffShareService->clearCounter($userId);
        }
        // 原来是内购员工解绑关系(管理后台异常保险起见)
        if (!$isStaff && !empty($userMemUid)) {
            $bindMemUser = (new StaffShare())->getStaffShareInfoByCondition(['mem_uid' => $userMemUid]);
            if (!empty($bindMemUser)) {
                $this->unBindingStaffShareBymemUid([$userMemUid]);
                $staffShareService->clearCounter($userId);
            }
        }
    }
}
