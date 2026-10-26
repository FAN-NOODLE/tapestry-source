<?php

namespace App\Http\Controllers\App\V1\EnterpriseNew;

use App\Jobs\SendEmailJob;
use CrsSdk\MicroAuthority\Role;
use CrsSdk\Shop\Store;
use App\Exceptions\Error;
use App\Models\Enterprise;
use Illuminate\Http\Request;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseConfig;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use CrsSdk\MicroAuthority\Permissions;
use App\Services\EnterpriseConfigService;
use CrsSdk\Open\Enterprise as OpenEnterprise;
use Illuminate\Validation\ValidationException;
use App\Http\Requests\V1\EnterpriseNewRequest;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterpriseNew\ChildAccountService;
use App\Services\EnterpriseNew\EnterpriseNewService;

class EditController extends BaseController
{
    /**
     * 编辑企业
     * @param Request $request
     * @param EnterpriseNewService $enterpriseNewService
     * @return array
     * @throws ValidationException
     */
    public function index(Request $request, EnterpriseNewService $enterpriseNewService)
    {
        // 接收参数
        $params = $request->only([
            'enterprise_id',
            'name',
            'contact',
            'mobile',
            'password',
            'password_confirmation',
            'address',
            'corp_id',
            'version',
            'version_tag',
            'account_limit',
            'account_start_at',
            'account_expire_at',
            'available_system',
            'business_source',
            'permission_list',
            'warehouse_type',
            'industry_id',
        ]);

        // 请求数据验证
        EnterpriseNewRequest::edit($params);

        $enterprise = Enterprise::query()->where('id', '=', $params['enterprise_id'])->first();
        if (empty($enterprise)) {
            error(201108);
        }
        $eConfig = EnterpriseConfigService::getWarehouseConfig($enterprise->id);
        // 如果为中心商城,多仓变单仓不允许
        if ($params['version'] === Enterprise::VERSION_CENTER_SHOP && $eConfig['value'] == EnterpriseConfig::MORE_WAREHOUSE && $params['warehouse_type'] == EnterpriseConfig::SINGLE_WAREHOUSE) {
            error(201114);
        }

        if (!isset($params['permission_list'])) {
            $params['permission_list'] = [];
        }

        // 有门店不能为单仓
        if ($params['version'] != Enterprise::VERSION_CENTER_SHOP && $params['warehouse_type'] == EnterpriseConfig::SINGLE_WAREHOUSE) {
            error(201110);
        }

        // TODO::账号数量限制的处理等账号卖点启动后再介入开发

        $params['name']    = trim($params['name']);
        $params['contact'] = trim($params['contact']);
        $params['address'] = trim($params['address']);
        $params['corp_id'] = trim($params['corp_id']);
        // 封装最后操作时间
        $params['edit_at'] = time();
        $editVersionTag = $enterprise->version_tag;
        try {
            DB::beginTransaction();

            // 新版本企业编辑走new
            $enterprise = $enterpriseNewService->enterpriseEdit($params);

            if (!empty($enterprise)) {
                // 同步企业名字和hash给Open配置表
                try {
                    (new OpenEnterprise())->update(['name' => $enterprise['name'], 'enterprise_hash' => $enterprise['hash']]);
                } catch (\Exception $e) {
                    Log::error('编辑企业:同步企业信息:open接口异常', [
                        'line'    => $e->getLine(),
                        'code'    => $e->getCode(),
                        'file'    => $e->getFile(),
                        'message' => $e->getMessage(),
                    ]);
                }

                $paramsStore = [
                    'enterprise_id' => $params['enterprise_id'],
                    'name'          => $params['name'],
                    'is_store'      => 2
                ];

                switchAppHost($params['enterprise_id']);

                try {
                    (new Store())->add($paramsStore);
                } catch (\Exception $e) {
                    Log::error('编辑企业:同步企业到门店异常', [
                        'line'    => $e->getLine(),
                        'code'    => $e->getCode(),
                        'file'    => $e->getFile(),
                        'message' => $e->getMessage(),
                    ]);

                    throw $e;
                }
            }

            DB::commit();

        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error($e->getCode(), $e->getMessage());
        }

        // 刷新权限
        $role     = new Role();
        $roleInfo = $role->detail([
            'enterprise_id'  => $params['enterprise_id'],
            'role_mark'      => 'super_admin',
            'version_config' => -1,// -1只查基础信息
        ]);
        (new Role())->edit([
            'name'               => $roleInfo['name'],
            'role_id'            => $roleInfo['id'],
            'role_pid'           => 0,
            'enterprise_id'      => $params['enterprise_id'],
            'permission_id_list' => $params['permission_list']
        ]);

        // 发送邮件, 当修改了版本再发送邮件
        if ($editVersionTag != $params['version_tag']) {
            (new SendEmailJob($enterprise, 0))->handle();
        }
        return success(['data' => []]);
    }
}