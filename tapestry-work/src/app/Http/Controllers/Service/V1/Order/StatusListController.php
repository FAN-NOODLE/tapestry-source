<?php
/**
 * Created by PhpStorm.
 * User: YuBin
 * Date: 2020年04月10日
 * Time: 14：42：33
 */

namespace App\Http\Controllers\Service\V1\Order;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CrsOrder;
use App\Models\Order;
use App\Models\User;
use App\Services\MemberService;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use App\Exceptions\Error;
use Illuminate\Http\Request;

class StatusListController extends BaseController
{
    const ORDER_REFUND_STATUS = [41, 61, 81, 82, 84];

    private $input = [];

    /**
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接收参数
        $params = $request->only([
            'store_id',
            'enterprise_id',
            'take_goods_code',
            'order_search',
            'worker_id',
            'page',
            'size',
            'status',
            'order_status',
            'order_type',
            'service_order_type',
            'delivery_type',
            'promoter_ids',
            'waiter_ids',
            'start_time',
            'end_time',
            'mem_uid'
        ]);

        // 验证规则
        $rules = [
            'enterprise_id' => 'integer',
            'promoter_ids'  => 'array',
            'start_time'    => 'integer',
            'end_time'      => 'integer',
        ];

        // 提示语
        $messages = [
            'promoter_ids.array'     => 'promoter_ids格式需为数组',
            'start_time.integer'     => 'start_time需为int类型',
            'end_time.integer'       => 'end_time需为int类型',
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);
        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }

        $params['size'] = $params['size'] ?? 10;
        $params['page'] = $params['page'] ?? 1;

        $this->input = $params;

        $params = ['enterprise_id' => $params['enterprise_id']];
        switchAppHost($params);

        $list = $this->getPromotionOrderList();

        return success($list);
    }

    /**
     * 推广订单
     *
     * @return array
     * @throws \Exception
     */
    private function getPromotionOrderList(): array
    {
        $params=$this->input;
        $list = CrsOrder::query()
            ->where("sync_source", CrsOrder::HAO_WU_JI_ORDER)
            ->whereIn("store_staff_id", $params['promoter_ids'])
            ->where("store_staff_id", '>',0)
            ->when(!empty($params['order_status']), function ($query)use ($params) {
                return $query->whereIn('status', $params['order_status']);
            })
            ->when(!empty($params['start_time']), function ($query) use ($params){
                return $query->where('pay_at',  '>=',$params['start_time']);
            })
            ->when(!empty($params['end_time']), function ($query) use ($params){
                return $query->where('pay_at',  '<=',$params['end_time']);
            })
            ->when(!empty($params['order_search']), function ($query) use ($params){
                return $query->where(function ($orderModel) use ($params){
                    return $orderModel
                        ->orWhere("order_no", 'like', '%' . $params['order_search'] . '%')
                        ->orWhere("delivery_mobile", 'like', '%' . $params['order_search'] . '%');
                });
            })
            ->orderBy("updated_at", 'desc')
            ->paginate($this->input['size'], ['*'], '', $this->input['page'])
            ->toArray();

        return $this->format($list);
    }

    /**
     * 格式化订单列表
     * @param array $orderData
     *
     * @return array|mixed
     * @throws \Exception
     */
    protected function format($orderData = [])
    {
        $orderList = $orderData['data'];
        $orderNos = array_column($orderList, 'order_no');
        if (empty($orderNos)) {
            $data = [
                'page'  => intval($this->input['page']),
                'size'  => intval($this->input['size']),
                'total' => 0,
                'list'  => []
            ];

            return $data;
        }

        $orderList = array_column($orderList, null, 'order_no');

        Log::info('gw-assistant订单列表请求参数' . json_encode(['order_nos' => $orderNos]));

        $orders = (new \CrsSdk\Shop\Order())->statusList(['order_nos' => $orderNos, 'order_type' => 1, 'status' => 'all']);

        Log::info('gw-assistant订单列表表返回参数' . json_encode($orders));
        if (!empty($orders['list'])) {
            $unionIds = array_unique(array_column($orderList, 'union_id'));

            $userLIst = User::query()
                ->whereIn('union_id', $unionIds)
                ->get()
                ->toArray();

            $userLIst = array_column($userLIst, null, 'union_id');

            foreach ($orders['list'] as $key => &$value) {
                // 计算实收金额 公式：实收金额 = 实付金额 + 其他优惠 - 退款总额
                $value['order_total_amount'] = $value['cash_fee'] + $value['preferential_price'] - $value['total_refund_fee'];
                //提货码隐藏
                if ($value['delivery_type'] == 2 && $value['take_goods_status'] == 1) {
                    $value['take_goods_code'] = '';
                }

                if (in_array($value['status'], self::ORDER_REFUND_STATUS)) {
                    $value['is_refund_order'] = 1;
                } else {
                    $value['is_refund_order'] = 0;
                }
                // 订单类型
                $value['service_order_type'] = $this->input['service_order_type'] ?? 0;
                $value['goods_total']        = $this->getGoodsNum($value['goods_list']);
                $value['parcel_button']      = false;
                // 需配送显示添加包裹按钮(类型为快递,状态已付款,有服务人)
                if ($value['delivery_type'] == 1 && $value['status'] == 60 && $value['service_staff_id'] > 0) {
                    $value['parcel_button'] = true;
                }

                $value['user_nick_name'] = '';
                $value['user_avatar']    = '';
                $orderInfo = $orderList[$value['order_no']] ?? [];
                $unionId   = $orderInfo['union_id'] ?? '';
                $userInfo = $userLIst[$unionId] ?? [];
                if (!empty($userInfo)) {
                    $value['user_nick_name'] = $userInfo['nick_name'];
                    $value['user_avatar']    = $userInfo['avatar'];
                }

                $value['ec_user_id'] = (new MemberService())->getUserEcUserId($this->input['mem_uid']??0, $unionId, $value['enterprise_id']);
            }
        }

        $orders = [
            'page'  => intval($this->input['page']),
            'size'  => intval($this->input['size']),
            'total' => $orderData['total'],
            'list'  => $orders['list']
        ];

        return $orders;
    }


    /**
     * 获取商品对数量
     *
     * @param array $goodsList
     *
     * @return int
     */
    protected function getGoodsNum($goodsList = [])
    {
        $num = 0;
        if (!empty($goodsList)) {
            $num = array_sum(array_column($goodsList, 'quantity'));
        }

        return $num;
    }

}