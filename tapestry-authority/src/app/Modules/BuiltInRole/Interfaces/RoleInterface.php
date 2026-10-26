<?php

namespace App\Modules\BuiltInRole\Interfaces;

interface RoleInterface
{
    // 初始化角色
    public function initRole();

    // 获取数据权限
    public function getDataRange($adminId, $type, $objectId);

    // 绑定数据权限
    public function adminBindDataRange($adminId, $storeId, $type);

    // 解绑数据权限关联
    public function unsetDataRangeRelation($adminId, $type, $storeOrderShop);

    // 获取角色
    public function getRole($adminId, $roleId);

    // 绑定角色
    public function adminBindRole($storeAdminId);

    // 解绑角色关联
    public function unsetRoleRelation($adminId);

    // 账号操作维护
    public function accountOperation($params);

    // 管理员角色互斥
    public function adminRoleMutex($admin);

    // 管理员状态启用
    public function adminStatusEnable($admin);

    // 获取权限
    public function getPermission();
}