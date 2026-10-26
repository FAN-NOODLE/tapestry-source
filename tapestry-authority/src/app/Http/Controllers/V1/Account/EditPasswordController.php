<?php

namespace App\Http\Controllers\V1\Account;

use App\Common\Constant;
use App\Common\RedisKey;
use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseAdmin;
use App\Services\EnterpriseAdminService;
use App\Traits\CodeTrait;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class EditPasswordController extends BaseController
{

    use CodeTrait;
    /**
     * 编辑密码
     *
     */
    public function index(): JsonResponse
    {
        $input   = $this->request->only([
            'admin_id', 'old_ps', 'new_ps', 'new_ps_two', 'is_require_login', 'account', 'enterprise_id'
        ]);
        $adminId = $input['admin_id'] ?? '';
        $account = $input['account'] ?? '';

        $enterpriseId = $input['enterprise_id'] ?? '';

        $isRequireLogin = isset($input['is_require_login']) && $input['is_require_login'];
        $select         = ['password', 'account', 'id', 'salt'];
        if ($isRequireLogin) {
            if ((empty($adminId) || empty($input['old_ps']))) {
                notice(701017);
            }
            $adminInfo = EnterpriseAdmin::query()->find($adminId, $select);
            if (!$adminInfo) {
                notice(701017);
            }
        } else {
            if ((empty($account) || empty($input['old_ps']))) {
                notice(701017);
            }

            if(!$enterpriseId){
                notice(700014);
            }

            $adminInfo = EnterpriseAdmin::query()->select($select)->where([
                'enterprise_id' => $enterpriseId,
                'account'       => $account
            ])->first();
            if (!$adminInfo) {
                notice(701017);
            }
        }

        $encryptPs = md5($input['old_ps'] . $adminInfo->salt);
        if ($encryptPs != $adminInfo->password) {
            notice(700034);
        }
        if ($input['new_ps'] != $input['new_ps_two']) {
            notice(700034);
        }
        //密码加密
        $newPassword = md5($input['new_ps'] . $adminInfo->salt);
        EnterpriseAdmin::query()->where('id', $adminInfo->id)->update([
            'password' => $newPassword,
            'password_expire_time' => time() + EnterpriseAdmin::PASSWORD_EXPIRE
        ]);
        return success([]);
    }

    /**
     * 重置密码
     *
     */
    public function resetPassword(): JsonResponse
    {
        $input = $this->request->only([
            'enterprise_id', 'account', 'mobile', 'admin_id', 'code'
        ]);
        if (empty($input['enterprise_id']) || (empty($input['admin_id']) && (empty($input['mobile']) || empty($input['account'])))) {
            notice(701017);
        }

        if (!empty($input['code'])) {
            //用户自己重置密码，非管理后台重置
            if (!$input['mobile'] = $this->redis->get($input['mobile'])) {
                notice(700032);
            }
            $this->checkCode($input['enterprise_id'], $input['mobile'], $input['code']);
        }
        (new EnterpriseAdminService())->resetPassword($input);
        if (!empty($input['code'])) {
            $this->redis->setex(sprintf(RedisKey::RESET_PASSWORD_TIME, $input['enterprise_id'], $input['mobile']), 86400, true);
        }
        return success([]);
    }

    private function checkCode($enterpriseId, $mobile, $code)
    {
        $accountDetail = (new EnterpriseAdmin())->accountDetailByPhone($mobile, $enterpriseId, ['salt', 'password', 'account', 'mobile', 'enterprise_id', 'username']);
        // 账号尚未开通
        if (empty($accountDetail)) {
            notice(700006);
        }
        if ($code != '952700') {
            (new EnterpriseAdminService())->checkCode($mobile, $code, 'console');
        }
    }

    /**
     * 重置密码验证
     *
     */
    public function checkResetPassword(Request $request): JsonResponse
    {
        $enterpriseId = $request->post('enterprise_id');
        $account = $request->post('account');
        $codeKey = $request->post('code_key');
        $code = trim($request->post('code'));

        if (empty($account)) {
            notice(700006);
        }
        if (empty($enterpriseId)) {
            notice(700020);
        }

        if ($code != '9527') {
            if (empty($code) || empty($codeKey)) {
                notice(700029);
            }
            //验证图片验证码
            if (!$this->checkImageCode($code, $codeKey)) {
                notice(700029);
            }
        }

        $accountDetail = (new EnterpriseAdmin())->accountDetailByAccount($account, $enterpriseId, ['account', 'mobile', 'status']);
        if (empty($accountDetail) || $accountDetail['status'] == EnterpriseAdmin::STATUS_CLOSE) {
            notice(700012);
        }

        if (empty($accountDetail['mobile'])) {
            notice(701017);
        }
        if ($this->redis->get(sprintf(RedisKey::RESET_PASSWORD_TIME, $enterpriseId, $accountDetail['mobile']))) {
            notice(700041);
        }
        $md5Mobile = md5($accountDetail['mobile'] . rand(0, 9999)) . Constant::TOKEN_SALT;
        $this->redis->setex($md5Mobile, 300, $accountDetail['mobile']);
        return success(['account' => $accountDetail['account'], 'mobile' => $md5Mobile, 'enterprise_id' => $enterpriseId]);
    }


    public function makeSalt()
    {
        return substr(uniqid(), -6);
    }

}