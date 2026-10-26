<?php

namespace App\Http\Controllers\App\V1\EnterpriseNew;

use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use App\Http\Requests\V1\EnterpriseNewRequest;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterpriseNew\EnterpriseNewService;

class AddController extends BaseController
{
    /**
     * 添加企业
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws ValidationException
     */
    public function index()
    {
        // 接收参数
        $params = $this->request->only([
            'name',
            'contact',
            'mobile',
            'address',
            'corp_id',
            'version',
            'version_tag',
            'account_limit',
            'account_start_at',
            'account_expire_at',
            'email',
            'password',
            'password_confirmation',
            'available_system',
            'permission_list',
            'warehouse_type',
            'industry_id',
            'service_merchant_id',
        ]);
        // 请求数据验证
        EnterpriseNewRequest::add($params);

        $enterprise = ((new EnterpriseNewService()))->add($params);

        return success(['data' => $enterprise]);
    }
}