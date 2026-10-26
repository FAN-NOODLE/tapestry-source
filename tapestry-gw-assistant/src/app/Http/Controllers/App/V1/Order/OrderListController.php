<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 17:37
 */

namespace App\Http\Controllers\App\V1\Order;


use App\Common\Constant;
use App\Models\StoreMemberViewTab;
use CrsSdk\Shop\Order;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Shop\User;
use CrsSdk\Work\Customer;
use CrsSdk\Work\PartnerMember;
use Illuminate\Support\Facades\Log;

class OrderListController extends OrderBaseController
{
    const ORDER_ALL = 'all';

    const ORDER_REFUND_STATUS = [41, 61, 81, 82, 84];

    const ORDER_RECEIVED_ARRAY = [25, 30, 40, 41, 60, 61, 80, 81, 83, 100, 200, 201];

    // 推广订单查询状态
    const ORDER_PROMOTER_STATUS = [40, 60, 41, 80, 81, 83, 100, 200, 205];

    // 我的订单
    const MY_ORDER = 2;
    // 待接订单
    const ORDER_TO_BE_RECEIVED = 1;

    //员工类型  0全部， 1、导购  2、合伙人
    const IDENTITY_ID_ARRAY = [0, 1, 2];

    /**
     * 不同状态下的订单列表
     * 订单状态
     * 需配送下 待发货 40 待收货 80 已签收 80
     * 需自提 待签收 60还是40 已签收 80
     * 待售后 21:待审核
     * 待退款 25:人工中
     * status 需配送：toBeDelivered，仅自提：onlyPickedUp
     * 待售后：afterSaleToBeChecked，待退款：toBeRefunded
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index()
    {
        $storeId = $this->request->get('store_id', 0);

        $input['status']             = $this->request->get('status', self::ORDER_ALL);
        $input['page']               = $this->request->get('page', Constant::PAGE);
        $input['size']               = $this->request->get('size', Constant::LIMIT);
        //订单类型：1推广订单   2服务订单
        $input['order_type']         = $this->request->get('order_type', 2);
        //服务订单类型：1待接订单
        $input['service_order_type'] = $this->request->get('service_order_type', 0);
        //推广人ids，格式如[234,345,5,4]
        $input['promoter_ids']       = $this->request->get('promoter_ids', []);
        $input['start_time']         = $this->request->get('start_time', '');
        $input['end_time']           = $this->request->get('end_time', '');
        //服务人ids，格式如[234,345,5,4]
        $input['waiter_ids']         = $this->request->get('waiter_ids', []);
        $input['order_status']       = $this->request->get('order_status', []);
        //订单快递类型：1快递   2自提  0全部
        $input['delivery_type']      = $this->request->get('delivery_type', '');
        //门店ids，格式如[234,345,5,4]
        $input['store_ids']          = $this->request->get('store_ids', []);
        //选中导购ID集合格式如[234,345,5,4]
        $input['check_worker_ids']   = $this->request->get('check_worker_ids', []);
        $input['check_all']          = $this->request->get('check_all', 0);
        //门店ID
        $input['store_id']           = !empty($storeId) ? [$storeId] : [];
        $input['enterprise_id']      = $this->enterpriseId;

        //初始promoter_ids,waiter_ids传的应该是mem_uid,由于老代码原因，这里先处理一下，方便后面使用
        $promoterMemUidArr = $input['promoter_ids'];
        $input['promoter_ids'] = [];
        $waiterMemUidArr = $input['waiter_ids'];
        $input['waiter_ids'] = [];

        // 如果不是推广订单则全部为空，这几个角色返回初始化数据
        if (in_array($this->role, [
                Constant::LOGIN_ROLE_ENTERPRISE,
                Constant::LOGIN_ROLE_COMMUNITY,
                Constant::LOGIN_ROLE_ENTERPRISE_ADMIN
            ]) && $input['order_type'] == self::MY_ORDER) {
            return success($this->getInitList($input));
        }

        $haoWuJiEnterpriseId = env('HAO_WU_JI_ENTERPRISE_ID');
        // 如果企业ID是好物品集合企业
        if (intval($this->enterpriseId) == $haoWuJiEnterpriseId
            &&
            $input['order_type'] !== Constant::EXTENSION_ORDER_TYPE) {
            return success($this->getInitList($input));
        }

        //员工类型：
        //0或者不带、全部
        //1、导购订单
        //2、合伙人订单
        $input['identity_type']      = $this->request->get('identity_type', 0);
        //检查identity_type
        if (!in_array($input['identity_type'], self::IDENTITY_ID_ARRAY)) {
            $input['identity_type'] = 0;
        }

        //(这个方法处理比较复杂，前面的人加的备注太少了，后面的记得在关键地方加备注)
        //提货码
        $takeGoodsCode = $this->request->get('take_goods_code');
        if (!empty($takeGoodsCode)) {
            $input['take_goods_code'] = $takeGoodsCode;
        }
        //关键词
        $keywords = $this->request->get('keywords');
        if (!empty($keywords)) {
            $input['order_search'] = $keywords;
        }
        // 如果传递的开始结束时间一致
        if (!empty($input['start_time']) && !empty($input['end_time'])
            && $input['start_time'] == $input['end_time']) {
            $input['start_time'] = strtotime(date('Y-m-d 00:00:00', $input['start_time']));
            $input['end_time']   = strtotime(date('Y-m-d 23:59:59', $input['end_time']));
        }

        // 如果选择了区域或门店
        if (!empty($input['store_ids'])) {
            $input['store_id'] = $input['store_ids'];
        }

        // 修改状态
        if ($input['order_type'] == 2 && $input['service_order_type'] == 0 && empty($input['order_status'])) {
            // 查询服务订单数据罗盘
            $input['order_status'] = [40];
        }

        // 服务人mem_uidID集合
        if (!empty($waiterMemUidArr)) {
            // 服务人员uid获取
            $input['waiter_ids'] = $this->getWorkerIds($storeId, $waiterMemUidArr);
        }

        //获取当前登录员工的信息
        $workerDetail = $this->storeWorkerDetail($input['enterprise_id'], $storeId);
        $workerDetailId = $workerDetail['id'] ?? 0;
        $workerDetailMemUid = $workerDetail['mem_uid'] ?? 0;

        // 推广订单
        if ($input['order_type'] == Constant::EXTENSION_ORDER_TYPE) {
            $input['order_status'] = self::ORDER_PROMOTER_STATUS;

            $this->saveViewTab(
                $input['enterprise_id'],
                $storeId,
                $workerDetailId,
                StoreMemberViewTab::EXTENSION_TAB_TYPE
            );
            //这些角色只查看自己范围的
            if (in_array($this->role, [
                Constant::LOGIN_ROLE_STORE_MEMBER,
                Constant::LOGIN_ROLE_ENTERPRISE,
                Constant::LOGIN_ROLE_COMMUNITY,
                Constant::LOGIN_ROLE_ENTERPRISE_ADMIN
            ])) {
                $input['store_id']     = [];
                $input['promoter_ids'] = [
                    $workerDetailId
                ];
            }

            // 查看店长的推广订单
            if (in_array($this->role, [
                Constant::LOGIN_ROLE_STORE_MANAGER,
            ])) {
                $input['store_id'] = [];
                // 如果选中推广导购则查询推广导购
                if (!empty($input['check_worker_ids'])) {
                    $input['promoter_ids'] = $input['check_worker_ids'];
                } elseif (!empty($input['check_all'])) {
                    $input['promoter_ids'] = $this->getWorkerIds($storeId, []);
                } else {
                    $input['promoter_ids'] = [$workerDetailId];
                }
            }
        }

        // 如果是我的接的服务订单
        if ($input['service_order_type'] == self::MY_ORDER) {
            // 如果选中推广导购则查询推广导购
            if (!empty($input['check_worker_ids'])) {
                $input['waiter_ids'] = $input['check_worker_ids'];
            } elseif (!empty($input['check_all'])) {
                $input['waiter_ids'] = $this->getWorkerIds($storeId, []);
            } else {
                $input['waiter_ids'] = [$workerDetailId];
            }
        }

        //合伙人信息
        $partnerIds = $this->getPartnerIds([
            'enterprise_id' => $input['enterprise_id'],
            'worker_ids' => $input['promoter_ids']
        ]);
        //按选择的身份设置查询条件
        if ($input['identity_type'] == parent::USERTYPE_PARTNER) {  //查询合伙人推广订单
            $input['promoter_ids'] = $partnerIds;
        } elseif ($input['identity_type'] == parent::USERTYPE_WORKER) {  //查询导购推广订单
            //$input['promoter_ids'] = $this->getWorkerIds($storeId, $promoterMemUidArr);
        } else {  //查询所有
            $input['promoter_ids'] = array_merge($input['promoter_ids'], $partnerIds);
        }

        // 如果是待接订单
        if ($input['service_order_type'] == self::ORDER_TO_BE_RECEIVED) {
            $this->saveViewTab(
                $input['enterprise_id'],
                $storeId,
                $workerDetailId,
                StoreMemberViewTab::SERVICE_ORDER_TO_BE_RECEIVED_TAB_TYPE
            );
            $input['order_status'] = self::ORDER_RECEIVED_ARRAY;

            $input['promoter_ids'] = [];
        }

        //检查查询条件 ，如果order_type=1是推广订单， check_worker_ids不为空，但是promoter_ids为空，则返回空
        if ($input['order_type'] == Constant::EXTENSION_ORDER_TYPE && $input['service_order_type'] == 0
            && !empty($input['check_worker_ids']) && empty($input['promoter_ids'])) {
            return success($this->getInitList($input));
        }

        // 如果是好物集企业则调用好物集专用
        if (intval($this->enterpriseId) == $haoWuJiEnterpriseId) {
            $input['mem_uid'] = $workerDetail['mem_uid'];
            $orders           = (new \CrsSdk\Work\Order())->statusList($input);
        } else {
            $orders = $this->getOrderList($input, $workerDetail);
        }
        // 操作退款的权限
        $orders['is_refund_authority'] = (int)config('app.is_refund_authority', true);
        return success($orders);
    }

    /**
     * 获取初始化列表
     *
     * @param array $input
     *
     * @return array
     */
    protected function getInitList($input = [])
    {

        $data = [
            'is_refund_authority' => (int)config('app.is_refund_authority', true),
            'page'                => intval($input['page']),
            'size'                => intval($input['size']),
            'total'               => 0,
            'list'                => []
        ];

        return $data;
    }

