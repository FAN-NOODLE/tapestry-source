<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/21 19:23
 */

namespace App\Services;

use App\Models\ApplicationPower;
use App\Models\BasePower;
use App\Models\Role;
use App\Models\RolePermission;
use App\Models\RoleStore;

class RoleService extends BaseService
{

    /**
     * 全部权限
     */
    const POWER_ALL = 'all';

    /**
     * 应用权限
     */
    const POWER_APP = 'app';

    /**
     * 企业ID
     */
    private $enterpriseId = null;


    /**
     * @param $enterpriseId int 企业ID
     */
    public function __construct($enterpriseId = null)
    {
        parent::__construct();
        $this->enterpriseId = $enterpriseId;
    }

    /**
     * 获取角色权限
     *
     * @param $roleId int 角色ID
     * @return array
     */
    public function getPower($roleId)
    {
        $permission = RolePermission::query()->select([
            'id',
            'bap_id',
            'type'
        ])->where('role_id', $roleId)->get()->toArray();

        // 定义存放权限ID数组
        $permissionIdArr = [
            'base' => [],
            'app' => []
        ];

        // 提出权限ID
        foreach ($permission as $key => $value) {

            // 判断权限类型
            if ($value['type'] === RolePermission::TYPE_BASE) {
                $permissionIdArr['base'][] = $value['bap_id'];
            } else if ($value['type'] === RolePermission::TYPE_APP) {
                $permissionIdArr['app'][] = $value['bap_id'];
            }
        }

        // 查询应用权限
        $appAll = ApplicationPower::query()->select([
            'id',
            'name',
            'alias',
            'component',
            'path',
            'parent',
            'serve_type',
            'power_alias'
        ])->where([
            [
                'available',
                '=',
                ApplicationPower::AVAILABLE_IS
            ],
            [
                'serve_type',
                '!=',
                ApplicationPower::SERVE_TYPE_ALL
            ],
            [
                'enterprise_id',
                '=',
                $this->enterpriseId
            ]
        ])->get()->toArray();

        $appShopPower = $this->getAppPower($appAll, ApplicationPower::SERVE_TYPE_SHOP);
        $appCrsPower = $this->getAppPower($appAll, ApplicationPower::SERVE_TYPE_CRS);

        $app = [
            'shop' => $this->buildPower($appShopPower, $permissionIdArr['app']),
            'crs' => $this->buildPower($appCrsPower, $permissionIdArr['app'])
        ];

        $data = [
            'shop' => $app['shop'],
            'crs' => $app['crs']
        ];
        return $data;
    }

    /**
     * 获取全部权限（含应用、基础）
     *
     * @return array
     */
    public function getAllPower()
    {
        // 查询应用权限
        $appAll = ApplicationPower::query()->select([
            'id',
            'name',
            'alias',
            'component',
            'path',
            'parent',
            'serve_type',
            'power_alias'
        ])->where([
            [
                'available',
                '=',
                ApplicationPower::AVAILABLE_IS
            ],
            [
                'serve_type',
                '!=',
                ApplicationPower::SERVE_TYPE_ALL
            ],
            [
                'enterprise_id',
                '=',
                $this->enterpriseId
            ]
        ])->get()->toArray();

        $appShopPower = $this->getAppPower($appAll, ApplicationPower::SERVE_TYPE_SHOP);
        $appCrsPower = $this->getAppPower($appAll, ApplicationPower::SERVE_TYPE_CRS);

        $data = [
            'shop' => $appShopPower,
            'crs' => $appCrsPower
        ];
        return $data;
    }

    /**
     * 组装角色拥有的权限菜单
     *
     * @param $data array 全部权限数据
     * @param $hasPowerIdArr array 拥有的权限id
     * @return array
     */
    private function buildPower($data, $hasPowerIdArr)
    {
        $tree = [];
        foreach ($data as $key => $value) {
            $value['selected'] = in_array($value['id'], $hasPowerIdArr);

            // 存在子节点
            if (!empty($value['child'])) {
                $value['child'] = $this->buildPower($value['child'], $hasPowerIdArr);
            }
            $tree[] = $value;
            unset($data[$key]);
        }
        return $tree;
    }

    /**
     * 获取应用权限树形菜单
     *
     * @param $all array 全部权限数据
     * @param $serveType int BasePower权限类型SERVE_TYPE_ALL=全部,SERVE_TYPE_SHOP=shop,SERVE_TYPE_CRS=crs
     * @return array 树形权限结构
     */
    public function getAppPower($all, $serveType = ApplicationPower::SERVE_TYPE_ALL)
    {
        $data = [];
        foreach ($all as $key => $value) {
            if ($value['serve_type'] === $serveType) {
                $data[] = $value;
            }
        }
        return $this->getChildPower($data);
    }

    /**
     * 递归生成权限树形菜单
     *
     * @param $data array 所有权限数据
     * @param $path string 查找节点
     * @return array 树形权限结构
     */
    private function getChildPower($data, $path = '')
    {
        $tree = [];
        foreach ($data as $key => $value) {
            if ($value['parent'] === $path) {
                $value['child'] = $this->getChildPower($data, $value['path']);
                $tree[] = $value;
                unset($data[$key]);
            }
        }
        return $tree;
    }

}