<?php

namespace App\Models;

class Permissions extends Base
{
    /**
     * 类型:1菜单
     */
    const TYPE_MENU = 1;

    /**
     * 类型:2视图
     */
    const TYPE_VIEW = 2;

    /**
     * 类型:3按钮
     */
    const TYPE_BUTTON = 3;

    /**
     * 类型:4菜单+视图
     */
    const TYPE_MENU_AND_VIEW = 4;

    /**
     * 版本与企业表version字段关系
     */
    const VERSION_RELATION = [
        Enterprise::VERSION_CENTER_SHOP_STORE  => 'center_shop_store',
        Enterprise::VERSION_CENTER_SHOP        => 'center_shop',
        Enterprise::VERSION_DIRECT_CHAIN_STORE => 'direct_chain_store'
    ];

    /**
     * 依赖关系
     * key为依赖权限 values为依赖的子权限
     * 目的: 如果初始化角色没有key权限,则不会有values权限
     */
    const DEPEND_RELATION = [
        'myapp_organization' => ['myapp_organization_smart-guide'],
        // 企微'myapp_enterprise' => ['myapp_wisdom'],
    ];

    /**
     * 版本配置:1标配
     */
    const VERSION_STANDARD_CONFIG = 1;

    /**
     * 版本配置:2选配
     */
    const VERSION_SELECT_CONFIG = 2;

    /**
     * 版本配置:3无此功能
     */
    const VERSION_NO_CONFIG = 3;

    /**
     * 系统类型:SCRM
     */
    const SYSTEM_TYPE_SCRM = 11;

    /**
     * 系统类型:智慧导购isg
     */
    const SYSTEM_TYPE_ISG = 12;

    /**
     * 系统类型:商城新版store
     */
    const SYSTEM_TYPE_STORE = 13;

    /**
     * 系统类型:CDP
     */
    const SYSTEM_TYPE_CDP = 14;

    /**
     * 系统类型:商城旧版shop
     */
    const SYSTEM_TYPE_SHOP = 15;

    /**
     * 管理后台所有系统类型
     */
    const CONSOLE_SYSTEM_TYPE_RELATION = [
        self::SYSTEM_TYPE_SCRM  => 'SCRM',
        self::SYSTEM_TYPE_ISG   => '智慧导购ISG',
        self::SYSTEM_TYPE_STORE => '新商城store',
        self::SYSTEM_TYPE_CDP   => 'CDP',
        self::SYSTEM_TYPE_SHOP  => '旧商城shop'
    ];

    /**
     * 系统类型:智慧导购小程序
     */
    const SYSTEM_TYPE_GUIDE = 21;

    /**
     * 表名称:权限表
     *
     * @var string
     */
    protected $table = 'permissions';

    public $incrementing = true;

    protected $guarded = [];

    protected $hidden = [];

    /**
     * 通过类型和权限IDs集合去除不满足类型的ID
     *
     * @param $permissions
     * @param $type
     * @return array
     */
    public static function filterPermissionId($permissions, $type)
    {
        return self::query()
            ->select(['key', 'name'])
            ->whereIn('id', $permissions)
            ->where('system_type', '=', $type)
            ->get()
            ->toArray();
    }

    /**
     * 过滤指定仓库类型的权限 根据仓库类型找出不要的权限ID
     *
     * @param int $type 仓库类型
     * @return array
     */
    public static function filterWarehouse($type)
    {
        $idList = [];
        if ($type == 1) {
            $idList = Permissions::query()
                ->whereIn('name', ['仓库管理'])
                ->where('type', '=', Permissions::TYPE_MENU)
                ->pluck('id')
                ->toArray();
        }

        return $idList;
    }
}
