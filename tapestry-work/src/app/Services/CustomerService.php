<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 10:35
 */

namespace App\Services;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Library\Uc\External;
use App\Library\Uc\Member;
use App\Library\Uc\Tag;
use App\Models\Address;
use App\Models\CrsOrder;
use App\Models\MemberOrderStatistics;
use App\Models\MemberUserRelation;
use App\Models\Order;
use App\Models\StoreWorker;
use App\Models\User as UserModel;
use App\Models\UserCoupon;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\Order as ShopOrderSdk;
use CrsSdk\Shop\User;
use CrsSdk\Shop\WxCoupon;
use Exception;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\Log;
use Psr\Http\Message\StreamInterface;

class CustomerService
{
    // 导购订单类型
    const  WORKER_ORDER_TYPE = 3;

    /**
     * 用户订单
     *
     * @param array $where
     * @param int $limit
     * @param array $orderFields
     * @param array $goodsFields
     *
     * @return array
     */
    public static function orders(
        array $where,
        int   $limit = 20,
        array $orderFields
        = [
            "id",
            "order_no",
            "final_amount",
            "order_at",
            "cash_fee",
            "pay_at",
            "status",
        ],
        array $goodsFields
        = [
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
        ]
    )
    {
        $paginator = CrsOrder::query()
            ->where($where)
            ->when($orderFields, function ($query) use ($orderFields) {
                $query->select($orderFields);
            })
            ->with(['goods' => function ($query) use ($goodsFields) {
                $goodsFields && $query->select($goodsFields);
            }])
            ->orderBy('order_at', 'desc')
            ->paginate($limit);

        foreach ($paginator->items() as &$order) {
            $order->status_alias = CrsOrder::getStatusAlias($order->status); // 状态别名返回
            $order->created_at = $order->order_at;
            foreach ($order->goods as &$good) {
                $good->category_name = json_decode($good->category_name ?? '[]', true);
                $good->sku_attribute = json_decode($good->sku_attribute ?? '[]', true);
                $good->sku_attribute_items = implode(" ", array_column($good->sku_attribute, 'attribute_item_value'));
            }
        }
        return [
            'data' => $paginator->items(),
            'meta' => [
                'total' => (int)$paginator->total(),
                'page' => (int)$paginator->currentPage(),
                'size' => (int)$paginator->perPage(),
            ]
        ];
    }

    /**
     * 用户信息
     *
     * @param array $where
     * @param array $fields
     *
     * @return Builder|Model|mixed|object|null
     */
    public static function user(
        array $where,
        array $fields
        = [
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
        ]
    )
    {
        return UserModel::query()
            ->where($where)
            ->when($fields, function ($query) use ($fields) {
                $query->select($fields);
            })
            ->first();
    }

    /**
     * 优惠券
     *
     * @param array $where
     * @param int $limit
     * @param array $fields
     *
     * @return Builder[]|Collection|mixed[]
     */
    public static function coupons(array $where, int $limit = 20, array $fields = [])
    {
        $paginator = UserCoupon::query()
            ->where($where)
            ->when($fields, function ($query) use ($fields) {
                $query->select($fields);
            })
            ->paginate($limit);

        return [
            'data' => UserCoupon::format($paginator->items()),
            'meta' => [
                'total' => (int)$paginator->total(),
                'page' => (int)$paginator->currentPage(),
                'size' => (int)$paginator->perPage(),
            ]
        ];
    }

    /**
     * 获取地址
     *
     * @param array $where
     * @param array $fields
     *
     * @return Builder|Model|mixed|object|null
     */
    public static function address(array $where, array $fields = ['id', 'province', 'city', 'address'])
    {
        return Address::query()
            ->where($where)
            ->when($fields, function ($query) use ($fields) {
                $query->select($fields);
            })
            ->first();
    }

    /**
     * @param $domain
     * @param $pltUserId
     *
     * @return string
     * @throws Exception
     */
    public static function pltUserIdToEcUnionId($domain, $pltUserId)
    {
        if (!$pltUserId) {
            return '';
        }

        $env = config('app.env');
        $identifier = config('work_plugin.crs-guide')[$env]['identify'];

        try {
            $exchangeInfo = (new External())->exchangeInfo($domain, $identifier, ['pltUserid' => $pltUserId, 'isIgnoreSuiteId' => 1]);
        } catch (Exception $e) {
            Log::info('获取ecUnionId失败:' . $e->getMessage());
        }

        if (empty($exchangeInfo['epExternalUnionid'])) {
            throw new Error(10000, '企业微信未找到用户信息');
        }

        return $exchangeInfo['epExternalUnionid'];
    }

    public static function ecUserIdToPltUserId($domain, $ecUserId, $identifier)
    {
        if (!$ecUserId) {
            return '';
        }

        try {
            $exchangeInfo = (new External())->exchangeInfo($domain, $identifier, ['ecUserid' => $ecUserId]);
            return $exchangeInfo['pltUserid'];
        } catch (Exception $e) {
            Log::info('获取pltUserId失败:' . $e->getMessage());
        }
    }

