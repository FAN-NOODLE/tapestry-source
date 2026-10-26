<?php

namespace App\Jobs;

use App\Common\Constant;
use App\Models\CrsOrder;
use App\Models\CrsOrderGoods;
use App\Models\CrsOrderStatusTrace;
use App\Models\CrsWorkerThirdPartRelation;
use App\Models\GuideUserRelation;
use App\Models\MemberUserRelation;
use App\Models\Order;
use App\Models\OrderGoods;
use App\Models\OrderStatusTrace;
use App\Models\User;
use App\Models\UserAddress;
use App\Models\UserCoupon;
use App\Models\WxCoupon;
use App\Services\PartnerMemberService;
use Illuminate\Queue\Jobs\Job;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Overtrue\Pinyin\Pinyin;
use App\Models\PartnerEnterpriseMap;

class CMQPlainJobHandler
{
    /**
     * Execute the job.
     *
     * @param Job $job
     * @param string $payload
     *
     * @return void
     */
    public function handle(Job $job, $payload)
    {
        $start = microtime(true);
        $userQueueName   = config('cmq_setting.user_queue_name');
        $orderQueueName  = config('cmq_setting.order_queue_name');
        $couponQueueName = config('cmq_setting.coupon_queue_name');
        $tagCustomerGroup = config('cmq_setting.tag_customer_group');
        $syncCrsOrder = config('cmq_setting.sync_crs_order');
        $syncHwjOrder = config('cmq_setting.sync_hwj_order');
        $syncHwjUser = config('cmq_setting.sync_hwj_user');

        $queue = $job->getQueue();

        try {
            Log::info('消费队列消息: ' . $payload);
            Log::info('订单队列信息: ' . $syncHwjOrder);
            Log::info('用户队列信息: ' . $syncHwjUser);

            switch ($queue) {
                case $userQueueName:
                    $this->userHandle($payload);
                    break;
                case $orderQueueName:
                    $this->orderHandle($payload);
                    break;
                case $couponQueueName:
                    $this->couponHandle($payload);
                    break;
                case $syncCrsOrder:
                    $this->syncCrsOrder($payload);
                    break;
                case $syncHwjOrder:
                    $this->syncHwjOrder($payload);
                    break;
                case $syncHwjUser:
                    $this->syncHwjUer($payload);
                    break;
                default:
                    // do nothing
            }

            // 删除消息
            if (!$job->isDeletedOrReleased()) {
                Log::info('delete the message:' . $job->getQueue());
                $job->delete();
            }
        } catch (\Exception $exception) {
            Log::info('消费队列消息异常: ' . json_encode($exception->getMessage()));
        }
        $costTime = round((microtime(true) - $start) * 1000, 2);
        Log::info('消费队列消息处理完成，耗时：' . $costTime . 'ms');
    }

