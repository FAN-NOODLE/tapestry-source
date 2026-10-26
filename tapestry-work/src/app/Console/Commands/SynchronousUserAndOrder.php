<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\Address;
use App\Models\Order;
use App\Models\OrderGoods;
use App\Models\OrderStatusTrace;
use App\Models\ShopOrder;
use App\Models\User;
use Illuminate\Console\Command;
use App\Models\ShopUser;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class SynchronousUserAndOrder extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'synchronous-user-and-order';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '同步旧数据';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }


    public function handle()
    {
        self::userAndAddress();
        self::orderAndGoodsAndTrace();
        echo 'synchronous:success';
    }

    // 同步客户和地址
    public function userAndAddress()
    {
        ShopUser::query()->with('addresses')
            ->chunk(500, function ($users) {
                $users = $users->toArray();
                $ids   = array_column($users, 'id');

                Log::info('shop-user-ids', $ids);

                if (!$ids) {
                    return;
                }

                // 已经同步过 id 集合(包括刪除的)
                $existIds = User::query()->withTrashed()->whereIn('id', $ids)->pluck('id')->toArray();

                Log::info('user-exist-ids', $existIds);

                $insertUsers   = [];
                $insertAddress = [];

                foreach ($users as $user) {
                    // 不存在就同步
                    if (!in_array($user['id'], $existIds)) {
                        $insertUsers[] = [
                            'id'              => $user['id'],
                            'enterprise_id'   => $user['enterprise_id'],
                            'union_id'        => $user['union_id'],
                            'open_id'         => $user['open_id'],
                            'nick_name'       => $user['nick_name'],
                            'real_name'       => $user['real_name'],
                            'real_name_index' => $user['real_name_index'],
                            'mobile'          => $user['mobile'],
                            'birthday'        => $user['birthday'],
                            'avatar'          => $user['avatar'],
                            'country'         => $user['country'],
                            'province'        => $user['province'],
                            'city'            => $user['city'],
                            'register_at'     => $user['register_at'],
                            'auth_at'         => $user['auth_at'],
                            'type'            => $user['auth_at'] ? 1 : 0,
                            'sex'             => $user['sex'],
                            'user_amount'     => $user['user_amount'],
                            'created_at'      => $user['created_at'],
                            'updated_at'      => $user['updated_at'],
                        ];

                        foreach ($user['addresses'] as $address) {
                            $insertAddress[] = [
                                'id'            => $address['id'],
                                'enterprise_id' => $address['enterprise_id'],
                                'user_id'       => $address['user_id'],
                                'province'      => $address['province'],
                                'city'          => $address['city'],
                                'area'          => $address['area'],
                                'address'       => $address['address'],
                                'province_code' => $address['province_code'],
                                'city_code'     => $address['city_code'],
                                'area_code'     => $address['area_code'],
                                'zip_code'      => $address['zip_code'],
                                'default'       => $address['default'],
                                'created_at'    => $address['created_at'],
                                'updated_at'    => $address['updated_at'],
                            ];
                        }
                    }
                }

                try {
                    DB::beginTransaction();
                    $insertUsers && User::query()->insert($insertUsers);
                    $insertAddress && Address::query()->insert($insertAddress);
                    DB::commit();
                } catch (\Exception $e) {
                    DB::rollBack();
                    Log::info('synchronous-user-error:', ['message' => $e->getMessage(), 'line' => $e->getLine()]);
                    exit;
                }
            });
    }

    // 同步订单和订单商品和订单变化
    public function orderAndGoodsAndTrace()
    {
        ShopOrder::query()
            ->with('goods')
            ->with('traces')
            ->chunk(500, function ($orders) {
                $orders = $orders->toArray();
                $ids    = array_column($orders, 'id');

                Log::info('shop-order-ids', $ids);

                if (!$ids) {
                    return;
                }

                // 已经同步过的 id 集合(包括删除的)
                $existIds = Order::query()->withTrashed()->whereIn('id', $ids)->pluck('id')->toArray();

                Log::info('order-exist-ids', $existIds);

                $insertOrder  = [];
                $insertGoods  = [];
                $insertTraces = [];

                foreach ($orders as $order) {
                    // 不存在就同步
                    if (!in_array($order['id'], $existIds)) {
                        $insertOrder[] = [
                            'id'                 => $order['id'],
                            'user_id'            => $order['user_id'],
                            'enterprise_id'      => $order['enterprise_id'],
                            'order_no'           => $order['order_no'],
                            'status'             => $order['status'],
                            'order_type'         => $order['order_type'],
                            'total_amount'       => $order['total_amount'],
                            'preferential_price' => $order['preferential_price'],
                            'final_amount'       => $order['final_amount'],
                            'payment_channel'    => $order['payment_channel'],
                            'pay_at'             => $order['pay_at'],
                            'refund_at'          => $order['refund_at'],
                            'order_at'           => $order['created_at'],
                            'delivery_name'      => $order['delivery_name'],
                            'delivery_mobile'    => $order['delivery_mobile'],
                            'delivery_address'   => $order['delivery_address'],
                            'remarks'            => $order['remarks'],
                            'province'           => $order['province'],
                            'city'               => $order['city'],
                            'area'               => $order['area'],
                            'logistic_identity'  => $order['logistic_identity'],
                            'logistic_no'        => $order['logistic_no'],
                            'send_at'            => $order['send_at'],
                            'sign_at'            => $order['sign_at'],
                            'cash_fee'           => $order['cash_fee'],
                            'coupon_fee'         => $order['coupon_fee'],
                            'coupon_count'       => $order['coupon_count'],
                            'other_fee'          => $order['other_fee'],
                            'created_at'         => $order['created_at'],
                            'updated_at'         => $order['updated_at'],
                        ];

                        // 商品
                        foreach ($order['goods'] as $good) {
                            $insertGoods[] = [
                                'id'              => $good['id'],
                                'enterprise_id'   => $good['enterprise_id'],
                                'order_no'        => $good['order_no'],
                                'sku_no'          => $good['sku_no'],
                                'spu_id'          => $good['spu_id'],
                                'sku_id'          => $good['sku_id'],
                                'spu_name'        => $good['spu_name'],
                                'type'            => $good['type'],
                                'goods_type'      => $good['goods_type'],
                                'category_name'   => $good['category_name'],
                                'sku_picture'     => $good['sku_picture'],
                                'sku_attribute'   => $good['sku_attribute'],
                                'quantity'        => $good['quantity'],
                                'original_price'  => $good['original_price'],
                                'sell_price'      => $good['sell_price'],
                                'single_discount' => 0,
                                'sub_price'       => $good['sub_price'],
                                'goods_pay_price' => 0,
                                'created_at'      => $good['created_at'],
                                'updated_at'      => $good['updated_at'],
                            ];
                        }

                        // 追踪
                        foreach ($order['traces'] as $trace) {
                            $insertTraces[] = [
                                'id'            => $trace['id'],
                                'enterprise_id' => $order['enterprise_id'],
                                'order_no'      => $order['order_no'],
                                'event'         => $trace['event'],
                                'event_desc'    => $trace['event_desc'],
                                'from'          => $trace['from_order_status'],
                                'to'            => $trace['to_order_status'],
                                'change_at'     => $trace['created_at'],
                                'created_at'    => $trace['created_at'],
                                'updated_at'    => $trace['updated_at'],
                            ];
                        }
                    }
                }

                // 事务处理
                try {
                    DB::beginTransaction();
                    $insertOrder && Order::query()->insert($insertOrder);
                    $insertGoods && OrderGoods::query()->insert($insertGoods);
                    $insertTraces && OrderStatusTrace::query()->insert($insertTraces);
                    DB::commit();
                    // 更新用户type
                    self::updateUserType($insertOrder);
                } catch (\Exception $e) {
                    DB::rollBack();
                    Log::info('synchronous-order-error:', ['message' => $e->getMessage(), 'line' => $e->getLine()]);
                    exit;
                }
            });
    }

    /**
     * 下单用户更新
     *
     * @param $params
     */
    public function updateUserType($params)
    {
        // 已成交订单id集合
        $hasDealUIds = [];
        // 未成交订单id集合
        $notDealUIds = [];

        foreach ($params as $order) {
            if (in_array($order['status'], ShopOrder::STATUS_PAY_ARRAY)) {
                $hasDealUIds[] = $order['user_id'];
            } else {
                $notDealUIds[] = $order['user_id'];
            }
        }

        $hasDealUIds && User::query()
            ->whereIn('id', $hasDealUIds)
            ->update([
                'type' => User::IS_ADDED
            ]);

        $notDealUIds && User::query()
            ->whereIn('id', $notDealUIds)
            ->where('type', '<>', User::IS_ADDED)
            ->update([
                'type' => User::IS_ADDITION
            ]);

        Log::info('has-deal-user-ids', $hasDealUIds);
        Log::info('not-deal-user-ids', $notDealUIds);
    }

}
