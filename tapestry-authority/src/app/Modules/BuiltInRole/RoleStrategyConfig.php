<?php

namespace App\Modules\BuiltInRole;

class RoleStrategyConfig
{
    /**
     * 店员级别权限
     */
    const ROLE = [
        'guide_store_member'       => [
            'name'                  => '智慧导购店员',
            'mark'                  => 'guide_store_member',
            'key_list'              => [
                'selectEnterprise',
                'selectStore',
                'storeQrcode',
                'setting',
                'index',
                'orderList',
                'orderListPromotion',
                'orderListService',
                'customerList',
                'promotionBanner',
                'serviceBanner',
            ],
            'namespace'             => 'GuideStoreMemberRole',
            'permission'            => [231, 232, 233, 234, 235, 236, 237, 238, 240, 245, 246],
            'role_available_system' => [2]
        ],
        'guide_store_manager'      => [
            'name'                  => '智慧导购店长',
            'mark'                  => 'guide_store_manager',
            'key_list'              => [
                'selectEnterprise',
                'selectStore',
                'storeQrcode',
                'setting',
                'index',
                'orderList',
                'orderListPromotion',
                'orderListService',
                'shopMemberManage',
                'customerList',
                'compass',
                'dataCompass',
                'customerData',
                'banner',
                'promotionBanner',
                'serviceBanner',
            ],
            'namespace'             => 'GuideStoreManagerRole',
            'permission'            => [231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246],
            'role_available_system' => [2]
        ],
        'console_store_manager'    => [
            'name'                  => '管理后台店长',
            'mark'                  => 'console_store_manager',
            'key_list'              => [
                'myapp_home',
                'myapp_store_store-manager',
                'myapp_store_store-manager_goods-decorate-edit',
                'myapp_store_store-manager_commodity-management',
                'myapp_store_store-manager_store-detail',
                'myapp_store_store-manager_store-decorate',
                'myapp_store_store-manager_performance-detail',
                'myapp_store_store-manager_performance-list',
                'myapp_store_store-manager_clerk-list',
                'myapp_store_store-manager_store-edit',
                'myapp_store_store-manager_store-list',
                'myapp_store',
                'myapp_order-manager_order',
                'myapp_order-manager_order_detail',
                'myapp_order-manager_order_list',
                'myapp_order-manager',
                'myapp_order-manager_after-sale-manage',
                'myapp_order-manager_after-sale-manage_list',
            ],
            'namespace'             => 'GuideStoreManagerRole',
            'permission'            => [1, 111, 115, 116, 112, 114, 113, 121, 125, 128, 133, 134, 130, 129, 131, 132, 127, 126, 117, 118, 119, 120, 122, 123, 124],
            'role_available_system' => [1]
        ],
        'guide_enterprise_guide'   => [
            'name'                  => '智慧导购企业导购',
            'mark'                  => 'guide_enterprise_guide',
            'key_list'              => [
                'selectEnterprise',
                'setting',
                'index',
                'orderList',
                'orderListPromotion',
                'customerList',
                'banner',
                'promotionBanner',
            ],
            'namespace'             => 'GuideEnterpriseGuideRole',
            'permission'            => [231, 234, 235, 236, 237, 240, 244, 245],
            'role_available_system' => [2]
        ],
        'console_enterprise_guide' => [
            'name'                  => '管理后台企业导购',
            'mark'                  => 'console_enterprise_guide',
            'key_list'              => ['myapp_home'],
            'namespace'             => 'ConsoleEnterpriseGuideRole',
            'permission'            => [1],
            'role_available_system' => [1]
        ],
        'guide_community_guide'    => [
            'name'                  => '智慧导购社群导购',
            'mark'                  => 'guide_community_guide',
            'key_list'              => [
                'selectEnterprise',
                'setting',
                'index',
                'orderList',
                'orderListPromotion',
                'customerList',
                'banner',
                'promotionBanner',
            ],
            'namespace'             => 'GuideCommunityGuideRole',
            'permission'            => [231, 234, 235, 236, 237, 240, 244, 245],
            'role_available_system' => [2]
        ],
        'console_community_guide'  => [
            'name'                  => '管理后台社群导购',
            'mark'                  => 'console_community_guide',
            'key_list'              => ['myapp_home'],
            'namespace'             => 'ConsoleCommunityGuideRole',
            'permission'            => [1],
            'role_available_system' => [1]
        ],
        'guide_enterprise_admin'   => [
            'name'                  => '智慧导购企业管理',
            'mark'                  => 'guide_enterprise_admin',
            'key_list'              => [
                'selectEnterprise',
                'setting',
                'compass',
                'dataCompass',
                'customerData',
                'banner',
                'promotionBanner',
                'serviceBanner',
            ],
            'namespace'             => 'GuideEnterpriseAdmin',
            'permission'            => [231, 234, 241, 242, 243, 244, 245, 246],
            'role_available_system' => [2]
        ],
        'console_enterprise_admin' => [
            'name'                  => '管理后台企业管理',
            'mark'                  => 'console_enterprise_admin',
            'key_list'              => [
                'myapp_home',
                'myapp_smart-guide',
                'myapp_smart-guide_data-compass',
                'myapp_smart-guide_guide-manage',
                'myapp_smart-guide_guide-manage_guide-list',
                'myapp_smart-guide_guide-manage_guide-list_list',
                'myapp_smart-guide_guide-manage_guide-set',
                'myapp_smart-guide_guide-manage_guide-set_list',
                'myapp_smart-guide_guide-manage_guide-pay',
                'myapp_smart-guide_guide-manage_guide-pay_list',
            ],
            'namespace'             => 'ConsoleEnterpriseAdminRole',
            'permission'            => [
                247, 289, 248, 290, 291, 292, 293, 294, 295
            ],
            'role_available_system' => [1]
        ],
        'console_admin_role' => [
            'name'                  => '管理员',
            'mark'                  => 'console_admin_role',
            'key_list'              => [],
            'namespace'             => 'ConsoleAdminRole',
            'permission'            => [],
            'role_available_system' => []
        ]
    ];

    /**
     * 获取角色配置
     * @param $mark
     * @return mixed|void
     */
    public static function getConfig($mark)
    {
        $keys = array_keys(self::ROLE);

        if (!in_array($mark, $keys)) {
            error(702023);
        }

        return self::ROLE[$mark];
    }
}