<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use App\Exceptions\Error;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\EnterpriseAdminRequest;
use App\Services\EnterpriseAdminService;
use App\Models\EnterpriseAdmin;

class PasswordController extends BaseController
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
     *
     * @param EnterpriseAdminRequest $enterpriseAdminRequest
     * @param EnterpriseAdminService $enterpriseAdminService
     */
    public function __construct(
        EnterpriseAdminRequest $enterpriseAdminRequest,
        EnterpriseAdminService $enterpriseAdminService
    ) {
        parent::__construct();

        // 企业管理员验证器赋值
        $this->enterpriseAdminRequest = $enterpriseAdminRequest;
        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 企业管理员
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'admin_id', 'old_password', 'password']);

        // 请求数据验证
        $validator = $this->enterpriseAdminRequest->passwordUpdate($params);
        if ($validator->fails()) {
            throw new Error(202200, $validator->errors()->first());
        }

        // 检查是否是该企业下管理员
        $admin = EnterpriseAdmin::query()
            ->where([
                'id'            => $params['admin_id'],
                'enterprise_id' => $params['enterprise_id']
            ])->first();

        if (!$admin) {
            return error(202007);
        }

        // 修改密码
        $result = $this->enterpriseAdminService
            ->passwordUpdate(
                $params['admin_id'],
                $params['old_password'],
                $params['password']
            );

        if (!$result) {
            return error(202106);
        }

        return success(['data' => []]);
    }
}
