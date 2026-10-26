<?php

namespace App\Http\Controllers\App\V1\Enterprise;

use Exception;
use App\Exceptions\Error;
use CrsSdk\Shop\UserLevel;
use CrsSdk\Open\Enterprise;
use CrsSdk\Shop\PagePostage;
use Illuminate\Http\Request;
use App\Models\EnterpriseApp;
use CrsSdk\Console\Application;
use CrsSdk\Work\CommissionSetting;
use Illuminate\Support\Facades\DB;
use App\Services\EnterpriseService;
use Illuminate\Support\Facades\Log;
use App\Models\EnterpriseAppInstall;
use App\Services\EnterpriseLogService;
use App\Services\EnterpriseAdminService;
use App\Http\Requests\EnterpriseRequest;
use App\Http\Requests\EnterpriseAdminRequest;
use App\Http\Controllers\App\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 企业验证器
     *
     * @var EnterpriseRequest
     */
    protected $enterpriseRequest;

    /**
     * 企业service
     *
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * 企业管理员验证器
     *
     * @var EnterpriseAdminRequest
     */
    protected $enterpriseAdminRequest;

    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * AddController constructor.
     * @param EnterpriseRequest $enterpriseRequest
     * @param EnterpriseService $enterpriseService
     * @param EnterpriseAdminRequest $enterpriseAdminRequest
     * @param EnterpriseAdminService $enterpriseAdminService
     * @throws Exception
     */
    public function __construct(
        EnterpriseRequest $enterpriseRequest,
        EnterpriseService $enterpriseService,
        EnterpriseAdminRequest $enterpriseAdminRequest,
        EnterpriseAdminService $enterpriseAdminService
    )
    {
        parent::__construct();

        // 企业验证器赋值
        $this->enterpriseRequest = $enterpriseRequest;
        // 企业service赋值
        $this->enterpriseService = $enterpriseService;
        // 企业管理员验证器赋值
        $this->enterpriseAdminRequest = $enterpriseAdminRequest;
        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 增加企业
     *
     * @param Request $request
     * @return array|mixed
     * @throws Exception
     * @author 贺树鹏 heshupeng@vchangyi.com
     * @date   2019/02/16 38:16
     */
    public function index(Request $request)
    {
        // 接收参数
        // $params = $request->only(['name', 'contact', 'mobile', 'address', 'crs_status', 'store_status', 'corp_id', 'username', 'account', 'app_ids']);

        $input = $request->all();

        $input['name']    = trim($input['name']);
        $input['contact'] = trim($input['contact']);
        $input['address'] = trim($input['address']);

        // 请求数据验证
        $validator = $this->enterpriseRequest->common($input);
        if ($validator->fails()) {
            throw new Error(201100, $validator->errors()->first());
        }

        try {
            $this->enterpriseService->add($input, $input['account']);
        } catch (Exception $exception) {
            throw new Exception('添加企业失败:'.$exception->getMessage());
        }

        return success(['data' => []]);
    }
}
