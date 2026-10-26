<?php
/**
 * Created by PhpStorm.
 * User: zhaoqi
 * Date: 2019-02-20
 * Time: 18:20
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterpriseAdminService;
use App\Models\EnterpriseAdmin;
use App\Http\Requests\EnterpriseAdminRequest;

class ResetPasswordController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * ResetPasswordController constructor.
     *
     * @param EnterpriseAdminService $enterpriseAdminService
     */
    public function __construct(EnterpriseAdminService $enterpriseAdminService)
    {
        parent::__construct();

        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 重置企业管理员密码
     *
     * @author 赵奇 zhaoqi@vchangyi.com
     * @date   2019/02/20 11:16
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'admin_id']);

        // 验证数据
        $validator = (new EnterpriseAdminRequest())->resetPassword($params);
        if ($validator->fails()) {
            return notice($validator->errors()->first(), [], $params);
        }

        // 根据企业ID查询管理员
        $admin = EnterpriseAdmin::query()
            ->where([
                'id'            => $params['admin_id'],
                'enterprise_id' => $params['enterprise_id']
            ])->where('role', '!=', EnterpriseAdmin::SUPER_ADMIN)
            ->first();

        if (!$admin) {
            return notice(202007, [], $params);
        }

        // 企业管理员密码重置
        $result = $this->enterpriseAdminService->enterpriseAdminPassword($params['admin_id']);
        if (!$result) {
            return notice(202005, [], $params);
        }

        return success(['data' => []]);
    }
}