    /**
     * 处理用户队列数据
     *
     * @param  $payload
     */
    public function userHandle($payload)
    {
        /**
         * 1、从CMQ消费用户消息(一次只有一条客户信息)；
         * 2、根据客户ID查询客户信息；
         * 3、查询订单信息，判断客户状态；
         * 4、判断客户信息，如存在，则到5，如不存在，则到6；
         * 5、更新客户信息，执行完到7；
         * 6、将客户信息、客户地址信息写入表中，执行完到7；
         * 7、查询客户地址信息，如存在，则到8，如不存在，则到9；
         * 8、更新地址信息；
         * 9、删除原有客户默认地址信息，添加新客户默认地址信息，执行完到10；
         * 10、查询客户分配导购信息，如果user_id为空(即为外部联系人绑定导购)，更新user_id；
         */
        $pinyinTool = new Pinyin();

        // 1、从CMQ消费用户消息(一次只有一条客户信息)；
        // 客户信息
        $userInfo = json_decode($payload, true)['userInfo'];
        // 客户地址信息
        $defaultAddress = $userInfo['defaultAddress'] ?? [];

        if (!isset($userInfo['id'])) {
            return;
        }
        if (empty($userInfo['unionId'])) {
            return;
        }

        // 2、根据客户ID查询客户信息
        $user = User::query()->where(['open_id' => $userInfo['openId']])->first(User::$selectField);

        $userInfoData = [
            'enterprise_id'   => $userInfo['enterpriseId'],
            'union_id'        => $userInfo['unionId'],
            'open_id'         => $userInfo['openId'],
            'nick_name'       => $userInfo['nickName'],
            'real_name'       => $userInfo['realName'],
            'real_name_index' => implode('', $pinyinTool->name($userInfo['realName'])),
            'mobile'          => $userInfo['mobile'],
            'birthday'        => $userInfo['birthday'],
            'avatar'          => $userInfo['avatar'],
            'country'         => $userInfo['country'],
            'province'        => $userInfo['province'],
            'city'            => $userInfo['city'],
            'register_at'     => $userInfo['registerAt'],
            'auth_at'         => $userInfo['authAt'],
            'sex'             => $userInfo['sex'],
            'user_amount'     => $userInfo['userAmount'],
        ];

        // 4、判断客户信息，如存在，则到5，如不存在，则到6
        if ($user) {
            // 5、更新客户信息，执行完到7
            User::query()->where(['id' => $user['id']])->update($userInfoData);
        } else {
            // 6、将客户信息写入表中，执行完到7
            $user = User::query()->create($userInfoData);
        }

        $this->updateOrderCouponUserId($userInfo['openId'], $user['id'], $userInfo['unionId']);

        // 查询订单信息，判断客户状态
        if (empty($userInfo['authAt'])) {
            $userInfoType = Constant::UN_AUTH; // 未授权
        } else {
            $order = Order::query()->where('user_id', $user['id'])->first(Order::$selectField);
            if ($order) {
                $paidOrder = Order::query()->where('user_id', $user['id'])
                    ->whereIn('status', [25, 30, 40, 41, 60, 61, 80, 81, 82, 83, 84, 100, 200, 201])->first(Order::$selectField);
                if ($paidOrder) {
                    $userInfoType = Constant::TRADED; // 已成交
                } else {
                    $userInfoType = Constant::UN_TRADE; // 未成交
                }
            } else {
                $userInfoType = Constant::AUTH; // 已授权
            }
        }

        User::query()->where('id', $user['id'])->update(['type' => $userInfoType]);

        $defaultAddressData = [];

        if ($defaultAddress) {
            $defaultAddressData = [
                'enterprise_id' => $userInfo['enterpriseId'],
                'user_id'       => $user['id'],
                'province'      => $defaultAddress['province'],
                'city'          => $defaultAddress['city'],
                'area'          => $defaultAddress['area'],
                'address'       => $defaultAddress['address'],
                'province_code' => $defaultAddress['provinceCode'] ?? '',
                'city_code'     => $defaultAddress['cityCode'] ?? '',
                'area_code'     => $defaultAddress['areaCode'] ?? '',
                'zip_code'      => $defaultAddress['zipCode'] ?? '',
                'default'       => 1, // 默认地址
                'deleted_at'    => 0,
            ];
        }

        if ($defaultAddressData) {
            // 7、查询客户地址信息，如存在，则到8，如不存在，则到9
            $userAddress = UserAddress::query()->withTrashed()->where([
                'address' => $defaultAddress['address'],
                'user_id' => $user['id'],
            ])->first(UserAddress::$selectField);

            if ($userAddress) {
                // 8、更新地址信息
                UserAddress::query()->where(['user_id' => $user['id']])->delete();

                UserAddress::query()->withTrashed()->where(['id' => $userAddress['id']])->update($defaultAddressData);
            } else {
                // 9、删除原有客户默认地址信息，添加新客户默认地址信息，执行完到10
                UserAddress::query()->where(['user_id' => $user['id']])->delete();

                UserAddress::query()->create($defaultAddressData);
            }
        }

        // 10、查询客户分配导购信息，如果user_id为空(即为外部联系人绑定导购)，更新user_id
        $guideUserRelation = GuideUserRelation::query()->where('union_id', $userInfo['unionId'])->first(GuideUserRelation::$selectField);

        if ($guideUserRelation && empty($guideUserRelation['user_id'])) {
            GuideUserRelation::query()->where('union_id', $userInfo['unionId'])->update(['user_id' => $user['id']]);
        }

        if ($userInfo['enterpriseId'] == 86) {
            GuideUserRelation::query()->where('user_id', $user['id'])->update(['union_id' => $userInfo['unionId']]);
        }

        // 创建|更新用户同时更新外联系人表中的user_id
        MemberUserRelation::query()->where('union_id', $userInfo['unionId'])->update(['user_id' => $user['id']]);
    }

