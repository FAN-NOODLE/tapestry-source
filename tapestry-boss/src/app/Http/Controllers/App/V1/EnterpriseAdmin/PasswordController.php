<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\EnterpriseAdmin;

use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterpriseAdminService;
use App\Models\EnterpriseAdmin;

class PasswordController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * PasswordController constructor.
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
     * @author 梁欢 lianghuan@vchangyi.com
     * @author 赵奇
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id']);

        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return notice(201101);
        }

        // 根据企业ID查询管理员
        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('role', EnterpriseAdmin::SUPER_ADMIN)
            ->first(['id']);

        if (!$admin) {
            return notice(202007);
        }

        $admin = $admin->toArray();
        // 企业管理员密码重置
        if (!$this->enterpriseAdminService->enterpriseAdminPassword($admin['id'])) {
            return notice(202005);
        }

        return success(['data' => []]);
    }
}
