<?php

namespace App\Http\Controllers\V1\Account;

use App\Modules\BuiltInRole\BuiltInRoleActive\CommonRole;
use App\Modules\BuiltInRole\RoleStrategy;
use Illuminate\Validation\Rule;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class InitAccountRoleController extends BaseController
{
    const STORE_MEMBER = 'guide_store_member';
    const GUIDE_STORE_MANAGER = 'guide_store_manager';
    const CONSOLE_STORE_MANAGER = 'console_store_manager';
    const GUIDE_ENTERPRISE_GUIDE = 'guide_enterprise_guide';
    const CONSOLE_ENTERPRISE_GUIDE = 'console_enterprise_guide';
    const GUIDE_COMMUNITY_GUIDE = 'guide_community_guide';
    const CONSOLE_COMMUNITY_GUIDE = 'console_community_guide';
    const GUIDE_ENTERPRISE_ADMIN = 'guide_enterprise_admin';
    const CONSOLE_ENTERPRISE_ADMIN = 'console_enterprise_admin';

    /**
     * 店员
     * @throws ValidationException
     */
    public function initStoreMember()
    {
        $this->validate($this->request, [
            'type'          => ['required', Rule::in([CommonRole::ONLINE, CommonRole::OFFLINE])],
            'mobile'        => 'required',
            'store_id'      => 'required',
            'username'      => 'required',
            'operator_id'   => 'required',
            'enterprise_id' => 'required',
        ], [
            'type.required'          => '类型必传',
            'mobile.required'        => '手机号必传',
            'store_id.required'      => '门店ID必传',
            'username.required'      => '店员名称必传',
            'operator_id.required'   => '操作人必传',
            'enterprise_id.required' => '企业ID必传',
        ]);

        $params = $this->request->all();

        $strategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::STORE_MEMBER);

        $strategy->roleStrategy->accountOperation($params);

        return success();
    }

    /**
     * 店长
     * @throws ValidationException
     */
    public function initStoreManager()
    {
        $this->validate($this->request, [
            'mobile'        => 'required',
            'store_id'      => 'required',
            'username'      => 'required',
            'operator_id'   => 'required',
            'enterprise_id' => 'required',
        ], [
            'mobile.required'        => '手机号必传',
            'store_id.required'      => '门店ID必传',
            'username.required'      => '店员名称必传',
            'operator_id.required'   => '操作人必传',
            'enterprise_id.required' => '企业ID必传',
        ]);

        $params = $this->request->all();

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::GUIDE_STORE_MANAGER);

        $guideStrategy->roleStrategy->accountOperation($params);

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::CONSOLE_STORE_MANAGER);

        $guideStrategy->roleStrategy->accountOperation($params);

        return success();
    }

    /**
     * 企业导购
     */
    public function initEnterpriseGuide()
    {
        $this->validate($this->request, [
            'type'          => ['required', Rule::in([CommonRole::ONLINE, CommonRole::OFFLINE])],
            'mobile'        => 'required',
            'username'      => 'required',
            'operator_id'   => 'required',
            'enterprise_id' => 'required',
        ], [
            'mobile.required'        => '手机号必传',
            'store_id.required'      => '门店ID必传',
            'username.required'      => '店员名称必传',
            'operator_id.required'   => '操作人必传',
            'enterprise_id.required' => '企业ID必传',
        ]);

        $params = $this->request->all();

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::GUIDE_ENTERPRISE_GUIDE);

        $guideStrategy->roleStrategy->accountOperation($params);

        return success();
    }

    /**
     * 社群导购
     */
    public function initCommunityGuide()
    {
        $this->validate($this->request, [
            'type'          => ['required', Rule::in([CommonRole::ONLINE, CommonRole::OFFLINE])],
            'mobile'        => 'required',
            'username'      => 'required',
            'operator_id'   => 'required',
            'enterprise_id' => 'required',
        ], [
            'mobile.required'        => '手机号必传',
            'store_id.required'      => '门店ID必传',
            'username.required'      => '店员名称必传',
            'operator_id.required'   => '操作人必传',
            'enterprise_id.required' => '企业ID必传',
        ]);

        $params = $this->request->all();

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::GUIDE_COMMUNITY_GUIDE);

        $guideStrategy->roleStrategy->accountOperation($params);

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::CONSOLE_COMMUNITY_GUIDE);

        $guideStrategy->roleStrategy->accountOperation($params);

        return success();
    }

    /**
     * 企业管理
     */
    public function initEnterpriseAdmin()
    {
        $this->validate($this->request, [
            'type'          => ['required', Rule::in([CommonRole::ONLINE, CommonRole::OFFLINE])],
            'mobile'        => 'required',
            'username'      => 'required',
            'operator_id'   => 'required',
            'enterprise_id' => 'required',
        ], [
            'mobile.required'        => '手机号必传',
            'store_id.required'      => '门店ID必传',
            'username.required'      => '店员名称必传',
            'operator_id.required'   => '操作人必传',
            'enterprise_id.required' => '企业ID必传',
        ]);

        $params = $this->request->all();

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::GUIDE_ENTERPRISE_ADMIN);

        $guideStrategy->roleStrategy->accountOperation($params);

        $guideStrategy = (new RoleStrategy())->getRoleStrategy($params['enterprise_id'], $params['operator_id'], self::CONSOLE_ENTERPRISE_ADMIN);

        $guideStrategy->roleStrategy->accountOperation($params);

        return success();
    }
}