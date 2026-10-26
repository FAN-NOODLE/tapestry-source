<?php
/**
 * 订单列表
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Performance;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\Order;
use \Illuminate\Http\Request;

class OrderListController extends BaseController
{
    /**
     * 是否有用户信息 1：是，2：否 2
     */
    const ORDER_USER_FALSE = 2;

    public function index(Request $request)
    {
        $input = $request->only(['search', 'page', 'size']);
        $input['search']['enterprise_id'] = $this->enterpriseId;
        $input['search']['is_user'] = self::ORDER_USER_FALSE;
        $input['search']['mem_uid'] = $this->memUid;
        $input['size'] = $input['size'] ?? Constant::DEFAULT_SIZE;

        // 接值
        $rules = [
            'search.status' => 'integer',
            'search.order_type' => 'integer',
            'page' => 'integer',
            'size' => 'integer',
        ];
        $messages = [
            'search.status.integer' => '订单状态必须为正整数',
            'search.order_type.integer' => '订单类型必须为正整数',
            'page.integer' => '页码必须为正整数',
            'size.integer' => '每页数量必须为正整数',
        ];
        $this->validate($request, $rules, $messages);

        // 订单列表
        $list = (new Order())->orderList($input);

        return $list;
    }
}
