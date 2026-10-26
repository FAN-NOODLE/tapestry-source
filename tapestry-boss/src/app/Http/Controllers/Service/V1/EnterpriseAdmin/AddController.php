<?php
/**
 * Created by PhpStorm.
 * User: zhaoqi
 * Date: 2019/02/20
 * Time: 15:26
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
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
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/20 15:26
     *
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'username', 'account', 'password', 'operator_id', 'role_id']);

        $params['username'] = trim($params['username']);
        $params['account']  = trim($params['account']);

        // 兼容门店添加店长时，创建默认管理员
        if (empty($params['password'])) {
            $params['password'] = md5(EnterpriseAdminService::PASSWORD);
        }

        // 请求数据验证
        $validator = $this->enterpriseAdminRequest->add($params);
        if ($validator->fails()) {
            return notice($validator->errors()->first(), [], $params);
        }

        // 企业管理员信息添加
        $params['salt']     = $this->enterpriseAdminService->makeSalt();
        $params['password'] = md5($params['password'] . $params['salt']);
        $params['role']     = $params['role_id'];
        unset($params['role_id']);

        $result = $this->enterpriseAdminService->adminAdd($params);

        if (!$result) {
            return notice(202001, [], $params);
        }

        return success(['data' => ['id' => $result->id]]);
    }
}
