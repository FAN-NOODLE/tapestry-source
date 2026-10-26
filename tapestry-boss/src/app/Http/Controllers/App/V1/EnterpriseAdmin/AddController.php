<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\EnterpriseAdmin;

use App\Exceptions\Error;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use App\Http\Requests\EnterpriseAdminRequest;
use App\Services\EnterpriseAdminService;

class AddController extends BaseController
{
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
     * @param EnterpriseAdminRequest $enterpriseAdminRequest
     * @param EnterpriseAdminService $enterpriseAdminService
     * @throws \Exception
     */
    public function __construct(
        EnterpriseAdminRequest $enterpriseAdminRequest,
        EnterpriseAdminService $enterpriseAdminService
    )
    {
        parent::__construct();

        // 企业管理员验证器赋值
        $this->enterpriseAdminRequest = $enterpriseAdminRequest;
        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 增加企业管理员
     *
     * @param Request $request
     *
     * @return array|mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/25 11:16
     *
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'username', 'account', 'password']);

        // 请求数据验证
        $validator = $this->enterpriseAdminRequest->common($params);
        if ($validator->fails()) {
            throw new Error(202200, $validator->errors()->first());
        }

        // 企业管理员信息添加
        $params['salt']     = $this->enterpriseAdminService->makeSalt();
        $params['password'] = md5($params['password'] . $params['salt']);
        $result             = $this->enterpriseAdminService->adminAdd($params);
        if (!$result) {
            return error(202001);
        }

        return success(['data' => []]);
    }
}
