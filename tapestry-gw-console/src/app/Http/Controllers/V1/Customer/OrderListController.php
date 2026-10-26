<?php

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use App\Common\Constant;
use CrsSdk\Shop\Order;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Exception;

class OrderListController extends BaseController
{
    /**
     * 客户订单列表
     *
     * @param Request $request
     *
     * @return JsonResponse|void
     * @throws Exception
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', Constant::DEFAULT_SIZE);
        // 查询条件
        $search = $request->get('search', []);

        // 客户id不能为空
        if (!isset($search['user_id']) || (isset($search['user_id']) && !$search['user_id'])) {
            error(320051);
        }

        // 企业id
        $search['enterprise_id'] = $this->enterpriseId;

        // 调用sdk
        $params = [
            'page'   => $page,
            'size'   => $size,
            'search' => $search
        ];
        $list = (new Order())->orderList($params);

        return success($list);
    }
}
