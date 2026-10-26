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
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->input();

        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return error(201101);
        }

        $limit = isset($params['limit']) ? $params['limit'] : 10;
        $where['enterprise_id'] = $params['enterprise_id'];
        isset($params['keyword']) && $where['keyword'] = $params['keyword'];

        // 列表数据查询
        $where = empty($where) ? [] : $where;
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
