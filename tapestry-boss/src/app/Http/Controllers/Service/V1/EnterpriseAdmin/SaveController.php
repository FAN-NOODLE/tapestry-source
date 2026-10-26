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
use App\Services\EnterpriseAdminService;
use App\Models\EnterpriseAdmin;
use App\Http\Requests\EnterpriseAdminRequest;

class SaveController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * SaveController constructor.
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
     * 编辑企业管理员
     *
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/20 15:26
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'admin_id', 'username', 'account', 'operator_id', 'role_id', 'member_id']);

        $params['username'] = trim($params['username']);
        $params['account'] = trim($params['account']);

        // 验证数据
        $validator = (new EnterpriseAdminRequest())->save($params);
        if ($validator->fails()) {
            return notice($validator->errors()->first(), [], $params);
        }

        // 检查是否是该企业下管理员
        $admin = EnterpriseAdmin::query()
            ->where([
                'id'            => $params['admin_id'],
                'enterprise_id' => $params['enterprise_id']
            ])->first();

        if (!$admin) {
            return notice(202007, [], $params);
        }

        if (EnterpriseAdmin::SUPER_ADMIN == $admin->role) {
            // return notice(202007, [], $params);
        }

        // 企业管理员信息参数处理
        $adminId = $params['admin_id'];
        unset($params['admin_id']);
        if (empty($params)) {
            return notice(202004, [], $params);
        }

        $params['role'] = $params['role_id'];
        unset($params['role_id']);

        // 企业管理员信息保存
        $result = $this->enterpriseAdminService->enterpriseAdminSave($adminId, $params);
        if (!$result) {
            return notice(202005, [], $params);
        }

        return success(['data' => []]);
    }
}
