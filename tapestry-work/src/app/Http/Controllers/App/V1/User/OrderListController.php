<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/12/4
 * Time: 17:38
 */

namespace App\Http\Controllers\App\V1\User;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\Order;

class OrderListController extends BaseController
{

    /**
     * 根据 客户user_id 获取订单列表
     *
     * @return array|mixed|\Psr\Http\Message\StreamInterface
     */
    public function index()
    {
        $page = request()->post('page', 1);
        $size = request()->post('size', 10);
        $search = request()->post('search', []);

        if (!isset($search['user_id']) || empty($search['user_id'])) {
            return success([
                'data' => [],
                'meta' => [
                    'total' => 0,
                    'page' => $page,
                    'size' => $size
                ]
            ]);
        }

        return (new Order())->orderList([
            'page' => $page,
            'size' => $size,
            'search' => [
                'user_id' => $search['user_id'],
                'enterprise_id' => $this->enterpriseId
            ]
        ]);
    }
}