    /**
     * getCustomerList
     * @desc 根据条件筛选顾客列表
     * @param $param
     *
     * @return array
     * @date 2020-07-16
     * @throws Exception
     * @author XuHuitao
     */
    public static function getCustomerList($param)
    {
        $memberInfo = StoreWorkerService::detail([
            'enterprise_id' => $param['enterprise_id'],
            'store_id' => $param['store_id'] ?? 0,
            'mobile' => $param['mobile'],
            'worker_status' => StoreWorker::IN_SERVICE
        ]);
        if (empty($memberInfo)) {
            throw new Error(600103, '店员信息不存在');
        }

        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $param['enterprise_id']])['data'];
        $param['qy_domain'] = $enterpriseInfo['domain'];
        $param['mem_uid'] = $memberInfo['mem_uid'];

        // 如果选择了企业标签，通过标签id从UC查询所有带有该标签的用户
        if (isset($param['tag_ids']) && !empty($param['tag_ids'])) {
            $param['tag_union_id'] = self::getUnionIdListFromEnterpriseTags($param['qy_domain'], $param['tag_ids']);
        }

        // 顾客列表只展示通过导购二维码添加的顾客
        $param['qr_code_id'] = $memberInfo['contact_qr_code_id'];

        // 选择了自建标签，通过自建标签查询带有该标签的用户
        if (isset($param['guide_tag_ids']) && !empty($param['guide_tag_ids'])) {
            $relationIds = UserCustomizeTagsService::getRelationIdsByTags($param['enterprise_id'], $memberInfo['id'], $param['guide_tag_ids']);
            $param['relation_id'] = $relationIds;
        }

        $param['worker_user_id'] = $memberInfo['user_id'];
        $param['worker_id'] = $memberInfo['id'];

        $customerList = MemberUserRelation::getMemberListWithOrder($param, 0, 10, true);
        $filterResult = self::filterCustomerListFromShop($customerList, $param);
        $customerList = $filterResult['list'];

        // 组装顾客数据
        $customerList = self::formatCustomerList($customerList, $param, $memberInfo);


        return [
            'customer_list' => $customerList,
            'worker_qr_code' => $memberInfo['qr_code'],
            'member_info' => [
                'mem_uid' => $memberInfo['mem_uid'],
                'is_boss' => $memberInfo['is_boss'],
                'type' => $memberInfo['type']
            ],
            'total' => $filterResult['total']
        ];
    }

    /**
     * getUnionIdListFromEnterpriseTags
     * @desc 通过企业标签列表查询对应标签下的用户
     * @param $qyDomain
     * @param $tagList
     *
     * @return array
     * @date 2020-08-02
     * @throws Exception
     */
    private static function getUnionIdListFromEnterpriseTags($qyDomain, $tagList)
    {
        $tagList = (new TagService())->getUserOfTagList($qyDomain, $tagList);
        $tagUserUnions = array_column($tagList, 'ecUnionid');
        $tagUserUnions = array_unique($tagUserUnions);
        return $tagUserUnions;
    }

    /**
     * filterCustomerListFromShop
     * @desc 调用gw_shop接口，根据筛选条件筛选顾客列表
     * @param $customerList
     * @param $param
     *
     * @return mixed|StreamInterface
     * @date 2020-08-04
     */
    private static function filterCustomerListFromShop($customerList, $param)
    {
        if (isset($customerList['data']) && isset($customerList['total'])) {
            $total = $customerList['total'];
            $customerList = $customerList['data'];
        }

        if (empty($customerList)) {
            return ['list' => [], 'total' => 0];
        }
        $unionIdList = [];
        foreach ($customerList as $value) {
            $unionIdList[] = [
                'relation_id' => $value['relation_id'],
                'union_id' => $value['union_id']
            ];
        }
        $param['customer_list'] = $unionIdList;
//        $result = (new Customer())->filterGuideCustomer($param);
//        $userList = array_column($result['list'], null, 'relation_id');
        //todo 兼容没有商城的尴尬
        $resultCustomers = [];
        foreach ($customerList as $key => $value) {
//            if (array_key_exists($value['relation_id'], $userList)) {
            $customerInfo = $value;
//                $customerInfo['mobile'] = $userList[$value['relation_id']]['user_info']['mobile'] ?? '';
//                $customerInfo['is_member'] = empty($userList[$value['relation_id']]['user_info']['third_member_no']) ? 2 : 1;
//                $customerInfo['third_member_no'] = $userList[$value['relation_id']]['user_info']['third_member_no'] ?? '';
//                $customerInfo['birthday'] = $userList[$value['relation_id']]['user_info']['birthday'] ?? '';
//                $customerInfo['real_name'] = $userList[$value['relation_id']]['user_info']['real_name'] ?? '';
//                $customerInfo['nick_name'] = $userList[$value['relation_id']]['user_info']['nick_name'] ?? $value['ec_name'];
//                $customerInfo['user_id'] = $userList[$value['relation_id']]['user_info']['id'] ?? 0;
//                $resultCustomers[] = $customerInfo;
            $customerInfo['mobile'] = '';
            $customerInfo['is_member'] = 2;
            $customerInfo['third_member_no'] = '';
            $customerInfo['birthday'] = '';
            $customerInfo['real_name'] = '';
            $customerInfo['nick_name'] = $value['ec_name'];
            $customerInfo['user_id'] = 0;
            $resultCustomers[] = $customerInfo;
//            } else {
//                unset($customerList[$key]);
//            }
        }
        return ['list' => $resultCustomers, 'total' => isset($total) ? $total : count($resultCustomers)];
    }

    /**
     * formatCustomerList
     * @desc 组装顾客列表数据
     * @param $customerList
     * @param $param
     * @param $memberInfo
     *
     * @return mixed
     * @date 2020-08-02
     * @throws Exception
     */
    private static function formatCustomerList($customerList, $param, $memberInfo)
    {
        if (empty($customerList)) {
            return [];
        }
        // 组装顾客数据
        $ecIdList = array_column($customerList, 'ec_id');
        $userTagList = self::getCustomerTagsFromUC($param['qy_domain'], $ecIdList);

        // 获取顾客自建标签
        $relationIdList = array_column($customerList, 'relation_id');
        $userGuideTags = UserCustomizeTagsService::getTagOfUsers($param['enterprise_id'], $memberInfo['id'], $relationIdList);
        $userGuideTagMap = [];
        foreach ($userGuideTags as $value) {
            $userGuideTagMap[$value['relation_id']][] = [
                'tagId' => $value['user_tags_id'],
                'tagName' => $value['tag_name'],
            ];
        }

        $exchange = false;
        $enterpriseIdentifierList = EnterpriseService::getEnterpriseIdentifierList($param['enterprise_id']);
        $identifierList = array_column($enterpriseIdentifierList, 'identifier');
        if (in_array(Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT, $identifierList)) {
            $identifier = Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT;
            if (env('APP_ENV') == 'dev') {
                $identifier = 'crs_store_assistant_dev';
            } else {
                $list = array_column($enterpriseIdentifierList, null, 'identifier');
                if ($list[$identifier]['created_at'] > env('ICLK_APP_CREATED_AT')) {
                    $identifier = 'iclk_crs_store_assistant';
                }
            }
            $exchange = true;
            Log::info('顾客ec_userid', array_column($customerList, 'ec_userid'));
            $customerExchangeInfoList = self::getCustomerExchangeList($param['qy_domain'],
                $identifier,
                array_column($customerList, 'ec_userid'));
            Log::info('顾客ec_userid转换后', $customerExchangeInfoList ?? []);
        }

        foreach ($customerList as $key => &$value) {
            $value['tags'] = array_values($userTagList[$value['ec_id']] ?? []);
            $value['guide_tags'] = $userGuideTagMap[$value['relation_id']] ?? [];
            $value['total_order_num'] = $value['total_order_num'] ?? 0;
            $value['total_order_amount'] = $value['total_order_amount'] ?? 0;
            if ($exchange) {
                $value['ec_userid'] = $customerExchangeInfoList[$value['ec_userid']]['pltUserid'] ?? $value['ec_userid'];
            }

            // 添加方式
            $value['add_type'] = ($value['source_type'] == $memberInfo['contact_qr_code_id']) ? Constant::ADD_TYPE_BUSINESS_CARD : Constant::ADD_TYPE_BUSINESS_CARD_OTHER;
        }
        return $customerList;
    }

    /**
     * getCustomerTagsFromUC
     *
     * @desc 从UC查询顾客标签信息
     *
     * @param $qyDomain
     * @param $ecIdList
     *
     * @return mixed
     * @date 2020-07-15
     * @throws Exception
     * @author XuHuitao
     */
    public static function getCustomerTagsFromUC($qyDomain, $ecIdList)
    {
        $info = (new Tag())->tagDetail($qyDomain, ['ecIds' => $ecIdList]);
        unset($info['requestId']);
        $tagList = [];
        array_walk($info, function ($item) use (&$tagList) {
            $tagList[$item['ecId']][$item['tagId']] = [
                'tagId' => $item['tagId'],
                'tagName' => $item['tagName'],
                'groupId' => $item['groupId'],
                'groupName' => $item['groupName']
            ];
        });
        return $tagList;
    }

    /**
     * getCustomerExchangeList
     *
     * @desc 通过外部联系人userId列表获取外部联系人平台userId列表
     *
     * @param $qyDomain
     * @param $identifier
     * @param $ecUserIdList
     *
     * @return array|mixed
     * @date 2020-07-26
     * @throws Exception
     */
    public static function getCustomerExchangeList($qyDomain, $identifier, $ecUserIdList)
    {
        $ecUserIdList = array_unique($ecUserIdList);
        $splitEcUserIds = array_chunk($ecUserIdList, Constant::UC_EXCHANGE_INTERFACE_PARAM_MAX);
        $result = [];
        foreach ($splitEcUserIds as $value) {
            $exchangeInfo = (new External())->exchangeInfoList($qyDomain, $identifier, ['ecUseridList' => $value]);
            if (empty($exchangeInfo)) {
                return [];
            }
            unset($exchangeInfo['requestId']);
            $exchangeInfo = array_column($exchangeInfo, null, 'epExternalUserid');
            $result[] = $exchangeInfo;
        }
        $result = array_merge([], ...$result);
        return $result;
    }

    /**
     * customerInfo
     *
     * @desc 获取顾客信息
     *
     * @param $enterpriseId
     * @param $memberInfo
     * @param $unionId
     * @param $qyDomain
     * @param $storeId
     * @param $relationId
     *
     * @return array|mixed|StreamInterface
     * @date 2020-07-17
     * @throws Exception
     */
    public static function customerInfo($enterpriseId, $memberInfo, $unionId, $qyDomain, $storeId, $relationId)
    {
        if ($relationId) {
            $customerRelation = MemberUserRelation::getDataById($relationId);
            if ($customerRelation) {
                $customerRelation = $customerRelation->toArray();
            } else {
                return [];
            }
            if ($customerRelation['mem_uid'] != $memberInfo['mem_uid']) {
                $memberInfo = StoreWorkerService::detail([
                    'enterprise_id' => $enterpriseId,
                    'store_id' => $storeId ?? 0,
                    'mem_uid' => $customerRelation['mem_uid'],
                    'worker_status' => StoreWorker::IN_SERVICE
                ]);
            }
        } else {
            if (empty($unionId)) {
                return [];
            }
            $customerRelation = MemberUserRelation::getCustomerLatestGuide([$memberInfo['mem_uid']], [$unionId]);
            $customerRelation = $customerRelation[$unionId] ?? [];
            if (empty($customerRelation) &&
                $memberInfo['type'] == Constant::STORE_WORKER_TYPE_STORE &&
                $memberInfo['is_boss'] == Constant::YES) {
                // 没有查询到店员和顾客关联数据，店长身份，查询该顾客在当前店铺下所有店员中最新的关联数据
                $memberList = self::getStoreWorkers($enterpriseId, $storeId)['list'];
                $customerRelation = MemberUserRelation::getCustomerLatestGuide(
                    array_column($memberList, 'mem_uid'),
                    [$unionId]);
                $customerRelation = $customerRelation[$unionId] ?? [];
                $memberList = array_column($memberList, null, 'mem_uid');
                $memberInfo = $memberList[$customerRelation['mem_uid']] ?? [];
            }
        }

        if (!empty($unionId)) {
            $customerInfo = self::getCustomerListFromShop($enterpriseId, [$unionId]);
            if (empty($customerInfo)) {
                $customerInfo = ['id' => 0];
            } else {
                $customerInfo = $customerInfo[0];
            }
        } else {
            $customerInfo = ['id' => 0];
        }

        $orderData = MemberOrderStatistics::getDataByUserId($customerInfo['id'], $enterpriseId, $unionId);
        $tagInfo = self::getCustomerTagsFromUC($qyDomain, [$customerRelation['ec_id']]);

        // 第三方应用更换ex_userid
        $enterpriseIdentifierList = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);
        $identifierList = array_column($enterpriseIdentifierList, 'identifier');

        if (in_array(Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT, $identifierList)) {
            $identifier = Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT;
            if (env('APP_ENV') == 'dev') {
                $identifier = 'crs_store_assistant_dev';
            } else {
                $list = array_column($enterpriseIdentifierList, null, 'identifier');
                if ($list[$identifier]['created_at'] > env('ICLK_APP_CREATED_AT')) {
                    $identifier = 'iclk_crs_store_assistant';
                }
            }
            $exchangeInfo = (new External())->exchangeInfo($qyDomain,
                $identifier,
                ['ecUserid' => $customerRelation['ec_userid']]);
            if (!$exchangeInfo) {
                throw new Error(10000, '用户信息查询失败');
            }
            $customerRelation['ec_userid'] = $exchangeInfo['pltUserid'] ?? $customerRelation['ec_userid'];
        }
        $customerData = [
            'ec_name' => $customerRelation['ec_name'] ?? '',
            'ec_avatar' => $customerRelation['ec_avatar'] ?? '',
            'ec_gender' => $customerRelation['ec_gender'] ?? 0,
            'ec_userid' => $customerRelation['ec_userid'],
            'mobile' => $customerInfo['mobile'] ?? '',
            'tags' => array_values($tagInfo[$customerRelation['ec_id']] ?? []),
            'total_order_num' => $orderData['total_order_num'] ?? 0,
            'total_order_amount' => $orderData['total_order_amount'] ?? 0,
            'relation_at' => $customerRelation['relation_at'] ?? 0,
            'user_id' => $customerInfo['id'],
            'mem_uid' => $customerRelation['mem_uid'] ?? '',
            'mem_name' => $memberInfo['name'] ?? ''
        ];
        return $customerData;
    }

    /**
     * getStoreWorkers
     *
     * @desc 获取门店下所有导购
     *
     * @param $enterpriseId
     * @param $storeId
     *
     * @return mixed
     * @date 2020-07-15
     * @throws Exception
     * @author XuHuitao
     */
    public static function getStoreWorkers($enterpriseId, $storeId)
    {
        return StoreWorkerService::list([
            'enterprise_id' => $enterpriseId,
            'store_id' => $storeId], false);
    }

    /**
     * getCustomerListFromShop
     * @desc 从shop查询顾客信息
     * @param $enterpriseId
     * @param $unionIds
     * @param $searchKeyWords
     *
     * @return mixed|StreamInterface
     * @date 2020-07-16
     * @author XuHuitao
     */
    public static function getCustomerListFromShop($enterpriseId, $unionIds, $searchKeyWords = '')
    {
        $param = [
            'enterprise_id' => $enterpriseId
        ];
        if (!empty($unionIds) && is_array($unionIds)) {
            $param['union_ids'] = $unionIds;
        }
        if (!empty($searchKeyWords)) {
            $param['search']['keywords'] = $searchKeyWords;
        }
        $userList = (new User())->users($param)['data'];
        return $userList;
    }

    /**
     * searchCustomer
     * @desc 判断搜索关键字分别通过手机号和姓名进行顾客搜索，返回
     * @param $enterpriseId
     * @param $searchKey
     *
     * @return array|Collection|mixed|StreamInterface|static[]
     * @date 2020-07-26
     */
    public static function searchCustomer($enterpriseId, $searchKey)
    {
        //$shopCustomerList = self::getCustomerListFromShop($enterpriseId, [], $searchKey);
        $shopCustomerList = [];
        $shopCustomerList = array_column($shopCustomerList, null, 'union_id');
        $relationList = MemberUserRelation::searchCustomerByName($enterpriseId, $searchKey, ['union_id', 'ec_name', 'id']);
        $relationList = $relationList->toArray();
        $relationList = array_column($relationList, null, 'id');
        return ['shop_customer' => $shopCustomerList, 'relations' => $relationList];
    }

    /**
     * 订单列表
     *
     * @param array $orderSearchParams
     *
     * @return array
     * @throws Exception
     */
    public static function customerOrder(array $orderSearchParams)
    {
        $returnData = Constant::LIST_DEFAULT_JSON;
        try {
            $orders = (new ShopOrderSdk())->orderList($orderSearchParams);
            $returnData['meta']['total'] = $orders['total'];
            $returnData['meta']['page'] = $orders['page'];
            $returnData['meta']['size'] = $orders['size'];
            $returnData['data'] = $orders['list'];

            if (!empty($orders['list'])) {
                array_walk($returnData['data'], function (&$order) {
                    $order['status_alias'] = Order::getStatusAlias($order['status']);
                    $order['goods'] = $order['goods_data'];
                    foreach ($order['goods'] as &$goodsItem) {
                        $goodsItem['sku_attribute_items'] = implode(" ", $goodsItem['sku_attribute']);
                    }
                    unset($order['goods_data']);
                });
            }
        } catch (Error $error) {
            Log::info("Shop:Order->orderList，错误信息：{$error->getMessage()},错误code：{$error->getCode()}", $orderSearchParams);
        }

        return $returnData;
    }

    /**
     * 用户优惠券
     *
     * @param int $enterpriseId
     * @param int $userId
     * @param int $couponStatus
     *
     * @return array
     * @throws Exception
     */
    public static function customerCoupons(int $enterpriseId, int $userId, int $couponStatus = 0): array
    {
        $coupons = Constant::LIST_DEFAULT_JSON;
        $couponParams = [
            'enterprise_id' => $enterpriseId,
            'user_id' => $userId,
            'status' => $couponStatus
        ];
        $list = [];
        try {
            $wxCoupons = (new WxCoupon())->userCoupons($couponParams);
            foreach ($wxCoupons['list'] as $coupon) {
                switch ($coupon['effective']) {
                    case UserCoupon::COUPON_LIMIT_NO:
                        $item = [
                            'id' => $coupon['coupon_user_id'],
                            'status' => $coupon['status'],
                            'effective' => $coupon['effective'],
                            'receive_time' => date('Y.m.d H:i', $coupon['receive_time']),
                            'preferential_rule' => '满' . ($coupon['coupon_mininumn'] / 100) . '减' . ($coupon['coupon_value'] / 100),
                            'preferential_price' => ($coupon['coupon_value'] / 100) . '元',
                            'preferential_time_limit' => '长期有效',
                            'order_amount' => empty($coupon['order_amount']) ? 0 : ($coupon['order_amount'] / 100),
                            'order_no' => $coupon['order_no'] ?? '',
                            'order_at' => (UserCoupon::COUPON_HAVE_USED == $coupon['status'] && isset($coupon['pay_at']) && $coupon['pay_at'] > 0) ? date('Y.m.d H:i', $coupon['pay_at']) : '',
                            'use_time' => UserCoupon::COUPON_HAVE_USED == $coupon['status'] ? date('Y.m.d H:i', $coupon['use_time']) : '',
                        ];
                        break;
                    case UserCoupon::COUPON_LIMIT_TINE;
                        $item = [
                            'id' => $coupon['coupon_user_id'],
                            'status' => UserCoupon::COUPON_WAIT_USE == $coupon['status'] ? (time() >= $coupon['end_time'] ? UserCoupon::COUPON_HAVE_EXPIRED : $coupon['status']) : $coupon['status'],
                            'effective' => $coupon['effective'],
                            'receive_time' => date('Y.m.d H:i', $coupon['receive_time']),
                            'preferential_rule' => '满' . ($coupon['coupon_mininumn'] / 100) . '减' . ($coupon['coupon_value'] / 100),
                            'preferential_price' => ($coupon['coupon_value'] / 100) . '元',
                            'preferential_time_limit' => date('Y.m.d', $coupon['start_time']) . '-' . date('Y.m.d', $coupon['end_time']),
                            'order_amount' => empty($coupon['order_amount']) ? 0 : ($coupon['order_amount'] / 100),
                            'order_no' => $coupon['order_no'] ?? '',
                            'use_time' => UserCoupon::COUPON_HAVE_USED == $coupon['status'] ? date('Y.m.d H:i', $coupon['use_time']) : '',
                            'order_at' => (UserCoupon::COUPON_HAVE_USED == $coupon['status'] && isset($coupon['pay_at']) && $coupon['pay_at'] > 0) ? date('Y.m.d H:i', $coupon['pay_at']) : ''
                        ];
                        break;
                    default:
                        $item = [];
                }

                $list[] = $item;
            }
            $coupons['data'] = $list;
            $coupons['meta']['total'] = $wxCoupons['total'];
            $coupons['meta']['page'] = $wxCoupons['page'];
            $coupons['meta']['size'] = $wxCoupons['size'];
        } catch (Error $error) {
            Log::error("", $couponParams);
        }

        return $coupons;
    }

    /**
     * statisticCustomerData
     *
     * @desc 根据导购列表统计顾客数据
     *
     * @param $enterpriseId
     * @param $workerList
     * @param $startTime
     * @param $endTime
     *
     * @return array
     * @date 2020-08-14
     */
    public static function statisticCustomerData($enterpriseId, $workerList, $startTime, $endTime)
    {
        $accessCustomerList = $historyUnionIds = $guideHistoryUnionIds = $otherHistoryUnionIds = $dailyData = $dealCustomerByTime = $dailyUnionId = $totalUnionIds = $totalDealUnionIds = $totalAccessUnionIds = [];
        $userService = new User();
        //获取筛选导购的二维码id
        $codeIdList = array_filter(array_column($workerList, 'contact_qr_code_id'));
        sort($codeIdList);
        //获取导购MemUid
        $memUidList = array_filter(array_column($workerList, 'mem_uid'));
        sort($memUidList);
        //获取开始时间前的顾客列表
        $historyCustomerList = MemberUserRelation::statisticCustomerList($enterpriseId, $memUidList, 0, $startTime);
        //获取开始时间前的顾客unionId，导购名片顾客unionId，其他方式顾客unionId
        foreach ($historyCustomerList as $item) {
            $historyUnionIds = array_merge($historyUnionIds, [$item['union_id']]);
            if (in_array($item['source_type'], $codeIdList)) {
                $guideHistoryUnionIds = array_merge($guideHistoryUnionIds, [$item['union_id']]);
            } else {
                $otherHistoryUnionIds = array_merge($otherHistoryUnionIds, [$item['union_id']]);
            }
        }
        //获取开始时间前的顾客总数，导购名片添加顾客总数，其他方式添加顾客总数
        $historyCustomer = count($historyUnionIds);
        $guideHistoryCustomer = count($guideHistoryUnionIds);
        $otherHistoryCustomer = count($otherHistoryUnionIds);
        //获取查询范围内的顾客列表
        $timeRangeResult = MemberUserRelation::statisticCustomerList($enterpriseId, $memUidList, $startTime, $endTime);
        //组合每日新增顾客，每日新增顾客数
        foreach ($timeRangeResult as &$value) {
            $created_date = date("Y-m-d", $value['created_date']);
            if (!isset($dailyData[$created_date])) {
                $dailyData[$created_date] = [
                    'created_date' => $value['created_date'],
                    'union_ids' => [],
                    'guide_union_ids' => [],
                    'other_union_ids' => [],
                    'customer_new' => 0,
                    'guide_customer_new' => 0,
                    'other_customer_new' => 0,
                ];
            }
            //每日新增顾客unionId
            $dailyData[$created_date]['union_ids'] = array_merge($dailyData[$created_date]['union_ids'], [$value['union_id']]);
            //导购名片每日新增顾客unionId
            if (in_array($value['source_type'], $codeIdList)) {
                $dailyData[$created_date]['guide_union_ids'] = array_merge($dailyData[$created_date]['guide_union_ids'], [$value['union_id']]);
            } else {
                //其他方式每日新增顾客unionId
                $dailyData[$created_date]['other_union_ids'] = array_merge($dailyData[$created_date]['other_union_ids'], [$value['union_id']]);
            }
            //每日新增顾客数，导购名片每日新增顾客数，其他方式每日新增顾客数量
            $dailyData[$created_date]['customer_new'] = count($dailyData[$created_date]['union_ids']);
            $dailyData[$created_date]['guide_customer_new'] = count($dailyData[$created_date]['guide_union_ids']);
            $dailyData[$created_date]['other_customer_new'] = count($dailyData[$created_date]['other_union_ids']);
        }
        //新增顾客数
        $newCustomer = array_sum(array_column($dailyData, 'customer_new'));
        //导购名片新增顾客数
        $guideNewCustomer = array_sum(array_column($dailyData, 'guide_customer_new'));
        //其他方式新增顾客数
        $otherNewCustomer = array_sum(array_column($dailyData, 'other_customer_new'));

        //初始化统计
        $detailData = self::dailyDataFormat($startTime, $endTime);

        //顾客总数
        $dailyTotalCustomer = $historyCustomer;
        //导购名片添加总数
        $dailyGuideCustomer = $guideHistoryCustomer;
        //其他方式添加总数
        $dailyOtherCustomer = $otherHistoryCustomer;

        foreach ($detailData as &$item) {
            $item['data']['customer_new'] = isset($dailyData[$item['time']]['customer_new']) ? $dailyData[$item['time']]['customer_new'] : 0;
            $item['data']['guide_customer_new'] = isset($dailyData[$item['time']]['guide_customer_new']) ? $dailyData[$item['time']]['guide_customer_new'] : 0;
            $item['data']['other_customer_new'] = isset($dailyData[$item['time']]['other_customer_new']) ? $dailyData[$item['time']]['other_customer_new'] : 0;
            $dailyTotalCustomer += $item['data']['customer_new'];
            $dailyGuideCustomer += $item['data']['guide_customer_new'];
            $dailyOtherCustomer += $item['data']['other_customer_new'];
            $item['data']['customer_total'] = $dailyTotalCustomer;
            $item['data']['guide_customer_total'] = $dailyGuideCustomer;
            $item['data']['other_customer_total'] = $dailyOtherCustomer;
            //按日期获取浏览商城顾客数
            $dailyUnionId = isset($dailyData[$item['time']]['union_ids']) ? $dailyData[$item['time']]['union_ids'] : [];
            //合并每天的数据去查询浏览顾客，成交顾客
            $totalUnionIds = array_merge($dailyUnionId, $totalUnionIds);
            $totalUnionIds = array_merge($historyUnionIds, $totalUnionIds);
            //去掉空的union_id
            $totalUnionIds = array_filter($totalUnionIds);
            sort($totalUnionIds);
            //获取筛选时间段内浏览商城顾客数
            if (!empty($totalUnionIds)) {
                $params['enterprise_id'] = $enterpriseId;
                $params['day_time'] = $item['day_time'];
                $params['union_ids'] = $totalUnionIds;
                $accessCustomerList = $userService->dailyAccessRecords($params)['list'];
            }
            $item['data']['access_union_ids'] = isset($accessCustomerList[$item['day_time']]['union_ids']) ? $accessCustomerList[$item['day_time']]['union_ids'] : [];
            $totalAccessUnionIds = array_merge($totalAccessUnionIds, $item['data']['access_union_ids']);
            $item['data']['customer_access'] = isset($accessCustomerList[$item['day_time']]['num']) ? $accessCustomerList[$item['day_time']]['num'] : 0;
            //获取成交顾客
            if (!empty($totalUnionIds)) {
                $dealCustomerByTime = self::orderOutStatistics($totalUnionIds, strtotime($item['time']), strtotime($item['time']) + 86399);
            }
            $totalDealUnionIds = array_merge(array_unique($dealCustomerByTime), $totalDealUnionIds);
            $item['data']['customer_deal'] = count(array_unique($dealCustomerByTime));
            unset($item['data']['access_union_ids']);
        }
        $detailData = array_column($detailData, 'data', 'time');
        $totalAccessCustomer = count(array_unique($totalAccessUnionIds));
        $totalDealCustomer = count(array_unique($totalDealUnionIds));
        return [
            'customer_total' => $dailyTotalCustomer,  //顾客总数
            'guide_customer_total' => $dailyGuideCustomer,  //导购名片添加总数
            'other_customer_total' => $dailyOtherCustomer,  //其他方式添加总数
            'customer_new' => $newCustomer,         //新增顾客数
            'guide_customer_new' => $guideNewCustomer,    //导购名片新增顾客数
            'other_customer_new' => $otherNewCustomer,    //其他方式新增顾客数
            'customer_deal' => $totalDealCustomer,   //成交顾客总数
            'customer_access' => $totalAccessCustomer, //浏览商城顾客数
            'detail_data' => $detailData,
        ];
    }

    /**
     * dailyDataFormat
     *
     * @desc 初始化时间范围内数据
     *
     * @param $startTime
     * @param $endTime
     *
     * @return array
     * @date 2020-08-18
     */
    private static function dailyDataFormat($startTime, $endTime)
    {
        $timeIndex = $startTime;
        $result = [];
        $initData = [
            'customer_total' => 0,  //顾客总数
            'guide_customer_total' => 0,  //导购名片顾客数
            'other_customer_total' => 0,  //其他方式顾客数
            'customer_access' => 0,  //浏览商城顾客数
            'customer_new' => 0,   //新增顾客数
            'guide_customer_new' => 0,  //导购名片新增顾客数
            'other_customer_new' => 0,  //其他方式新增顾客数
            'customer_deal' => 0,  //成交顾客数
        ];
        while ($timeIndex < $endTime) {
            $initData['time'] = date('Y-m-d', $timeIndex);
            $result[] = [
                'time' => date('Y-m-d', $timeIndex),
                'data' => $initData,
                'day_time' => date('Ymd', $timeIndex)
            ];
            $timeIndex += Constant::ONE_DAY_SEC;
        }
        return $result;
    }

    /**
     * 查询顾客对应数据
     * @param array $customerUnionIds
     * @param int $startTime
     * @param int $endTime
     *
     * @return array
     */
    public static function orderOutStatistics($customerUnionIds = [], $startTime = 0, $endTime = 0)
    {
        $payedUnionList = [];
        if (count($customerUnionIds) > 1000) {

            $groupUserIds = array_chunk($customerUnionIds, 900);
            foreach ($groupUserIds as $userIdList) {
                // 获取导购成交用户数
                $result = CrsOrder::query()
                    ->whereIn('union_id', $userIdList)
                    ->where('pay_at', '>', 0)
                    ->when(!empty($startTime), function ($query) use ($startTime) {
                        return $query->where('created_at', '>=', $startTime);
                    })
                    ->when(!empty($endTime), function ($query) use ($endTime) {
                        return $query->where('created_at', '<=', $endTime);
                    })
                    ->get(['union_id'])
                    ->toArray();
                $payedUnionList = array_merge($payedUnionList, array_column($result, 'union_id') ?? []);
            }
        } else {

            // 获取导购成交用户数
            $result = CrsOrder::query()
                ->whereIn('union_id', $customerUnionIds)
                ->where('pay_at', '>', 0)
                ->when(!empty($startTime), function ($query) use ($startTime) {
                    return $query->where('created_at', '>=', $startTime);
                })
                ->when(!empty($endTime), function ($query) use ($endTime) {
                    return $query->where('created_at', '<=', $endTime);
                })
                ->get(['union_id'])
                ->toArray();


            $payedUnionList = array_column($result, 'union_id') ?? [];
        }

        return $payedUnionList;
    }

    /**
     * 根据导购获取顾客数据
     *
     * @param $enterpriseId
     * @param $workerList
     * @param $startTime
     * @param $endTime
     *
     * @return mixed
     */
    public static function statisticCustomerDataByGuide($enterpriseId, $workerList, $startTime, $endTime)
    {
        $customerData = $payedUnionList = $customerSourceList = [];
        $guideMemUidList = array_column($workerList, 'mem_uid');
        $result = MemberUserRelation::statisticCustomerByGuide($enterpriseId, $guideMemUidList, $startTime, $endTime);
        $timeData = $result['time_data'];
        $totalData = $result['total_data'];
        // 获取截止时间前的顾客数据
        $customerList = self::getCustomerMemberList($enterpriseId, $guideMemUidList, 0, $endTime);
        if (!empty($customerList)) {
            $customerUnionIds = array_column($customerList, 'union_id');
            $payedUnionList = self::orderOutStatistics($customerUnionIds, $startTime, $endTime);
        }
        // 获取每个导购对应的union_id集合
        foreach ($customerList as $item) {
            $customerData[$item['mem_uid']][] = $item['union_id'];
            if (!empty($customerSourceList[$item['mem_uid']]['num'])) {
                $customerSourceList[$item['mem_uid']]['num'] += in_array($item['union_id'], $payedUnionList) ? 1 : 0;
            } else {
                $customerSourceList[$item['mem_uid']]['num'] = in_array($item['union_id'], $payedUnionList) ? 1 : 0;
            }
        }

        $shopUser = (new User());
        $customerInfoData = [];
        foreach ($customerData as $key => $unionIds) {
            $totalAccessCustomer = 0;
            $unionIds = array_filter(array_unique($unionIds));
            if (!empty($unionIds)) {
                $params['enterprise_id'] = $enterpriseId;
                $params['start_time'] = $startTime;
                $params['end_time'] = $endTime;
                $params['union_ids'] = $unionIds;
                $accessCustomerList = $shopUser->dailyAccessRecords($params)['list'];

                if (!empty($accessCustomerList)) {
                    $accessCustomerUnionIds = [];
                    foreach ($accessCustomerList as $item) {
                        $accessCustomerUnionIds = array_merge($accessCustomerUnionIds, $item['union_ids']);
                    }
                    $totalAccessCustomer = count(array_unique($accessCustomerUnionIds)) ?? 0;
                }
            }
            $customerInfoData[$key] = [
                'customer_deal' => $customerSourceList[$key]['num'] ?? 0, // 顾客成交数量
                'customer_access' => $totalAccessCustomer,   //浏览商城顾客数
            ];
        }

        foreach ($workerList as &$value) {
            $guideCustomerTotal = $otherCustomerTotal = $guideCustomerNew = $otherCustomerNew = 0;
            $customerInfo = $customerInfoData[$value['mem_uid']] ?? [];
            $data['customer_total'] = isset($totalData[$value['mem_uid']]) ? count($totalData[$value['mem_uid']]) : 0;
            if ($data['customer_total'] > 0) {
                $totalMemUidList = array_column($totalData[$value['mem_uid']], 'source_type');
                $guideCustomerTotal = isset(array_count_values($totalMemUidList)[$value['contact_qr_code_id']]) ? array_count_values($totalMemUidList)[$value['contact_qr_code_id']] : 0;
                $otherCustomerTotal = (count($totalMemUidList) - $guideCustomerTotal) > 0 ? count($totalMemUidList) - $guideCustomerTotal : 0;
            }
            $data['guide_customer_total'] = $guideCustomerTotal;
            $data['other_customer_total'] = $otherCustomerTotal;

            $data['customer_new'] = isset($timeData[$value['mem_uid']]) ? count($timeData[$value['mem_uid']]) : 0;
            if ($data['customer_new'] > 0) {
                $timeMemUidList = array_column($timeData[$value['mem_uid']], 'source_type');
                $guideCustomerNew = isset(array_count_values($timeMemUidList)[$value['contact_qr_code_id']]) ? array_count_values($timeMemUidList)[$value['contact_qr_code_id']] : 0;
                $otherCustomerNew = (count($timeMemUidList) - $guideCustomerNew) > 0 ? count($timeMemUidList) - $guideCustomerNew : 0;
            }

            $data['guide_customer_new'] = $guideCustomerNew;
            $data['other_customer_new'] = $otherCustomerNew;
            $data['customer_deal'] = $customerInfo['customer_deal'] ?? 0;
            $data['customer_access'] = $customerInfo['customer_access'] ?? 0;
            $data['id'] = $value['id'];
            $data['worker_no'] = $value['worker_no'];
            $data['name'] = $value['name'];
            $data['store_id'] = $value['store_id'];
            $data['contact_qr_code_id'] = $value['contact_qr_code_id'];
            $data['avatar'] = $value['avatar'];
            $data['mem_uid'] = $value['mem_uid'];
            $value = $data;
        }

        return $workerList;
    }

    /**
     * @param int $enterpriseId
     * @param array $codeIdList
     * @param int $startTime
     * @param int $endTime
     * @param array $field
     *
     * @return array
     */
    public static function getCustomerMemberList($enterpriseId = 0, $memUidList = [], $startTime = 0, $endTime = 0, $field = ['*'])
    {
        return MemberUserRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('union_id', '<>', '')
            ->where('source_type', '>', 0)
            ->when(!empty($memUidList), function ($query) use ($memUidList) {
                return $query->whereIn('mem_uid', $memUidList);
            })
            ->where('delete_type', 0)
            ->when(!empty($startTime) && !empty($endTime), function ($query) use ($startTime, $endTime) {
                return $query->whereBetween('created_at', [$startTime, $endTime]);
            })
            ->select($field)
            ->get()
            ->toArray();
    }

    /**
     * 获取客户详情
     *
     * @param int $enterpriseId 企业id
     * @param string $qyDomain 企业domain
     * @param string $appIdentifier 企业应用标识
     * @param int $customerId 客户id（user表主键id）
     * @param string $ecId 外部联系人记录id
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerDetail($enterpriseId, $qyDomain, $appIdentifier, $customerId = 0, $ecId = '')
    {
        // 用户信息
        $userDetail = [];
        if ($customerId) {
            // 调用sdk获取用户信息
            $params = [
                'enterprise_id' => $enterpriseId,
                'id' => $customerId
            ];
            $userDetail = (new User())->detail($params)['data'];
        }

        // 标记联系人关系数据
        $relation = [];
        if ($ecId) {
            // 查询条件
            $where = [
                'enterprise_id' => $enterpriseId,
                'ec_id' => $ecId,
                'is_exclusive' => MemberUserRelation::IS_EXCLUSIVE_TRUE // 专属导购
            ];
            // 专属导购mem_uid
            $relation = MemberUserRelation::where($where)->first();

            $relation = $relation ? $relation->toArray() : [];
        }

        // 数据格式化
        return $this->formatCustomerDetail([
            'enterprise_id' => $enterpriseId,
            'qy_domain' => $qyDomain,
            'app_identifier' => $appIdentifier,
            'relation' => $relation,
            'user_detail' => $userDetail
        ]);
    }

    /**
     * 客户详情数据格式化
     *
     * @param array $params 参数列表
     *                      + int enterprise_id 企业id
     *                      + string qy_domain 企业domain
     *                      + string app_identifier 企业应用标识
     *                      + array relation 员工-外部联系人标记关系
     *                      + array user_detail 用户信息
     *
     * @return array
     * @throws Exception
     */
    protected function formatCustomerDetail($params)
    {
        // 企业domain
        $qyDomain = $params['qy_domain'];
        // 应用标识
        $appIdentifier = $params['app_identifier'];

        // 员工-外部联系人标记关系
        $relation = $params['relation'];
        // 专属导购mem_uid赋值
        $memUid = $relation['mem_uid'] ?? '';

        // 专属导购名称
        $salesUsername = '';
        // 专属导购所属门店名称
        $salesDpName = '';
        if ($memUid) {
            // 调用uc sdk获取导购信息
            $member = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $memUid]);
            if ($member && 3 != $member['memStatus']) {
                $salesUsername = $member['memUsername'] ?? '';
                $salesDpName = $member['dpName'][0]['dpName'] ?? '';
            }
        }

        $ecName = '';
        $ecAvatar = '';
        $ecGender = 0;
        if (!empty($relation)) {
            $ecName = $relation['ec_name'];
            $ecAvatar = $relation['ec_avatar'];
            $ecGender = $relation['ec_gender'];
        }

        // 用户信息
        $userDetail = $params['user_detail'];

        return [
            'ec_id' => $relation['ec_id'] ?? '',
            'user_id' => $userDetail['id'] ?? 0,
            'union_id' => $userDetail['union_id'] ?? '',
            'nick_name' => $userDetail['nick_name'] ?? $ecName,
            'real_name' => $userDetail['real_name'] ?? '',
            'avatar' => $userDetail['avatar'] ?? $ecAvatar,
            'sex' => $userDetail['sex'] ?? $ecGender,
            'birthday' => $userDetail['birthday'] ?? '',
            'mobile' => $userDetail['mobile'] ?? '',
            'sales_username' => $salesUsername,
            'sales_dp_name' => $salesDpName,
            'register_at' => $userDetail['register_at'] ?? 0,
            'level_num' => $userDetail['level_num'] ?? 0,
            'level_name' => $userDetail['level_name'] ?? '',
            'level_is_highest' => $userDetail['level_is_highest'] ?? 0, // 用户当前会员是否最高等级会员（1：是，0：否）
            'expiration_at' => $userDetail['expiration_at'] ?? 0,
            'amount' => $userDetail['amount'] ?? 0
        ];
    }
}
