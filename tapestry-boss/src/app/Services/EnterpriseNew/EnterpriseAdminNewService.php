<?php


namespace App\Services\EnterpriseNew;


use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Services\BaseService;
use Exception;
use Illuminate\Support\Facades\Log;

class EnterpriseAdminNewService extends BaseService
{
    public function addSuperAdmin($data)
    {
        if (EnterpriseAdmin::query()->where('email', '=', $data['email'])->exists()) {
            return notice(710003);
        }

        try {
            return EnterpriseAdmin::query()->create($data);
        } catch (Exception $e) {
            Log::info('插入企业超管数据失败:' . $e->getMessage(), $data);
            return error(710004);
        }
    }

    /**
     * 获取企业超管信息
     *
     * @param $enterpriseId
     * @return mixed
     */
    public function getSuperAdmin($enterpriseId)
    {
        try {
            return EnterpriseAdmin::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('is_super', EnterpriseAdmin::IS_SUPER)
                ->first();
        } catch (Exception $e) {
            Log::error('获取企业超管信息失败：' . $e->getMessage());
            return error(600000, ['获取企业超管信息失败']);
        }
    }

    /**
     * 更新超管信息
     *
     * @param $where
     * @param $data
     * @return int|void
     */
    public function updateSuperAdmin($where, $data)
    {
        try {
            return EnterpriseAdmin::query()->where($where)->update($data);
        } catch (Exception $e) {
            Log::error('更新超管信息：' . $e->getMessage());
            return error(600000, ['更新超管信息']);
        }
    }
}
