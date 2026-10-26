<?php

namespace App\Http\Controllers\V1\Rbac;

use App\Http\Controllers\V1\BaseController;
use App\Models\Permissions as PermissionsModel;
use App\Modules\Rbac\Permissions;
use App\Services\PermissionReleaseService;

class PermissionsController extends BaseController
{
    /**
     * @var Permissions
     */
    private $permissions;

    public function __construct()
    {
        parent::__construct();

        $this->validate($this->request, [
            'version' => 'integer'
        ]);
        $input = $this->request->all();
        if (isset($input['version']) === false) {
            $this->request->offsetSet('version', 0);
            $input['version'] = 0;
        }

        $this->permissions = new Permissions($input['version']);
    }

    /**
     * 所有权限树形列表
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function allTree()
    {
        $tree = $this->permissions->allTree();

        return success($tree);
    }

    /**
     * 附加功能列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function additionalList()
    {
        $tree = $this->permissions->additionalConfigList();

        return success($tree);
    }

    /**
     * 新增权限
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function add()
    {
        $this->validate($this->request, [
            'pid'                => 'required|integer',
            'name'               => 'required|string',
            'path'               => 'required|string',
            'key'                => 'required|string',
            'icon'               => 'present|string',
            'type'               => 'required|integer',
            'system_type'        => 'required|integer|in:' . implode(',', [
                    PermissionsModel::SYSTEM_TYPE_SCRM,
                    PermissionsModel::SYSTEM_TYPE_ISG,
                    PermissionsModel::SYSTEM_TYPE_STORE,
                    PermissionsModel::SYSTEM_TYPE_CDP,
                    PermissionsModel::SYSTEM_TYPE_SHOP,
                ]),
            'weight'             => 'integer|min:0',
            'center_shop_store'  => 'required|integer|in:1,2,3',
            'center_shop'        => 'required|integer|in:1,2,3',
            'direct_chain_store' => 'required|integer|in:1,2,3',
            'version_tag'        => 'required|string',
            'is_full_release'    => 'required|integer|in:0,1'
        ]);
        $input = $this->request->all();

        // key唯一
        if (PermissionsModel::query()->where('key', '=', $input['key'])->exists()) {
            notice(100000, ['key已存在']);
        }

        $data = [
            'pid'                => $input['pid'],
            'name'               => $input['name'],
            'path'               => $input['path'],
            'key'                => $input['key'],
            'icon'               => $input['icon'] ?? '',
            'type'               => $input['type'],
            'system_type'        => $input['system_type'],
            'weight'             => $input['weight'] ?? 0,
            'center_shop_store'  => $input['center_shop_store'],
            'center_shop'        => $input['center_shop'],
            'direct_chain_store' => $input['direct_chain_store'],
            'version_tag'        => $input['version_tag'],
            'is_full_release'    => $input['is_full_release']
        ];

        // 兼容其他环境调用同步权限,保持权限ID一致
        if (isset($input['id'])) {
            $data['id'] = $input['id'];
        }
        PermissionsModel::query()->create($data);

        return success();
    }

    /**
     * 编辑权限接口
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function edit()
    {
        $this->validate($this->request, [
            'permissions_id'     => 'required|integer',
            'pid'                => 'required|integer',
            'name'               => 'required|string',
            'path'               => 'required|string',
            'key'                => 'required|string',
            'icon'               => 'present|string',
            'type'               => 'required|integer',
            'system_type'        => 'required|integer|in:' . implode(',', [
                    PermissionsModel::SYSTEM_TYPE_SCRM,
                    PermissionsModel::SYSTEM_TYPE_ISG,
                    PermissionsModel::SYSTEM_TYPE_STORE,
                    PermissionsModel::SYSTEM_TYPE_CDP,
                    PermissionsModel::SYSTEM_TYPE_SHOP,
                ]),
            'weight'             => 'integer|min:0',
            'center_shop_store'  => 'required|integer|in:1,2,3',
            'center_shop'        => 'required|integer|in:1,2,3',
            'direct_chain_store' => 'required|integer|in:1,2,3',
            'version_tag'        => 'required|string',
            'is_full_release'    => 'required|integer|in:0,1'
        ]);
        $input = $this->request->all();

        $permissions = PermissionsModel::query()->where('id', '=', $input['permissions_id'])->first();
        if ($permissions === null) {
            notice(100000, ['权限不存在']);
        }

        // key唯一
        if (PermissionsModel::query()
            ->where('id', '!=', $permissions->getAttribute('id'))
            ->where('key', '=', $input['key'])
            ->exists()) {
            notice(100000, ['key已存在']);
        }

        // 防止当前权限放入自己的子级下
        $childrenList = $this->permissions->getChildren($this->permissions->all(), $permissions->getAttribute('id'));
        if (in_array($input['pid'], array_column($childrenList, 'id'))) {
            notice(100000, ['父级禁止归属于子级下']);
        }

        $permissions->update([
            'pid'                => $input['pid'],
            'name'               => $input['name'],
            'path'               => $input['path'],
            'key'                => $input['key'],
            'icon'               => $input['icon'] ?? '',
            'type'               => $input['type'],
            'system_type'        => $input['system_type'],
            'weight'             => $input['weight'] ?? 0,
            'center_shop_store'  => $input['center_shop_store'],
            'center_shop'        => $input['center_shop'],
            'direct_chain_store' => $input['direct_chain_store'],
            'version_tag'        => $input['version_tag'],
            'is_full_release'    => $input['is_full_release']
        ]);

        return success();
    }

    /**
     * 权限列表接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function lists()
    {
        $this->validate($this->request, [
            'system_type' => 'integer|min:0',
            'version_tag' => 'string'
        ]);
        $input = $this->request->all();

        // 默认不需要软删除数据
        if (!isset($input['need_deleted'])) {
            $input['need_deleted'] = 0;
        }

        $data = PermissionsModel::withTrashed()
            ->when(isset($input['system_type']) && $input['system_type'] !== 0, function ($query) use ($input) {
                return $query->where('system_type', '=', $input['system_type']);
            })
            ->when(isset($input['version_tag']) && $input['version_tag'] !== '', function ($query) use ($input) {
                return $query->where('version_tag', '=', $input['version_tag']);
            })
            ->when(isset($input['need_deleted']) && $input['need_deleted'] === 0, function ($query) {
                return $query->where('deleted_at', '=', 0);
            })
            ->get()
            ->toArray();

        return success([
            'list' => $data
        ]);
    }

    /**
     * 发布权限接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function release()
    {
        $this->validate($this->request, [
            'system_type'   => 'required|integer',
            'version_tag'   => 'required|string',
            'enterprise_id' => 'required|integer|min:0'
        ]);
        $input = $this->request->all();

        $prs = new PermissionReleaseService($input['system_type'], $input['version_tag'], $input['enterprise_id']);
        $prs->enterpriseAdminRelationList();
        $result = $prs->release();

        return success([
            'result' => $result
        ]);
    }

    /**
     * 系统类型列表
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function systemTypeList()
    {
        $list = [];
        foreach (PermissionsModel::CONSOLE_SYSTEM_TYPE_RELATION as $type => $desc) {
            $list[] = [
                'value' => $type,
                'name'  => $desc
            ];
        }

        return success([
            'list' => $list
        ]);
    }

    /**
     * 同步权限 list接口返回的数据作为请求参数
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function sync()
    {
        $this->validate($this->request, [
            'list' => 'required|array'
        ]);
        $input = $this->request->all();

        $insertData = [];
        foreach ($input['list'] as $permission) {
            $permissionObj = PermissionsModel::withTrashed()->where('id', '=', $permission['id'])->first();
            if ($permissionObj === null) {
                $insertData[] = $permission;
            } else {
                unset($permission['id']);
                $permissionObj->update($permission);
            }
        }
        PermissionsModel::query()->insert($insertData);

        return success();
    }

    /**
     * 子级菜单ID列表
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function childrenIDList()
    {
        $this->validate($this->request, [
            'id' => 'required|integer|min:0'
        ]);
        $pid = $this->request->input('id');

        $data = $this->permissions->getChildrenIDList($pid);

        return success($data);
    }
}
