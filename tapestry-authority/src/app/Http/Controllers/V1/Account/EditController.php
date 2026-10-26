<?php

namespace App\Http\Controllers\V1\Account;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Models\EnterpriseAdmin;
use App\Modules\BuiltInRole\RoleStrategy;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class EditController extends BaseController
{
    /**
     * 编辑账号
     *
     * @return JsonResponse
     * @throws ValidationException
     * @api {domain}/micro-authority/v1/account/edit
     * @author lijian lijian@vchangyi.com 2020/04/14
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'operator_id',
            'enterprise_id',
            'account',
            'username',
            'password',
            'mobile',
            'email',
            'range_type',
            'range',
            'admin_available_system',
            'role_ids',
        ]);

        $input['admin_available_system'] = $input['admin_available_system'] ?? [];

//        $this->checkPermissionsKey($input['enterprise_id'], $input['id']);

        // 验证
//        AccountRequest::edit($input);

        $adminInfo = EnterpriseAdmin::query()->find($input['id']);
        if (!$adminInfo) {
            notice(701017);
        }
        //验证手机号或Email或账户名称是否有重复
        if (
            $adminInfo->account != $input['account'] ||
            $adminInfo->mobile != $input['mobile']
        ) {
            $count = EnterpriseAdmin::query()
                ->where(function ($q) use ($input) {
                    $q->where('account', $input['account'])
                        ->orWhere('mobile', $input['mobile']);
                    //->orWhere('email', $input['email']);
                })
                ->get();
            if ($count) {
                notice(701018);
            }
        }

        if (isset($input['password']) && !empty($input['password'])) {
            $password = $input['password'];
            $rule     = '/^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{8,}$/';

            //如果有修改密码
            if ($adminInfo->password !== $password) {
                //验证密码规则
                if (strlen($password) < 8 || !preg_match($rule, $password)) {
                    notice(701019);
                }
                //密码加密
                $input['password'] = md5(md5($password) . $adminInfo->salt);
            }
        }

        // 更新数据
        try {
            DB::beginTransaction();
            // 编辑管理员表信息
            (new EnterpriseAdminService())->editEnterpriseAdmin($input);

            // 编辑范围表信息
            (new EnterpriseDataRangeService())->editDataRange($input);

            // 编辑角色表信息
            (new EnterpriseAdminRoleService())->editEnterpriseAdminRole($input);


            DB::commit();
        } catch (Exception $exception) {
            DB::rollBack();

            Log::info('编辑账号失败' . $exception->getMessage(), $input);
            throw new Error($exception->getCode(), $exception->getMessage());
        }

        return success([]);
    }
}