    /**
     * 处理用户队列数据-new
     *
     * @param  $payload
     */
    public function syncHwjUer($payload)
    {
        /**
         * 1、从CMQ消费用户消息(一次只有一条客户信息)；
         * 2、根据客户ID查询客户信息；
         * 3、查询订单信息，判断客户状态；
         * 4、判断客户信息，如存在，则到5，如不存在，则到6；
         * 5、更新客户信息，执行完到7；
         * 6、将客户信息、客户地址信息写入表中，执行完到7；
         * 7、查询客户地址信息，如存在，则到8，如不存在，则到9；
         * 8、更新地址信息；
         * 9、删除原有客户默认地址信息，添加新客户默认地址信息，执行完到10；
         * 10、查询客户分配导购信息，如果user_id为空(即为外部联系人绑定导购)，更新user_id；
         */
        $pinyinTool = new Pinyin();

        // 1、从CMQ消费用户消息(一次只有一条客户信息)；
        // 客户信息
        $userInfo = json_decode($payload, true)['userInfo'];

        Log::info('用户队列信息json: ' . json_encode($userInfo));
        // 客户地址信息
        $defaultAddress = $userInfo['defaultAddress'] ?? [];

        if (!isset($userInfo['id'])) {
            return;
        }
        if (empty($userInfo['unionId'])) {
            return;
        }

        // 2、根据客户ID查询客户信息
        $user = User::query()->where(['open_id' => $userInfo['openId']])->first(User::$selectField);

        $userInfoData = [
            'enterprise_id'   => $userInfo['enterpriseId'],
            'union_id'        => $userInfo['unionId'],
            'open_id'         => $userInfo['openId'],
            'nick_name'       => $userInfo['nickName'],
            'real_name'       => $userInfo['realName'],
            'real_name_index' => implode('', $pinyinTool->name($userInfo['realName'])),
            'mobile'          => $userInfo['mobile'],
            'birthday'        => $userInfo['birthday'],
            'avatar'          => $userInfo['avatar'],
            'country'         => $userInfo['country'],
            'province'        => $userInfo['province'],
            'city'            => $userInfo['city'],
            'register_at'     => $userInfo['registerAt'],
            'auth_at'         => $userInfo['authAt'],
            'sex'             => $userInfo['sex'],
            'user_amount'     => $userInfo['userAmount'],
            'sync_source'     => CrsWorkerThirdPartRelation::SYNC_SOURCE_HWJ
        ];

        // 4、判断客户信息，如存在，则到5，如不存在，则到6
        if ($user) {
            // 5、更新客户信息，执行完到7
            User::query()->where(['id' => $user['id']])->update($userInfoData);
        } else {
            // 6、将客户信息写入表中，执行完到7
            $user = User::query()->create($userInfoData);
        }

        $this->updateOrderCouponUserId($userInfo['openId'], $user['id'], $userInfo['unionId']);

        // 查询订单信息，判断客户状态
//        if (empty($userInfo['authAt'])) {
//            $userInfoType = Constant::UN_AUTH; // 未授权
//        } else {
//            $order = Order::query()->where('user_id', $user['id'])->first(Order::$selectField);
//            if ($order) {
//                $paidOrder = Order::query()->where('user_id', $user['id'])
//                    ->whereIn('status', [25, 30, 40, 41, 60, 61, 80, 81, 82, 83, 84, 100, 200, 201])->first(Order::$selectField);
//                if ($paidOrder) {
//                    $userInfoType = Constant::TRADED; // 已成交
//                } else {
//                    $userInfoType = Constant::UN_TRADE; // 未成交
//                }
//            } else {
//                $userInfoType = Constant::AUTH; // 已授权
//            }
//        }

//        User::query()->where('id', $user['id'])->update(['type' => $userInfoType]);

        $defaultAddressData = [];

        if ($defaultAddress) {
            $defaultAddressData = [
                'enterprise_id' => $userInfo['enterpriseId'],
                'user_id'       => $user['id'],
                'province'      => $defaultAddress['province'],
                'city'          => $defaultAddress['city'],
                'area'          => $defaultAddress['area'],
                'address'       => $defaultAddress['address'],
                'province_code' => $defaultAddress['provinceCode'] ?? '',
                'city_code'     => $defaultAddress['cityCode'] ?? '',
                'area_code'     => $defaultAddress['areaCode'] ?? '',
                'zip_code'      => $defaultAddress['zipCode'] ?? '',
                'default'       => 1, // 默认地址
                'deleted_at'    => 0,
            ];
        }

        if ($defaultAddressData) {
            // 7、查询客户地址信息，如存在，则到8，如不存在，则到9
            $userAddress = UserAddress::query()->withTrashed()->where([
                'address' => $defaultAddress['address'],
                'user_id' => $user['id'],
            ])->first(UserAddress::$selectField);

            if ($userAddress) {
                // 8、更新地址信息
                UserAddress::query()->where(['user_id' => $user['id']])->delete();

                UserAddress::query()->withTrashed()->where(['id' => $userAddress['id']])->update($defaultAddressData);
            } else {
                // 9、删除原有客户默认地址信息，添加新客户默认地址信息，执行完到10
                UserAddress::query()->where(['user_id' => $user['id']])->delete();

                UserAddress::query()->create($defaultAddressData);
            }
        }

        // 10、查询客户分配导购信息，如果user_id为空(即为外部联系人绑定导购)，更新user_id
        $guideUserRelation = GuideUserRelation::query()->where('union_id', $userInfo['unionId'])->first(GuideUserRelation::$selectField);

        if ($guideUserRelation && empty($guideUserRelation['user_id'])) {
            GuideUserRelation::query()->where('union_id', $userInfo['unionId'])->update(['user_id' => $user['id']]);
        }

        if ($userInfo['enterpriseId'] == 86) {
            GuideUserRelation::query()->where('user_id', $user['id'])->update(['union_id' => $userInfo['unionId']]);
        }

        // 创建|更新用户同时更新外联系人表中的user_id
        MemberUserRelation::query()->where('union_id', $userInfo['unionId'])->update(['user_id' => $user['id']]);
    }

