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
use App\Http\Resources\EnterpriseAdminResource;

class DetailController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * DetailController constructor.
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
     * 企业管理员详情
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['admin_id']);

        // 企业管理员ID不能为空
        if (!isset($params['admin_id']) || empty($params['admin_id'])) {
            return error(202101);
        }

        // 企业管理员信息查询
        $enterprise = $this->enterpriseAdminService->enterpriseAdmin($params['admin_id']);
        if (!$enterprise) {
            return error(202002);
        }

        return success(['data' => new EnterpriseAdminResource($enterprise)]);
    }
}
