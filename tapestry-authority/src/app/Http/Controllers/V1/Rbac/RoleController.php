<?php

namespace App\Http\Controllers\V1\Rbac;

use App\Models\AdminRoleRelation;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseConfig;
use App\Models\Permissions;
use App\Modules\BuiltInRole\RoleStrategy;
use App\Modules\Rbac\Role;
use App\Models\Role as RoleModel;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\Rule;

class RoleController extends BaseController
{
    /**
     * @var Role
     */
    private $role;

    public function __construct()
    {
        parent::__construct();

        $this->validate($this->request, [
            'admin_id'      => 'integer|min:1',
            'enterprise_id' => 'required|integer|min:1'
        ]);
        $input = $this->request->all();

        // 为了兼容初始化角色 admin_id不传则为0
        $this->role = new Role($input['admin_id'] ?? 0, $input['enterprise_id']);
    }

    public function init()
    {
        $this->validate($this->request, [
            'admin_id'           => 'integer|min:1',
            'permission_id_list' => 'array',
            'mark_list'          => 'required|array'
        ]);
        $input = $this->request->all();

        // 过滤特殊权限
        $warehouseType = $this->getWarehouseType();

        $res = $this->role->init($input['mark_list'], $input['permission_id_list'] ?? [], $input['admin_id'] ?? 0, Permissions::filterWarehouse($warehouseType));

        // 内置管理员
        $guideStrategy = (new RoleStrategy())->getRoleStrategy($input['enterprise_id'], $input['admin_id'], 'console_admin_role');
        $guideStrategy->roleStrategy->accountOperation([]);

        return success($res);
    }

    public function add()
    {
        $this->validate($this->request, [
            'name'               => 'required|string',
            'role_pid'           => 'integer|min:0',
            'admin_id'           => 'required|integer|min:1',
            'weight'             => 'integer|min:0',
            'permission_id_list' => 'present|array'
        ]);
        $input = $this->request->all();

        $res = $this->role->add($input, $input['permission_id_list'], $input['admin_id'], $input['role_pid'] ?? RoleModel::SUPER_ADMIN_ROLE_PID);

        return success($res);
    }

    public function edit()
    {
        $this->validate($this->request, [
            'role_id'            => 'required|integer|min:1',
            'role_pid'           => 'required|integer|min:0',
            'name'               => 'required|string',
            'weight'             => 'integer|min:0',
            'permission_id_list' => 'present|array',
            'enterprise_id'      => 'integer'
        ]);
        $input = $this->request->all();

        // 过滤特殊权限
        $warehouseType = $this->getWarehouseType();

        $res = $this->role->edit($input['role_id'], $input, $input['permission_id_list'], Permissions::filterWarehouse($warehouseType));

        // 内置管理员
        Log::info('执行刷新管理员权限');
        $guideStrategy = (new RoleStrategy())->getRoleStrategy($input['enterprise_id'], $input['role_id'], 'console_admin_role');
        $guideStrategy->roleStrategy->accountOperation([]);
        Log::info('执行刷新完成');

        return success($res);
    }

    public function detail()
    {
        $this->validate($this->request, [
            'role_id'        => 'integer|min:1',
            'mark_list'      => [
                'string',
                Rule::in(RoleModel::ROLE['super_admin']['name'], RoleModel::ROLE['area']['name'], RoleModel::ROLE['store_manager']['name'], RoleModel::ROLE['store_member']['name'])
            ],
            'version_config' => [
                'integer',
                Rule::in(-1, 0, Permissions::VERSION_STANDARD_CONFIG, Permissions::VERSION_SELECT_CONFIG, Permissions::VERSION_NO_CONFIG)
            ]
        ]);
        $input = $this->request->all();

        if (empty($input['role_id']) && empty($input['role_mark'])) {
            notice(100000, ['请传入角色ID或标识']);
        }

        if (!empty($input['role_id'])) {
            $input['role_id_or_mark'] = $input['role_id'];
        } else {
            $input['role_id_or_mark'] = $input['role_mark'];
        }

        $data = $this->role->detail($input['role_id_or_mark'], $input['version_config'] ?? -1);

        // 获取企业关联的管理员数量
        $data['admin_name']           = EnterpriseAdmin::adminUsername($data['admin_id']);
        $data['admin_relation_count'] = AdminRoleRelation::adminRelationNums($data['id']);

        return success($data);
    }

    public function lists()
    {
        $this->validate($this->request, [
            'page'       => 'required|integer|min:1',
            'size'       => 'integer|max:100',
            'keywords'   => 'string',
            'sort_type'  => 'string|in:ASC,DESC',
            'role_name'  => 'string',
            'admin_name' => 'string'
        ]);
        $input = $this->request->all();

        if (!isset($input['size'])) {
            $input['size'] = 100;

        }
        if (!isset($input['sort_type'])) {
            $input['sort_type'] = 'DESC';
        }

        // 处理操作人查询
        $operator = [];
        if (!empty($input['admin_name'])) {
            $operator = EnterpriseAdmin::query()->where('username', 'like', "%{$input['admin_name']}%")->pluck('id')->toArray();
            $operator = empty($operator) ? [1] : $operator;
        }

        $input['operator'] = $operator;

        list($total, $list) = $this->role->lists($input['page'], $input['size'], $input);

        $names = EnterpriseAdmin::query()->whereIn('id', array_column($list, 'admin_id'))->pluck('username', 'id')->toArray();
        foreach ($list as &$item) {
            $item['admin_name'] = $names[$item['admin_id']] ?? '系统';
        }

        return success([
            'list'  => $list,
            'total' => $total,
            'page'  => $input['page'],
            'size'  => $input['size']
        ]);
    }

    public function delete()
    {
        $this->validate($this->request, [
            'role_id' => 'required|integer|min:1'
        ]);
        $input = $this->request->all();

        $this->role->delete($input['role_id']);

        return success();
    }

    private function getWarehouseType()
    {
        return EnterpriseConfig::query()
            ->where('enterprise_id', '=', $this->request->input('enterprise_id'))
            ->where('key', '=', 'warehouse_type')
            ->value('value');
    }
}