    /**
     * 处理订单队列数据
     *
     * @param  $payload
     */
    public function orderHandle($payload)
    {
        /**
         * 1、从CMQ消费订单消息(一次只有一条订单信息)；
         * 2、查询客户信息，如不存在，新增客户信息，判断客户状态；
         * 3、根据订单ID查询订单信息，如存在，则到4，如不存在，则到5；
         * 4、更新订单信息，执行完到6；
         * 5、新增订单信息，将订单信息写入表中，执行完到6；
         * 6、查询客户分配导购信息，如果user_id为空(即为外部联系人绑定导购)，更新user_id。
         */
        $pinyinTool = new Pinyin();

        // 1、从CMQ消费订单消息(一次只有一条订单信息)；
        $msgBody = json_decode($payload, true);
        // 订单信息
        $orderInfo = $msgBody['orderInfo'];
        // 订单商品
        $orderGoods = $msgBody['orderGoods'];
        // 订单状态变化
        $orderStatusTrace = $msgBody['orderStatusTrace'];
        // 订单用户
        $orderUser = $msgBody['orderUser'];

        if (!isset($orderUser['id'])) {
            return;
        }
        if (empty($orderUser['openId'])) {
            return;
        }

        $user = [];

        if (!empty($orderUser['unionId'])) {
            // 2、查询客户信息，判断客户状态，如不存在，新增客户信息，
            $user = User::query()->where('open_id', $orderUser['openId'])->first(User::$selectField);

            try {
                $userInfoData = [
                    'enterprise_id'   => $orderUser['enterpriseId'],
                    'union_id'        => $orderUser['unionId'],
                    'open_id'         => $orderUser['openId'],
                    'nick_name'       => $orderUser['nickName'],
                    'real_name'       => $orderUser['realName'],
                    'real_name_index' => implode('', $pinyinTool->name($orderUser['realName'])),
                    'mobile'          => $orderUser['mobile'],
                    'birthday'        => $orderUser['birthday'],
                    'avatar'          => $orderUser['avatar'],
                    'country'         => $orderUser['country'],
                    'province'        => $orderUser['province'],
                    'city'            => $orderUser['city'],
                    'register_at'     => $orderUser['registerAt'],
                    'auth_at'         => $orderUser['authAt'],
                    'sex'             => $orderUser['sex'],
                    'user_amount'     => $orderUser['userAmount'],
                ];

                // 如不存在，新增客户信息
                if ($user) {
                    User::query()->where('id', $user['id'])->update($userInfoData);
                } else {
                    $user = User::query()->create($userInfoData);
                }

                $defaultAddressData = [];
                if ($orderUser['defaultAddress'] ?? []) {
                    $defaultAddressData = [
                        'enterprise_id' => $orderUser['enterpriseId'],
                        'user_id'       => $user['id'],
                        'province'      => $orderUser['defaultAddress']['province'],
                        'city'          => $orderUser['defaultAddress']['city'],
                        'area'          => $orderUser['defaultAddress']['area'],
                        'address'       => $orderUser['defaultAddress']['address'],
                        'province_code' => $orderUser['defaultAddress']['provinceCode'] ?? '',
                        'city_code'     => $orderUser['defaultAddress']['cityCode'] ?? '',
                        'area_code'     => $orderUser['defaultAddress']['areaCode'] ?? '',
                        'zip_code'      => $orderUser['defaultAddress']['zipCode'] ?? '',
                        'default'       => 1, // 默认地址
                        'deleted_at'    => 0,
                    ];
                }

                if ($defaultAddressData) {
                    // 根据客户地址ID查询客户地址信息
                    $userAddress = UserAddress::query()->withTrashed()->where([
                        'address' => $orderUser['defaultAddress']['address'],
                        'user_id' => $user['id'],
                    ])->first(UserAddress::$selectField);

                    if ($userAddress) {
                        // 更新地址信息
                        UserAddress::query()->where(['user_id' => $user['id']])->delete();

                        UserAddress::query()->withTrashed()->where(['id' => $userAddress['id']])->update($defaultAddressData);
                    } else {
                        // 删除原有客户默认地址信息，添加新客户默认地址信息
                        UserAddress::query()->where(['user_id' => $user['id']])->delete();

                        UserAddress::query()->create($defaultAddressData);
                    }
                }
            } catch (\Exception $exception) {
                Log::info('消费订单队列消息用户数据异常: ' . json_encode($exception->getMessage()));
            }
        }

        // 3、根据订单ID查询订单信息，如存在，则到4，如不存在，则到5
        $order = Order::query()->where('order_no', $orderInfo['orderNo'])->first(Order::$selectField);

        $orderInfoData = [
            'user_id'            => $user['id'] ?? 0,
            'union_id'           => $orderUser['unionId'],
            'open_id'            => $orderUser['openId'],
            'enterprise_id'      => $orderInfo['enterpriseId'],
            'order_no'           => $orderInfo['orderNo'],
            'status'             => $orderInfo['status'],
            'order_type'         => $orderInfo['type'],
            'total_amount'       => $orderInfo['totalAmount'],
            'preferential_price' => $orderInfo['preferentialPrice'],
            'final_amount'       => $orderInfo['finalAmount'],
            'payment_channel'    => $orderInfo['paymentChannel'],
            'pay_at'             => $orderInfo['payAt'],
            'refund_at'          => $orderInfo['refundAt'],
            'order_at'           => $orderInfo['createdAt'],
            'delivery_name'      => $orderInfo['deliveryName'],
            'delivery_mobile'    => $orderInfo['deliveryMobile'],
            'delivery_address'   => $orderInfo['deliveryAddress'],
            'remarks'            => $orderInfo['remarks'],
            'province'           => $orderInfo['province'],
            'city'               => $orderInfo['city'],
            'area'               => $orderInfo['area'],
            'logistic_identity'  => $orderInfo['logisticIdentity'],
            'logistic_no'        => $orderInfo['logisticNo'],
            'send_at'            => $orderInfo['sendAt'],
            'sign_at'            => $orderInfo['signAt'],
            'cash_fee'           => $orderInfo['cashFee'],
            'coupon_fee'         => $orderInfo['couponFee'],
            'coupon_count'       => $orderInfo['couponCount'],
            'other_fee'          => $orderInfo['otherFee'],
        ];

        $orderGoodsData = [];
        foreach ($orderGoods as $value) {
            $orderGoodsData[] = [
                'enterprise_id'   => $orderInfo['enterpriseId'],
                'order_no'        => $orderInfo['orderNo'],
                'sku_no'          => $value['skuNo'],
                'spu_id'          => $value['spuId'],
                'sku_id'          => $value['skuId'],
                'spu_name'        => $value['spuName'],
                'type'            => $value['type'],
                'goods_type'      => $value['goodsType'],
                'category_name'   => $value['categoryName'],
                'sku_picture'     => $value['skuPicture'],
                'sku_attribute'   => $value['skuAttribute'],
                'quantity'        => $value['quantity'],
                'original_price'  => $value['originalPrice'],
                'sell_price'      => $value['sellPrice'],
                'single_discount' => $value['singleDiscount'],
                'sub_price'       => $value['subPrice'],
                'goods_pay_price' => $value['goodsPayPrice'] ?? 0,
                'created_at'      => time(),
                'updated_at'      => time(),
            ];
        }

        $orderStatusTraceData = [];
        if ($orderStatusTrace) {
            $orderStatusTraceData = [
                'enterprise_id' => $orderInfo['enterpriseId'],
                'order_no'      => $orderInfo['orderNo'],
                'event'         => $orderStatusTrace['event'] ?? '',
                'event_desc'    => $orderStatusTrace['eventDesc'] ?? '',
                'from'          => $orderStatusTrace['from'] ?? 0,
                'to'            => $orderStatusTrace['to'] ?? 0,
                'change_at'     => $orderStatusTrace['createdAt'] ?? 0,
            ];
        }

        if ($order) {
            // 4、更新订单信息，执行完到6
            Order::query()->where('order_no', $orderInfo['orderNo'])->update($orderInfoData);
        } else {
            // 5、新增订单信息，将订单信息写入表中，执行完到6
            Order::query()->create($orderInfoData);

            if ($orderGoodsData) {
                OrderGoods::query()->insert($orderGoodsData);
            }
        }

        if ($orderStatusTraceData) {
            OrderStatusTrace::query()->create($orderStatusTraceData);
        }

        if ($user) {
            // 6、查询客户分配导购信息，如果user_id为空(即为外部联系人绑定导购)，更新user_id
            $guideUserRelation = GuideUserRelation::query()->where('union_id', $orderUser['unionId'])->first(GuideUserRelation::$selectField);

            if ($guideUserRelation && empty($guideUserRelation['user_id'])) {
                GuideUserRelation::query()->where('union_id', $orderUser['unionId'])->update(['user_id' => $user['id']]);
            }

            if ($orderUser['enterpriseId'] == 86) {
                GuideUserRelation::query()->where('user_id', $user['id'])->update(['union_id' => $orderUser['unionId']]);
            }

            // 创建|更新用户同时更新外联系人表中的user_id
            MemberUserRelation::query()->where('union_id', $orderUser['unionId'])->update(['user_id' => $user['id']]);

            $this->updateOrderCouponUserId($orderUser['openId'], $user['id'], $orderUser['unionId']);

            // 判断客户状态
            if (in_array($user['type'], [Constant::UN_AUTH, Constant::AUTH, Constant::UN_TRADE])) {
                // 30 待确认, 40 待发货, 41 售后中(待发货)，60 待收货，61 售后中(待收货)，80 已签收，81 售后中(退货退款)，82 售后中(仅退款)，83 部分退款，84 售后中(部分退款)，100 已完成，200退款中，201 退款失败
                if (in_array($orderInfo['status'], [25, 30, 40, 41, 60, 61, 80, 81, 82, 83, 84, 100, 200, 201])) {
                    $userType = Constant::TRADED;
                } else {
                    $userType = Constant::UN_TRADE;
                }

                User::query()->where('id', $user['id'])->update(['type' => $userType]);
            }
        }
    }

