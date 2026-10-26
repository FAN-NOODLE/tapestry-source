<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\EnterpriseAdmin;

use App\Models\EnterpriseAdmin;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterpriseAdminService;

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
     * @param EnterpriseAdminService $enterpriseAdminService
     * @throws \Exception
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
     * @param Request $request
     *
     * @return array|mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/25 11:16
     *
     */
    public function index(Request $request)
    {
        $params = $request->only(['admin_id', 'username', 'account', 'status']);

        // 企业管理员ID不能为空
        if (!isset($params['admin_id']) || empty($params['admin_id'])) {
            return error(202101);
        }

        // 查询管理员信息
        $admin = $this->enterpriseAdminService->enterpriseAdmin($params['admin_id']);
        if (!$admin) {
            return error(202002);
        }

        // 更新管理员信息
        $result = EnterpriseAdmin::query()->where([
            'id' => $params['admin_id']
        ])->update($params);

        if (!$result) {
            return error(202005);
        }

        return success(['data' => []]);
    }
}
