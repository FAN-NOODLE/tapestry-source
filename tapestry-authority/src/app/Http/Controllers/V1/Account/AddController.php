<?php

namespace App\Http\Controllers\V1\Account;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Models\EnterpriseAdmin;
use App\Models\Role;
use App\Modules\BuiltInRole\RoleStrategy;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;
use App\Services\EnterpriseService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class AddController extends BaseController
{
    /**
     * 创建企业管理员账号
     *
     * @return JsonResponse
     * @throws ValidationException
     * @api {domain}/micro-authority/v1/account/add
     * @author lijian lijian@vchangyi.com 2020/04/13
     */
    public function index()
    {
        $input = $this->request->only([
            'operator_id',
            'enterprise_id',
            'username',
            'account',
            'mobile',
            'password',
            'email',
            'range_type',
            'range',
            'admin_available_system',
            'role_ids',
        ]);

        // 验证
        AccountRequest::add($input);
        // 验证手机号
        $mobileArr = explode('-', $input['mobile']);
        $mobile    = $mobileArr[1] ?? $input['mobile'];
        $preg      = "/^1[3-9]\d{9}$/";
        if (!preg_match($preg, $mobile)) {
            notice(701001); // 请输入正确的手机号码
        }

        // 验证创建账号规则
        $detail = (new EnterpriseService())->enterpriseDetail($input['enterprise_id']);
        if (empty($detail)) {
            notice(701005); // 企业不存在
        }

        // 验证角色互斥
//        if (!empty($input['role_ids'])) {
//            Role::roleMutexLimit($input['role_ids']);
//            $ids                = Role::roleByMark(['guide_enterprise_admin'], $input['enterprise_id']);
//            $hasEnterpriseAdmin = array_intersect($ids, $input['role_ids']);
//            // 处理智慧导购企业管理角色
//            if ($hasEnterpriseAdmin) {
//                $input['admin_available_system'] = [2];
//            }
//        }

//        $accountCount = (new EnterpriseAdminService())->getEnterpriseChildAdminCount($input['enterprise_id']);
        $accountCount = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $input['enterprise_id'])
            ->where('admin_available_system', '!=', json_encode([2])) // 过滤智慧导购
            ->where('admin_available_system', '!=', '[]') // 过滤空的可用系统
            ->where('is_super', EnterpriseAdmin::NO_SUPER)
            ->where('status', EnterpriseAdmin::STATUS_OPEN)
            ->count();
        if ($detail['account_limit'] <= $accountCount) {
            notice(701004, [$detail['account_limit']]); // 已达到允许创建子账号上限
        }

        if (count($mobileArr) > 1) {
            $input['global_roaming'] = $mobileArr[0];
        }
        //$input['account'] = account;
        $input['mobile'] = $mobile;

        $this->addAccountDataByOLY($input);

        return success([]);
    }

    public function addAccountData($input)
    {
        try {
            DB::beginTransaction();

            // 在账号表创建角色
            $enterpriseAdminObj = (new EnterpriseAdminService())->addEnterpriseAdmin($input);
            if (!$enterpriseAdminObj) {
                Log::info('插入管理员表出现错误');
                throw new Error(1000000, '创建账号失败');
            }

            $input['enterprise_admin_id'] = $enterpriseAdminObj->id;

            // 插入账号范围数据
            if ($input['range_type'] == EnterpriseAdmin::RANGE_TYPE_ONLY) {
                (new EnterpriseDataRangeService())->addEnterpriseDataRange($input);
            }

            // 插入账号角色关联数据
            if (!empty($input['role_ids'])) {
                (new EnterpriseAdminRoleService())->addEnterpriseAdminRole($input);
            }

            DB::commit();
        } catch (Exception $exception) {
            DB::rollBack();
            $data = [
                'code'    => $exception->getCode(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'message' => $exception->getMessage(),
            ];
            Log::info('创建账号失败', $data);
            throw new Error($exception->getCode(), $exception->getMessage(), $data);
        }
    }

    public function addAccountDataByOLY($input)
    {
        try {
            DB::beginTransaction();

            // 在账号表创建角色
            $enterpriseAdminObj = (new EnterpriseAdminService())->addEnterpriseAdminByOLY($input);
            if (!$enterpriseAdminObj) {
                Log::info('插入管理员表出现错误');
                throw new Error(1000000, '创建账号失败');
            }

            $input['enterprise_admin_id'] = $enterpriseAdminObj->id;

            // 插入账号范围数据
            if ($input['range_type'] == EnterpriseAdmin::RANGE_TYPE_ONLY) {
                (new EnterpriseDataRangeService())->addEnterpriseDataRange($input);
            }

            // 插入账号角色关联数据
            if (!empty($input['role_ids'])) {
                (new EnterpriseAdminRoleService())->addEnterpriseAdminRole($input);
            }

            DB::commit();
        } catch (Exception $exception) {
            DB::rollBack();
            $data = [
                'code'    => $exception->getCode(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'message' => $exception->getMessage(),
            ];
            Log::info('创建账号失败', $data);
            throw new Error($exception->getCode(), $exception->getMessage(), $data);
        }
    }
}