    /**
     * 获取订单列表
     *
     * @param $input
     * @param $workerDetail
     *
     * @return mixed|void
     */
    protected function getOrderList($input, $workerDetail)
    {

        try {
            Log::info('gw-assistant订单列表表请求参数' . json_encode($input));

            $orders = (new Order())->statusList($input);
            Log::info('gw-assistant订单列表表返回参数' . json_encode($orders));
            if (!empty($orders['list'])) {
                $userIds = array_unique(array_column($orders['list'], 'user_id'));

                $shopUserLIst = (new User())->users(
                    ['enterprise_id' => $this->enterpriseId, 'user_ids' => $userIds]
                )['data'];

                $shopUserLIst = array_column($shopUserLIst, null, 'id');

                $params = [
                    'union_ids'     => array_column($shopUserLIst, 'union_id'),
                    'enterprise_id' => $this->enterpriseId,
                    'mem_uid'       => $workerDetail['mem_uid'],
                    'qy_domain'     => $this->loginAdminEnterprise['domain']
                ];

                $list = (new Customer())->conversationDetail($params);

                $storeStaffIdArr = array_column($orders['list'], 'store_staff_id');
                $partnerInfoArr = [];
                $workerInfoArr = [];
                if (!empty($storeStaffIdArr)) {
                    //到work项目查询
                    //获取出$orders['list']中的partner_id，到worker项目查询这些合人人的信息
                    $input['partner_ids'] = $storeStaffIdArr;
                    $requestPartners = $input;
                    $requestPartners['page'] = 1;
                    $requestPartners['size'] = count($storeStaffIdArr);
                    $partnerInfoArr = (new PartnerMember())->partnerList($requestPartners);
                    $partnerInfoArr = $partnerInfoArr['list'] ?? [];
                    $partnerInfoArr = array_column($partnerInfoArr, null, 'id');

                    //查询导购的信息
                    $workerIdArr = array_column($partnerInfoArr, 'worker_id');
                    $workerInfoArr = (new StoreWorker())->getList([
                        'enterprise_id' => $input['enterprise_id'],
                        'ids' => $workerIdArr,
                        'page' => 1,
                        'size' => count($workerIdArr)
                    ]);
                    $workerInfoArr = $workerInfoArr['list'] ?? [];
                    $workerInfoArr = array_column($workerInfoArr, null, 'id');
                }

                $userList = [];
                if (!empty($list)) {
                    $userList = array_column($list, 'ec_userid', 'union_id');
                }
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
                    $value['service_order_type'] = $input['service_order_type'];
                    $value['goods_total']        = $this->getGoodsNum($value['goods_list']);
                    $value['parcel_button']      = false;
                    // 需配送显示添加包裹按钮(类型为快递,状态已付款,有服务人)
                    if ($value['delivery_type'] == 1 && $value['status'] == 60 && $value['service_staff_id'] > 0) {
                        $value['parcel_button'] = true;
                    }

                    $value['user_nick_name'] = '';
                    $value['user_avatar']    = '';
                    $userInfo                = $shopUserLIst[$value['user_id']] ?? [];
                    if (!empty($userInfo)) {
                        $value['user_nick_name'] = $userInfo['nick_name'];
                        $value['user_avatar']    = $userInfo['avatar'];
                    }

                    $tmp_union_id =  isset($userInfo['union_id']) ? $userInfo['union_id'] : '';
                    $value['ec_user_id'] = isset($userList[$tmp_union_id]) ? $userList[$tmp_union_id] : '';

                    //推广导购
                    $currentPartnerBindWorkerId = $partnerInfoArr[ $value['store_staff_id']??0 ]['worker_id'] ?? 0;
                    if (empty($value['promoter_name'])) {
                        $value['promoter_name'] = $workerInfoArr[ $currentPartnerBindWorkerId ]['name'] ?? '';
                    }
                    if (empty($value['promoter_mobile'])) {
                        $value['promoter_mobile'] = $workerInfoArr[ $currentPartnerBindWorkerId ]['mobile'] ?? '';
                    }

                    //加入合伙人电话与姓名
                    if (isset($value['store_staff_id'])) {
                        $arrKey = $value['store_staff_id'];
                        $value['partner_id'] = $partnerInfoArr[$arrKey]['id'] ?? 0;
                        $value['partner_name'] = $partnerInfoArr[$arrKey]['partner_name'] ?? '';
                        $value['partner_mobile'] = $partnerInfoArr[$arrKey]['partner_mobile'] ?? '';
                        $value['partner_mobile_sensitive'] = $partnerInfoArr[$arrKey]['mobile_sensitive'] ?? 0;
                        if ($value['partner_mobile_sensitive'] == 1 && !empty($value['partner_mobile'])) {
                            $value['partner_mobile'] = substr_replace($value['partner_mobile'], '****', 3, 4);
                        }
                    }
                }
            }
        } catch (\Exception $exception) {
            error('654009');
        }