    /**
     * 处理优惠券队列数据
     *
     * @param $payload
     */
    public function couponHandle($payload)
    {
        /**
         * 1、从CMQ消费优惠券消息(一次只有一条优惠券信息)；
         * 2、根据优惠券ID查询优惠券信息，如存在，则到3，如不存在，则到4；
         * 3、更新优惠券信息，执行完到5；
         * 4、新增优惠券信息，将优惠券信息写入表中，执行完到5；
         * 5、查询微信优惠券信息，如不存在，新增微信优惠券信息；
         */

        // 1、从CMQ消费优惠券消息(一次只有一条优惠券信息)；
        $userCouponCMQInfo = json_decode($payload, true)['userCouponInfo'];

        if (!isset($userCouponCMQInfo['id'])) {
            return;
        }
        if (!isset($userCouponCMQInfo['wx_coupon'])) {
            return;
        }
        if (empty($userCouponCMQInfo['open_id'])) {
            return;
        }

        // 微信优惠券信息
        $wxCouponCMQInfo = $userCouponCMQInfo['wx_coupon'];

        // 2、根据优惠券ID查询优惠券信息，如存在，则到3，如不存在，则到4；
        $userCouponInfo = UserCoupon::query()->where('shop_user_coupon_id', $userCouponCMQInfo['id'])->first(UserCoupon::$selectField);

        $userId = 0;

        if (!empty($userCouponCMQInfo['union_id'] ?? '')) {
            $user = User::query()->where('union_id', $userCouponCMQInfo['union_id'])->first(User::$selectField);
            if ($user) {
                $userId = $user->id;
            }
        }

        $userCouponData = [
            'shop_user_coupon_id' => $userCouponCMQInfo['id'],
            'enterprise_id'       => $userCouponCMQInfo['enterprise_id'],
            'user_id'             => $userId,
            'union_id'            => $userCouponCMQInfo['union_id'] ?? '',
            'open_id'             => $userCouponCMQInfo['open_id'] ?? '',
            'effective'           => $userCouponCMQInfo['effective'],
            'start_time'          => $userCouponCMQInfo['start_time'],
            'end_time'            => $userCouponCMQInfo['end_time'],
            'status'              => $userCouponCMQInfo['status'],
            'wx_coupon_id'        => $userCouponCMQInfo['wx_coupon_id'],
            'receive_time'        => $userCouponCMQInfo['receive_time'],
            'use_time'            => $userCouponCMQInfo['use_time'] ?? 0,
            'order_no'            => $userCouponCMQInfo['order_no'] ?? '',
            'order_amount'        => $userCouponCMQInfo['order_amount'] ?? 0,
        ];

        $wxCouponData = [
            'shop_wx_coupon_id' => $wxCouponCMQInfo['id'],
            'enterprise_id'     => $wxCouponCMQInfo['enterprise_id'],
            'coupon_name'       => $wxCouponCMQInfo['coupon_name'],
            'coupon_value'      => $wxCouponCMQInfo['coupon_value'],
            'coupon_minimum'    => $wxCouponCMQInfo['coupon_mininumn'],
        ];

        if (!$userCouponInfo) {
            // 3、更新优惠券信息，执行完到5；
            UserCoupon::query()->create($userCouponData);
        } else {
            // 4、新增优惠券信息，将优惠券信息写入表中，执行完到5；
            UserCoupon::query()->where('shop_user_coupon_id', $userCouponCMQInfo['id'])->update($userCouponData);
        }

        // 5、查询微信优惠券信息，如不存在，新增微信优惠券信息；
        $wxCouponInfo = WxCoupon::query()->where('shop_wx_coupon_id', $wxCouponCMQInfo['id'])->first(WxCoupon::$selectField);

        if (!$wxCouponInfo) {
            WxCoupon::query()->create($wxCouponData);
        }
    }

