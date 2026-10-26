<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterpriseService;

class ListController extends BaseController
{
    /**
     * 企业service
     *
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * ListController constructor.
     *
     * @param EnterpriseService $enterpriseService
     */
    public function __construct(EnterpriseService $enterpriseService)
    {
        parent::__construct();

        // 企业service赋值
        $this->enterpriseService = $enterpriseService;
    }

    /**
     * 企业列表
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
        $limit = isset($params['limit']) ? $params['limit'] : 10;

        // 列表数据查询
        $where = empty($where) ? [] : $where;
        $enterprise = $this->enterpriseService->paginate($limit, $where);

        return success([
            'data' => $enterprise->items(),
            'meta' => [
                'total' => $enterprise->total(),
                'size'  => $enterprise->perPage(),
                'page'  => $enterprise->currentPage()
            ]
        ]);
    }
}
