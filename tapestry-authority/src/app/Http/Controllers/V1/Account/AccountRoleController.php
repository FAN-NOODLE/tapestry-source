<?php

namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use App\Models\AdminRoleRelation;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Models\Permissions;
use App\Models\Role;
use App\Modules\BuiltInRole\BuiltInRoleActive\CommonRole;
use App\Modules\BuiltInRole\RoleStrategy;
use Illuminate\Validation\Rule;

class AccountRoleController extends BaseController
{
    // 店长
    const IS_BOSS = 1;

    const STORE_MEMBER = 'guide_store_member';
    const GUIDE_STORE_MANAGER = 'guide_store_manager';

    public function accountByRole()
    {
        $params = $this->request->all();

        if (empty($params['enterprise_id'])) {
            error(700014);
        }

        if (empty($params['role_ids']) && empty($params['marks'])) {
            error(702027);
        }

        // 如果标识不为空, 以标识为主, 即使传递了角色ID, 也会被覆盖
        if (!empty($params['marks'])) {
            $params['role_ids'] = Role::roleByMark($params['marks'], $params['enterprise_id']);
        }

        // 使用角色ID集合获取管理员ID集合
        $adminIds = [];
        if (!empty($params['role_ids'])) {
            $adminIds = AdminRoleRelation::adminRelationIds($params['role_ids']);
        }

        $adminList = EnterpriseAdmin::adminListByIds($params['enterprise_id'], $adminIds);

        return success([
            'list' => $adminList->toArray()
        ]);
    }

    public function storeWorkerPermission()
    {
        $this->validate($this->request, [
            'is_boss'  => 'required',
            'store_id' => 'required',
            'admin_id' => 'required',
        ], [
            'is_boss.required'  => '是否店长必传',
            'store_id.required' => '门店ID必传',
            'admin_id.required' => '管理员ID必传',
        ]);

        $params = $this->request->all();

        $dataRange = EnterpriseDataRange::getDataRange($params['admin_id'], $params['store_id'], EnterpriseDataRange::TYPE_STORE);
        if (empty($dataRange)) {
            error(702028);
        }

        // 是店长
        if ($params['is_boss'] === self::IS_BOSS) {
            $permissions = RoleStrategy::getPermission($params['admin_id'], Permissions::SYSTEM_TYPE_GUIDE, [self::GUIDE_STORE_MANAGER]);
        } else {
            $permissions = RoleStrategy::getPermission($params['admin_id'], Permissions::SYSTEM_TYPE_GUIDE, [self::STORE_MEMBER]);
        }

        return success($permissions);
    }
}