<?php

namespace App\Services;

use App\Common\Constant;
use App\Models\CrsOrder;
use App\Models\MemberExtend;
use CrsSdk\Shop\Order;
use App\Library\Uc\Member;

class OrderService
{

    /**
     * @param $input
     *
     * @return array
     * @throws \Exception
     */
    public function orderRanking($input)
    {
        $type       = $input['type'] ?? Constant::ORDER_SORT;
        $sort       = $input['sort_type'] ?? 'desc';
        $dataType   = $input['date_type'] ?? Constant::DAY;

        // 查询开始和结束日期
        switch ($dataType) {
            case 1:
                // 今日
                $start = date('Y.m.d');
                break;
            case 2:
                // 本周一
                $start = date('Y.m.d',strtotime('this week Monday', time()));
                break;
            default:
                // 当前月1号
                $start = date('Y.m.1');
                break;
        }

        // 获取当前登录用户的详情信息
        $ucResult = (new Member())->detail($input['qyDomain'], $input['appIdentifier'], ['memUid' => $input['memUid']]);

        if (empty($ucResult)) {
            error(600103);
        }

        $dpId = '';
        if (isset($ucResult['dpName'])) {
            // 取出所属部门id
            $dpId = $ucResult['dpName'][0]['dpId'];
        }

        if (empty($dpId)) {
            // 部门信息不存在
            error(600102);
        }

        $where = [
            'dpId'      => $dpId,
            'memAll'    => Constant::MEM_ALL_UC
        ];

        // 获取人员列表
        $memList = $this->listAllMember($input, $where);

        if (empty($memList)){
            [
                'list'          => [],
                'order_total'   => 0,
                'price_total'   => 0,
                'start_time'    =>  $start,
                'end_time'      =>  date('Y.m.d')
            ];
        }


        // 查询sdk所需条件
        $sdkWhere = [
            'enterprise_id' => $input['enterpriseId'],
            'date_type'     => $dataType,
            'sort_type'     => $sort,
            'mem_uid'       => array_keys($memList)
        ];

        // 调取排行SDK
        $list = (new Order())->managerRanking($sdkWhere);

        $lists = [];

        if (!empty($list['data'])){
            foreach ($list['data'] as $value){
                $lists[$value['mem_uid']] = [
                    'order_price'   => $value['order_price'],
                    'order_num'     => $value['order_num']
                ];
            }
        }

        // 组装数据
        foreach ($memList as $key=>$value){
            // if ($value['memStatus'] != Constant::MEM_DEALETE_UC){
                $memLists[$key] = [
                    'mem_uid'       => $key,
                    'mem_face'      => $value['memFace'],
                    'mem_username'  => $value['memUsername'],
                    // 人员是否离职：0=否；1=是
                    'mem_status'    => $value['memStatus'] != Constant::MEM_DEALETE_UC ? 0 : 1,
                    'order_price'   => $lists[$key]['order_price'] ?? 0,
                    'order_num'     => $lists[$key]['order_num'] ?? 0
                ];
            // }
        }

        // 排序显示
        switch ($type) {
            case 1:
                // 销售额
                $file = array_column($memLists,'order_price');
                break;
            default:
                // 订单数
                $file = array_column($memLists,'order_num');
                break;
        }

        if ($sort == 'asc'){
            array_multisort($file,SORT_ASC,$memLists);
        }else{
            array_multisort($file,SORT_DESC,$memLists);
        }

        // 上一条数据的销售额,默认为数组第一条
        $prevOrderPrice = 0;
        // 上一条数据的订单数,默认为数组第一条
        $prevOrderNum = 0;
        $i = 0;
        foreach ($memLists as $key => &$m){
            // 非第一个数据且按销售额排序
            if ($type == 1) {
                // 当前导购订单额不等于上一个导购订单额，排名自增
                if ($m['order_price'] != $prevOrderPrice) {
                    $prevOrderPrice = $m['order_price'];
                    $i++;
                } 
            }

            // 非第一个数据且按订单数排序
            if ($type == 2) {
                // 当前导购订单数不等于上一个导购订单数，排名自增
                if ($m['order_num'] != $prevOrderNum) {
                    $prevOrderNum = $m['order_num'];
                    $i++;
                }
            }

            // 如果排名没有自增则默认为1
            if ($i == 0) {
                $i = $i + 1;
            }

            $m['index_num'] = $i;
        }

        return [
            'list'          => $memLists,
            'order_total'   => array_sum(array_column($lists,'order_num')),
            'price_total'   => array_sum(array_column($lists,'order_price')),
            'start_time'    =>  $start,
            'end_time'      =>  date('Y.m.d')
        ];
    }


