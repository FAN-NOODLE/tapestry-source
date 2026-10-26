<?php

namespace App\Http\Controllers\V1\Account;

use App\Modules\Rbac\Role;
use Illuminate\Http\Request;
use App\Modules\Rbac\AdminRole;
use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Crypt;

class GetTokenController extends BaseController
{
    public function index(Request $request)
    {
        $adminId       = trim($request->post('admin_id', 0));
        $enterpriseId  = trim($request->post('enterprise_id', 0));
        $client        = trim($request->post('client'));
        $ttl           = trim($request->post('ttl', EnterpriseAdminService::TOKEN_EXPIRE));

        if (empty($client)) {
            notice(700013);
        }
        if (empty($enterpriseId)) {
            notice(700020);
        }

        // 该账号对应的企业列表
        $adminInfo = (new EnterpriseAdminService())->enterpriseAdminDetail($adminId, $client);

        if (empty($adminInfo)) {
            notice(700012);
        }
        if ((string)$adminInfo['enterprise_id'] !== (string)$enterpriseId) {
            notice(700012);
        }

        if ($adminInfo['account_status'] == EnterpriseAdminService::ACCOUNT_STATUS_BAN) {
            notice(700015);
        }

        $token = (new EnterpriseAdminService())->generateToken($adminInfo, false, $ttl);

        $adminInfo['enterprise_role_list'] = (new EnterpriseAdminService())->enterpriseAdminRoleList($adminId);

        if ($adminInfo['version'] > 1) {
            $role        = new Role($adminInfo['admin_id'], $adminInfo['enterprise_id']);
            $newRoleList = (new AdminRole($adminInfo['admin_id']))->roles();
        }

        $permissions_list = isset($role) ? $role->permissions->lists() : [];
        // 缓存权限api_list
        $api_list = [];
        foreach ($permissions_list as $val) {
            if (!empty($val['api_list'])) {
                $api_list = array_merge($api_list, json_decode($val['api_list']));
            }
        }
        $this->redis->set($adminInfo['admin_id'], json_encode($api_list));

        return success([
            'token'                 => $token,
            'tokenBytes'            => json_decode(Crypt::decrypt($token), true),
            'role_list'             => $newRoleList ?? [],
            'permissions_tree'      => isset($role) ? $role->permissions->myTree() : [],
            'permissions_list'      => $permissions_list,
            'enterprise_admin_info' => $adminInfo,
        ]);
    }
}
