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

class ListController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * ListController constructor.
     *
     * @param EnterpriseAdminService $enterpriseAdminService
     */
    public function __construct(EnterpriseAdminService $enterpriseAdminService)
    {
        parent::__construct();

        //企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 企业管理员列表
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/20 15:26
     *
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->input();

        // 企业ID不能为空
        if (empty($params['enterprise_id']) || !is_numeric($params['enterprise_id'])) {
            return notice(201101, [], $params);
        }

        $limit = $params['limit'] ?? 10;

        $where = [
            'enterprise_id' => $params['enterprise_id'],
            'role'          => EnterpriseAdmin::SUPER_ADMIN
        ];

        // 默认只查超管，传多个角色类型只查超管，查多角色类型
        if (isset($params['role'])) {
            $where['role'] = (array)$params['role'];
        }

        if (isset($params['keyword'])) {
            $where['keyword'] = $params['keyword'];
        }

        // 状态检索
        if (isset($params['status'])) {
            $where['status'] = (int)$params['status'];
        }

        if (isset($params['admin_ids'])) {
            $where['id'] = (array)$params['admin_ids'];
        }

        // 列表数据查询
        $enterpriseAdmin = $this->enterpriseAdminService->paginate($limit, $where);

        return success([
            'data' => $enterpriseAdmin->items(),
            'meta' => [
                'total' => $enterpriseAdmin->total(),
                'size'  => $enterpriseAdmin->perPage(),
                'page'  => $enterpriseAdmin->currentPage()
            ]
        ]);
    }
}
