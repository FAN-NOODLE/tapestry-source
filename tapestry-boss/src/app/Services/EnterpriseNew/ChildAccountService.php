<?php


namespace App\Services\EnterpriseNew;

use App\Models\EnterpriseAdmin;
use App\Services\BaseService;

class ChildAccountService extends BaseService
{
    /**
     * 获取管理员列表
     * @param $enterpriseId
     * @param int $isSuper
     * @return array
     */
    public function getEnterpriseChildAdminList($enterpriseId, $isSuper = EnterpriseAdmin::NO_SUPER, $status = EnterpriseAdmin::STATUS_OPEN)
    {
        return EnterpriseAdmin::query()
            ->where('enterprise_id', $enterpriseId)
            ->when(in_array($isSuper, [EnterpriseAdmin::NO_SUPER, EnterpriseAdmin::IS_SUPER]), function ($query) use ($isSuper) {
                return $query->where('is_super', $isSuper);
            })
            ->when(in_array($status, [EnterpriseAdmin::STATUS_OPEN, EnterpriseAdmin::STATUS_CLOSE]), function ($query) use ($status) {
                return $query->where('status', $status);
            })
            ->get()
            ->toArray();
    }

    /**
     * 修改账号状态 1-启用 2-禁用
     *
     * @param $input
     *
     * @return bool|int
     */
    public function updateEnterpriseAdminStatus($input)
    {
        return EnterpriseAdmin::query()
            ->when(isset($input['id']) && $input['id'] > 0, function ($query) use ($input) {
                return $query->where('id', $input['id']);
            })
            ->when(
                isset($input['ids']) && is_array($input['ids']) && !empty($input['ids']),
                function ($query) use ($input) {
                    return $query->whereIn('id', $input['ids'])->where('is_super', EnterpriseAdmin::NO_SUPER);
                }
            )
            ->when(isset($input['all']) && $input['all'] == 'all', function ($query) {
                return $query->where('is_super', EnterpriseAdmin::NO_SUPER);
            })
            ->where('enterprise_id', $input['enterprise_id'])
            ->update(['status' => $input['status']]);
    }

}