    /**
     * @param $openId
     * @param $userId
     * @param $unionId
     */
    protected function updateOrderCouponUserId($openId, $userId, $unionId)
    {
        $orderList = Order::query()->where('open_id', $openId)->where('user_id', 0)
            ->pluck('id')
            ->toArray();
        if (!empty($orderList)) {
            Order::query()->whereIn('id', $orderList)
                ->update([
                    'user_id'  => $userId,
                    'union_id' => $unionId,
                ]);
        }

        $userCouponList = UserCoupon::query()->where('open_id', $openId)->where('user_id', 0)
            ->pluck('id')
            ->toArray();

        if (!empty($userCouponList)) {
            UserCoupon::query()->whereIn('id', $userCouponList)
                ->update([
                    'user_id'  => $userId,
                    'union_id' => $unionId,
                ]);
        }
    }

    /**
     * syncCrsOrder
     * @desc crs订单同步
     * @param $payload
     * @date 2020-09-09
     *
     * @throws \Exception
     */
    protected function syncCrsOrder($payload)
    {
        Log::info('syncCrsOrder');
        $data = json_decode($payload, true);
        $this->syncOrderData($data, 1);
    }

    /**
     * syncHwjOrder
     * @desc 好物集订单同步
     * @param $payload
     * @date 2020-09-09
     *
     * php artisan queue:work --queue=test-work-sync-hwj-order --delay=0 --sleep=3 --tries=3 --memory=1024 --env=$APP_ENV
     *
     * @throws \Exception
     */
    protected function syncHwjOrder($payload)
    {
        Log::info("syncHwjOrder");
        $data = json_decode($payload, true);
        $data = $this->formatOrderData($data);
        $this->syncOrderData($data, 2);
    }