        return $orders;
    }

    /**
     * 获取店员IDs
     *
     * @param int   $storeId
     * @param array $memUids
     *
     * @return array
     * @throws \Exception
     */
    public function getWorkerIds($storeId = 0, $memUids = [])
    {
        // 获取指定条件服务人对应数据
        $workerMemList = (new StoreWorker())->getList([
            'enterprise_id' => $this->enterpriseId,
            'store_id'      => $storeId,
            'mem_uids'      => $memUids,
            'is_page'       => false
        ]);

        return array_column($workerMemList['list'], 'id');
    }

    /**
     * 获取合伙人的ID
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2020/11/30 11:04
     * @param $input
     * @return array
     * @throws \Exception
     */
    public function getPartnerIds($input): array
    {
        //$input['partner_mobile'] = $input['account_mobile'];
        $requestData['partner_ids'] = $input['partner_ids'] ?? [];
        $requestData['worker_id'] = $input['worker_id'] ?? 0;
        $requestData['worker_ids'] = $input['worker_ids'] ?? [];
        $requestData['enterprise_id'] = $input['enterprise_id'];
        $requestData['page'] = $input['page'] ?? 1;
        $requestData['size'] = $input['size'] ?? 10000;
        $partnerInfoArr = (new PartnerMember())->partnerList($requestData);
        $partnerList = $partnerInfoArr['list'] ?? [];

        if (empty($partnerList)) {
            return [];
        }
        return array_column($partnerList, 'id');
    }
}
