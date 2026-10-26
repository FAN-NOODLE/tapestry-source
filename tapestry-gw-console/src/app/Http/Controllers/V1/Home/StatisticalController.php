<?php
/**
 * 首页客户概况数据统计
 * User: zhuhongda
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Home;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use Carbon\Carbon;
use CrsSdk\Shop\NewOrder;
use CrsSdk\Shop\Order;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class StatisticalController extends BaseController
{
    /**
     *  服务类型： 商城
     */
    const SERVE_TYPE_SHOP = 1;
    /**
     *  服务类型： crs
     */
    const SERVE_TYPE_CRS = 2;

    /**
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only(['date_type']);
        // 企业id
        $input['enterprise_id'] = $this->enterpriseId;

        // 默认今日
        if (!isset($input['date_type'])) {
            $input['date_type'] = Constant::DAY;
        }

        // 非超管按照门店筛选
        $roleInfo = $this->getStores($this->enterpriseAdmin['role']);
        if (!$roleInfo['is_super_admin']) {
            $input['stores'] = $roleInfo['stores'];
            $input['type'] = $roleInfo['type'];
        }

        // 查询下单客户数、会员客户数、好友客户数、总客户数
        $statistical = (new Order())->salesStatistical($input);

        $result = $this->getDataByServeType($statistical)['data'];

        return success($result);
    }

    /**
     * 首页数据统计
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function statistics(Request $request)
    {
        // 首页默认数据
        $data = [
            // 成交订单数
            "complete_order"  => 0,
            // 成交金额
            "complete_amount" => 0,
            // 退款金额(单位分)
            "refund_amount"   => 0,
            // 客单价(单位分)
            "price"           => 0,
            // 访客数
            "visit"           => 0,
            // 下单买家数
            "order_person"    => 0,
            // 支付买家数
            "pay_person"      => 0,
            // 支付转化率
            "pay_rate"        => '0.00%',
            // 下单转换率
            "order_rate"      => '0.00%',
            // 下单支付转化率
            "pays_rate"       => '0.00%',
        ];

        // 先查询订单类统计数据
        $statisticsParams = [
            'enterprise_id'      => $this->enterpriseId,
            'time_type'          => 0,
            'start_time'         => $request->get('start_time', Carbon::today()->timestamp),
            'end_time'           => $request->get('end_time', Carbon::now()->timestamp),
        ];
        $storeId = $request->get('store_id');
        if (!empty($storeId)) {
            $statisticsParams['store_id'] = $storeId;
        }

        try {
            $orderStatisticsData = app(NewOrder::class)->orderStatistics($statisticsParams);
            // 成交订单数
            if (isset($orderStatisticsData['payed_order_count'])) {
                $data['complete_order'] = intval($orderStatisticsData['payed_order_count']);
            }

            // 成交金额
            if (isset($orderStatisticsData['payed_order_amount'])) {
                $data['complete_amount'] = intval($orderStatisticsData['payed_order_amount']);
            }

            // 退款金额
            if (isset($orderStatisticsData['total_order_refund_amount'])) {
                $data['refund_amount'] = intval($orderStatisticsData['total_order_refund_amount']);
            }

            // 客单价
            if (isset($orderStatisticsData['per_user_amount'])) {
                $data['price'] = intval($orderStatisticsData['per_user_amount']);
            }

            // 下单买家数
            if (isset($orderStatisticsData['place_order_user_count'])) {
                $data['order_person'] = intval($orderStatisticsData['place_order_user_count']);
            }

            // 支付买家数
            if (isset($orderStatisticsData['payed_user_count'])) {
                $data['pay_person'] = intval($orderStatisticsData['payed_user_count']);
            }
        } catch (\Exception $exception) {
            Log::error("[console]orderStatistics 调用出错了，错误信息:{$exception->getMessage()}", $statisticsParams);
        }
        // 查询访客数
        try {
            $userVisit = app(User::class)->accessRecords($statisticsParams);
            $data['visit'] = isset($userVisit['uv']) ? intval($userVisit['uv']) : 0;
        } catch (\Exception $exception) {
            Log::error("[console]accessRecords 调用出错了，错误信息:{$exception->getMessage()}", $statisticsParams);
        }

        if ($data['visit'] > 0 && $data['pay_person'] > 0) {
            // 支付转化率 = 支付买家数 / 访客数 * 100
            $data['pay_rate'] = number_format($data['pay_person'] * 100 / $data['visit'], 2) . '%';
        }

        if ($data['order_person'] > 0 && $data['visit'] > 0) {
            // 下单转换率 = 下单人数 / 访客数 * 100
            $data['order_rate'] = number_format($data['order_person'] * 100 / $data['visit'], 2). '%';
        }

        if ($data['order_person'] > 0 && $data['pay_person'] > 0) {
            // 下单支付转化率 = 支付买家数 / 下单买家数 * 100
            $data['pays_rate'] = number_format($data['pay_person']  * 100/ $data['order_person'], 2) . '%';
        }

        return success($data);
    }

    public function getDataByServeType($statistical)
    {
        // 默认全部开通 shop+crs
        $result = [
            'data' => [
                'user_all'       => $statistical['data']['user_all'],
                'user_order'     => $statistical['data']['user_order'],
                'member_number'  => $statistical['data']['user_number'],
                'friends_number' => $statistical['data']['customer_number']

            ]
        ];

        // 仅开通shop
        if ($this->serveType == self::SERVE_TYPE_SHOP) {
            $result = [
                'data' => [
                    'user_all'       => $statistical['data']['user_order'] + $statistical['data']['user_number'],
                    'user_order'     => $statistical['data']['user_order'],
                    'member_number'  => $statistical['data']['user_number']
                ]
            ];
        };
        // 仅开通Crs
        if ($this->serveType == self::SERVE_TYPE_CRS) {
            $result = [
                'data' => [
                    'user_all'       => $statistical['data']['customer_number'],
                    'friends_number' => $statistical['data']['customer_number']
                ]
            ];
        };
        return $result;
    }
}