    //格式化代码
    protected function formatOrderData($data){
        $thirdPartUserId = $data['mainOrderInfo']['shareUserId'];
        $model = new CrsWorkerThirdPartRelation();
        $haoWuJiEnterpriseId = env('HAO_WU_JI_ENTERPRISE_ID');
        $data['orderInfo']['store_staff_id'] = $model->getStoreStaffIdByThirdPartUserId(
            $thirdPartUserId,
            $haoWuJiEnterpriseId
        );
        return $data;
    }

    /**
     * syncOrderData
     * @desc 同步订单数据
     * @param $data
     * @date 2020-09-08
     *
     * @throws \Exception
     */
    protected function syncOrderData($data, $syncSource = 1)
    {
        Log::info('syncOrderData', $data);
        try {

            // 判断订单数据是否已存在
            $orderInfo = $data['orderInfo'];
            $existOrder = CrsOrder::getOrder($orderInfo['orderNo']);

            $updateOrderData = $this->convertOrderData($data);

            switchAppHost(['enterprise_id' => $updateOrderData['enterprise_id']]);

            // 新订单数据，创建订单和订单商品数据
            $updateOrderData['sync_source'] = $syncSource;

            //查询参数中的store_staff_id是否合伙人，如果是，则更新store_staff_id与partner_id
            if ($updateOrderData['store_staff_id'] > 0) {
                $partnerInfo = (new PartnerMemberService())->getPartnerInfo([
                    'enterprise_id' => $updateOrderData['enterprise_id'],
                    'id' => $updateOrderData['store_staff_id']
                ]);
                //存在合伙人ID，通过关系查询关联导购的ID
                if (!empty($partnerInfo)) {
                    $updateOrderData['partner_id'] = $partnerInfo['id'];
                    $updateOrderData['store_staff_id'] = $partnerInfo['worker_id'];
                }
            }

            // 如果是好物集订单则更新推广人
            if(CrsWorkerThirdPartRelation::SYNC_SOURCE_HWJ == $syncSource){
                $params = [
                    'order_no'       => $updateOrderData['order_no'],
                    'store_id'       => $updateOrderData['enterprise_id'],
                    'store_staff_id' => $updateOrderData['store_staff_id'],
                ];
                (new \CrsSdk\Shop\Order())->updateOrderStoreStaff($params);
            };

            if (empty($existOrder)) {
                $updateOrderData['id'] = idGenerate();
                CrsOrder::query()->create($updateOrderData);
                $insertGoodsData = $this->convertOrderGoodData($data);
                array_walk($insertGoodsData, function ($value) {
                    $value['id'] = idGenerate();
                    CrsOrderGoods::query()->create($value);
                });
            } else {

                //判定当是好物集订单后，不在对syncSource做处理
                if(CrsWorkerThirdPartRelation::SYNC_SOURCE_HWJ == $existOrder->sync_source){
                    $updateOrderData['sync_source'] = CrsWorkerThirdPartRelation::SYNC_SOURCE_HWJ;
                };
                // 订单数据已存在，全量更新订单数据
                CrsOrder::query()
                    ->where('order_no', $orderInfo['orderNo'])
                    ->update($updateOrderData);
            }

            if (!empty($data['orderStatusTrace']) && !empty($data['orderStatusTrace']['id'])) {
                if (empty($this->getOrderStatusTrace($data['orderStatusTrace']['id']))) {
                    $insertOrderStatusTrace = $this->convertOrderStatusTraceData($data);
                    CrsOrderStatusTrace::query()->create($insertOrderStatusTrace);
                }
            }

        } catch (\Exception $e) {
            Log::warning("同步订单数据异常: " . json_encode($e->getMessage()));
            throw $e;
        }
    }

