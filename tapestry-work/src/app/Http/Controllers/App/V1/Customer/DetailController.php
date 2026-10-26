<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\External;
use App\Models\Address;
use App\Models\MemberUserRelation;
use App\Models\Order;
use App\Models\User as UserModel;
use App\Models\UserCoupon;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class DetailController extends BaseController
{
    /**
     * 客户详情
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $userId    = $request->input('user_id', 0);
        $pltUserId = $request->input('plt_user_id', '');

        if (empty($userId) && empty($pltUserId)) {
            throw new Error(10000, '用户[user_id]不能为空');
        }

        $env        = config('app.env');
        $identifier = config('work_plugin.crs-guide')[$env]['identify'];

        $return    = [];
        $ecUnionId = '';
        if ($pltUserId) {
            $exchangeInfo = (new External())->exchangeInfo($this->qyDomain, $identifier, ['pltUserid' => $pltUserId, 'isIgnoreSuiteId' => 1]);
            if (!$exchangeInfo) {
                throw new Error(10000, '用户信息查询失败');
            }
            $ecUnionId = $exchangeInfo['epExternalUnionid'];
        }
        // 用户查询
        $user = UserModel::query()
            ->where([
                'enterprise_id' => $this->enterpriseId,
            ])
            ->when(!empty($ecUnionId), function ($query) use ($ecUnionId) {
                $query->where('union_id', $ecUnionId);
            })
            ->when(!empty($userId), function ($query) use ($userId) {
                $query->where('id', $userId);
            })
            ->first([
                'id',
                'union_id',
                'open_id',
                'sex',
                'nick_name',
                'real_name',
                'real_name_index',
                'mobile',
                'birthday',
                'avatar',
                'user_amount',
                'country',
                'province',
                'city',
                'register_at',
                'auth_at',
                'type',
            ]);

        if (!$user) {
            throw new Error(100000, '客户不存在');
        }

        //添加地址
        $addressInfo = Address::query()
            ->where([
                'user_id'       => $user['id'],
                'default'       => Address::DEF_ADDR,
                'enterprise_id' => $this->enterpriseId,
            ])
            ->first([
                'province',
                'city',
                'area',
                'address',
            ]);

        $address = '';
        if ($addressInfo) {
            $address = $addressInfo->province . $addressInfo->city . $addressInfo->area . $addressInfo->address;
        }

        $user->address = $address;

        if ($userId) {
            // 外部联系人id获取
            $ecUserId = MemberUserRelation::query()
                ->where([
                    'enterprise_id' => $this->enterpriseId,
                    'union_id'      => $user->union_id,
                    'mem_uid'       => $this->memUid,
                ])
                ->value('ec_userid');
            if (!empty($ecUserId)) {
                try {
                    $exchangeInfo = (new External())->exchangeInfo($this->qyDomain, $identifier, ['ecUserid' => $ecUserId, 'isIgnoreSuiteId' => 1]);
                    $pltUserId    = $exchangeInfo['pltUserid'];
                } catch (\Exception $e) {
                    Log::info('获取pltUserId失败:' . $e->getMessage());
                }
            }
        }

        $user->ec_user_id = $pltUserId;

        $return['user'] = $user;

        $orders = Order::query()
            ->where([
                'enterprise_id' => $this->enterpriseId,
                'user_id'       => $user['id'],
            ])
            ->with([
                'goods' => function ($query) {
                    $query->select([
                        "order_no",
                        "spu_id",
                        "sku_id",
                        "spu_name",
                        "category_name",
                        "sku_picture",
                        "sku_attribute",
                        "quantity",
                        "sell_price",
                        "sub_price"
                    ]);
                }
            ])
            ->orderBy('order_at', 'desc')
            ->get([
                "id",
                "order_no",
                "cash_fee as final_amount",
                "cash_fee",
                "order_at",
                "pay_at",
                "status",
            ])
            ->toArray();

        foreach ($orders as &$order) {
            // 状态别名返回
            $order['status_alias'] = order::getStatusAlias($order['status']);
            foreach ($order['goods'] as &$good) {
                $good['category_name'] = json_decode($good['category_name'] ?? '[]', true);
                $good['sku_attribute'] = json_decode($good['sku_attribute'] ?? '[]', true);
            }
        }

        $return['orders'] = $orders;

        return success(['data' => $return]);
    }

}
