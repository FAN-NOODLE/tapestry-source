<?php
/**
 * 订单详情
 * Created by PhpStorm.
 * User: 梁欢
 * Date: 2019/01/03
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Performance;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\UserServices;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;

class OrderDetailController extends BaseController
{
    /**
     * 用户services
     *
     * @var UserServices
     */
    protected $userServices;

    /**
     * OrderDetailController constructor.
     * @param UserServices $userServices
     */
    public function __construct(UserServices $userServices)
    {
        parent::__construct();

        // 用户services赋值
        $this->userServices = $userServices;
    }

    /**
     * 订单详情
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2019/01/03 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 详情
        $params = $request->only(['order_id']);
        $gets = (new Order())->orderDetail(['order_id' => $params['order_id']]);

        // 订单信息
        $orderData = [
            'id' => $gets['data']['id'],
            'order_no' => $gets['data']['order_no'],
            'pay_at' => $gets['data']['pay_at'],
            'created_at' => (int)$gets['data']['created_at'],
            'status' => (int)$gets['data']['status'],
            'payment_channel' => (int)$gets['data']['payment_channel'],
            'payment_no' => $gets['data']['pay_data']['third_payment_no'] ?? '',
            'final_amount' => (int)$gets['data']['final_amount'],
            'preferential_price' => (int)$gets['data']['preferential_price'],
            'total_amount' => (int)$gets['data']['total_amount'],
        ];

        // 用户列表
        $userData = $this->userServices->byList([$gets['data']['user_id']])[$gets['data']['user_id']] ?? [];

        return [
            'data' => [
                'order_data' => $orderData,
                // 用户信息
                'user_data' => [
                    'user_nick' => $userData['user_nick'],
                    'mobile' => $userData['mobile'],
                    'mem_face' => $userData['avatar'],
                ],
                // 商品列表
                'goods_data' => $gets['data']['goods_data'],
            ]
        ];
    }
}