    /**
     * getOrderStatusTrace
     * @desc 查询订单状态追踪表
     * @param $id
     *
     * @return mixed
     * @date 2020-09-08
     */
    protected function getOrderStatusTrace($id)
    {
        return CrsOrderStatusTrace::query()
            ->where('id', $id)
            ->first();
    }

    /**
     * getOrderGoods
     * @desc 查询订单商品数据
     * @param $orderNo
     *
     * @return mixed
     * @date 2020-09-08
     */
    protected function getOrderGoods($orderNo)
    {
        return CrsOrderGoods::query()
            ->where('order_no', $orderNo)
            ->first();
    }

    /**
     * convertOrderStatusTraceData
     * @desc 订单状态变化表数据组装
     * @param $data
     *
     * @return array
     * @date 2020-09-08
     */
    protected function convertOrderStatusTraceData($data)
    {
        return [
            'id' => $data['orderStatusTrace']['id'] ?? idGenerate(),
            'order_no' => $data['orderInfo']['orderNo'],
            'event' => $data['orderStatusTrace']['event'] ?? '',
            'event_desc' => $data['orderStatusTrace']['eventDesc'] ?? '',
            'from' => $data['orderStatusTrace']['from'] ?? 0,
            'to' => $data['orderStatusTrace']['to'] ?? 0,
            'change_at' => $data['orderStatusTrace']['createdAt'] ?? 0,
        ];
    }

    /**
     * convertOrderGoodData
     * @desc 组装订单商品数据
     * @param $data
     *
     * @return array
     * @date 2020-09-08
     */
    protected function convertOrderGoodData($data)
    {
        $goods = [];
        foreach ($data['orderGoods'] as $orderGood) {
            $goods[] = [
                'enterprise_id' => $data['orderInfo']['enterpriseId'],
                'order_no' => $data['orderInfo']['orderNo'],
                'sku_no' => $orderGood['skuNo'],
                'spu_id' => $orderGood['spuId'],
                'sku_id' => $orderGood['skuId'],
                'spu_name' => $orderGood['spuName'],
                'type' => $orderGood['type'],
                'goods_type' => $orderGood['goodsType'],
                'category_name' => $orderGood['categoryName'],
                'sku_picture' => $orderGood['skuPicture'],
                'sku_attribute' => $orderGood['skuAttribute'],
                'quantity' => $orderGood['quantity'],
                'original_price' => $orderGood['originalPrice'],
                'sell_price' => $orderGood['sellPrice'],
                'single_discount' => $orderGood['singleDiscount'],
                'sub_price' => $orderGood['subPrice'],
                'goods_pay_price' => $orderGood['goodsPayPrice']
            ];
        }
        return $goods;
    }

    /**
     * convertOrderData
     * @desc 组装订单表数据
     * @param $data
     *
     * @return array
     * @date 2020-09-08
     */
    protected function convertOrderData($data)
    {
        return [
            'enterprise_id' => $data['orderInfo']['enterpriseId'],
            'user_id' => $data['orderUser']['id'] ?? 0,
            'open_id' => $data['orderUser']['openId'],
            'union_id' => $data['orderUser']['unionId'],
            'order_no' => $data['orderInfo']['orderNo'],
            'status' => $data['orderInfo']['status'],
            'order_type' => $data['orderInfo']['type'],
            'payment_channel' => $data['orderInfo']['paymentChannel'],
            'total_amount' => $data['orderInfo']['totalAmount'],
            'preferential_price' => $data['mainOrderInfo']['preferentialPrice'],
            'final_amount' => $data['orderInfo']['finalAmount'],
            'pay_at' => $data['orderInfo']['payAt'],
            'refund_at' => $data['orderInfo']['refundAt'],
            'order_at' => $data['orderInfo']['createdAt'],
            'delivery_name' => $data['orderInfo']['deliveryName'],
            'delivery_mobile' => $data['orderInfo']['deliveryMobile'],
            'delivery_address' => $data['orderInfo']['deliveryAddress'],
            'remarks' => $data['orderInfo']['remarks'],
            'province' => $data['orderInfo']['province'],
            'city' => $data['orderInfo']['city'],
            'area' => $data['orderInfo']['area'],
            'logistic_identity' => $data['orderInfo']['logisticIdentity'],
            'logistic_no' => $data['orderInfo']['logisticNo'],
            'send_at' => $data['orderInfo']['sendAt'],
            'sign_at' => $data['orderInfo']['signAt'],
            'cash_fee' => $data['orderInfo']['cashFee'],
            'coupon_fee' => $data['orderInfo']['couponFee'],
            'coupon_count' => $data['orderInfo']['couponCount'],
            'other_fee' => $data['orderInfo']['otherFee'],
            'logistics_price' => $data['orderInfo']['logisticsPrice'] ?? 0,
            'third_payment_no' => $data['payInfo']['thirdPaymentNo'] ?? '',
            'store_staff_id' => $data['orderInfo']['storeStaffId'] ?? 0,
        ];
    }
}