    /**
     * 根据条件获取人员列表
     *
     * @param $input
     * @param $condition
     *
     * @return array
     * @throws \Exception
     */
    public function listAllMember($input, $condition = [])
    {

        $page     = 1;
        $page_max = 0;
        $list     = [];

        // 获取所有成员
        do {
            $condition['pageNum']  = $page;
            $condition['pageSize'] = Constant::UC_MEMBER_NUMBER;

            $result = (new Member())->ignorePermissionsList($input['qyDomain'], $input['appIdentifier'], $condition);

            if (isset($result['list']) && !empty($result['list'])) {
                $list = array_merge($list, $result['list']);
            }

            // 计算总页数
            if ($page_max === 0) {
                $page_max = ceil($result['total'] / Constant::UC_MEMBER_NUMBER);
            }

            $page++;
        } while ($page <= $page_max);

        if (empty($list)){
            return [];
        }

        $avatarLists = [];
        $memUid = array_column($list, 'memUid');

        // 查询会员头像
        $avatarList = MemberExtend::select('mem_uid','avatar')
            ->where(['enterprise_id' => $input['enterpriseId']])
            ->whereIn('mem_uid',$memUid)
            ->get()
            ->toArray();

        if (!empty($avatarList)){
            foreach ($avatarList as $value){
                $avatarLists[$value['mem_uid']] = $value['avatar'];
            }
        }

        foreach ($list as $value){
            $lists[$value['memUid']] = [
                'memFace'       => isset($avatarLists[$value['memUid']])?$avatarLists[$value['memUid']]:$value['memFace'],
                'memUsername'   => $value['memUsername'],
                'memStatus'     => $value['memStatus']
            ];
        }

        return $lists;
    }

    /**
     * getOrderDetail
     * @desc 获取订单详情
     * @param $orderNo
     *
     * @return array
     * @date 2020-09-15
     */
    public static function getOrderDetail($orderNo)
    {
        $orderInfo = CrsOrder::query()
            ->where('order_no', $orderNo)
            ->first();
        $orderInfo->goods;
        $orderDetail = $orderInfo->toArray();
        $orderDetail['goods_data'] = $orderDetail['goods'];
        unset($orderDetail['goods']);
        $goodsTotalFee = 0;
        foreach ($orderDetail['goods_data'] as &$good) {
            $good = [
                'spu_name' => $good['spu_name'],
                'sku_attribute' => array_values(array_column(json_decode($good['sku_attribute'] ?? '[]', true), 'attribute_item_value')),
                'original_price' => $good['original_price'],
                'sell_price' => $good['sell_price'],
                'quantity' => $good['quantity'],
                'sku_picture' => $good['sku_picture'],
                'sub_price' => $good['sell_price'],
            ];
            $goodsTotalFee += $good['sell_price'] * $good['quantity'];
        }
        $orderDetail['pay_data'] = [
            'third_payment_no' => $orderDetail['third_payment_no'],
            'payment_channel' => $orderDetail['payment_channel']
        ];
        $orderDetail['order_logistic'] = [
            'logistic_identity' => $orderDetail['logistic_identity'],
            'logistic_no' => $orderDetail['logistic_no']
        ];
        $orderDetail['goods_total_fee'] = $goodsTotalFee;
        $orderDetail['created_at'] = $orderDetail['order_at'];
        $orderDetail['wx_coupon_fee'] = $orderDetail['coupon_fee'];
        $orderDetail['delivery_address'] = $orderDetail['province'] . $orderDetail['city'] . $orderDetail['area'] . $orderDetail['delivery_address'];
        return $orderDetail;
    }

}
