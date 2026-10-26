<?php


namespace App\Services;


use App\Common\Constant;
use App\Library\Uc\Member;
use App\Models\StaffOrderDpRelation;
use App\Models\StaffOrderRelation;
use App\Models\StaffSetting;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\StaffOrder;
use CrsSdk\Shop\StaffRight;
use Illuminate\Support\Facades\DB;

class StaffOrderService extends BaseService
{
    /**
     * @param $input
     *
     * @return bool
     */
    public function createStaffOrder($input)
    {
        $enterpriseId   = $input['enterprise_id'];
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
        $orderNo        = $input['order_no'];
        $userId         = $input['user_id'];
        $qyDomain       = $enterpriseInfo['data']['domain'];
        $orderMoney     = $input['order_money'];
        //查询订单是否存在
        $order_info = StaffOrderRelation::query()
            ->where('order_no', '=', $orderNo)
            ->first();
        if (!empty($order_info)) {
            return false;
        }

        $params = [
            'enterprise_id' => $enterpriseId,
            'user_id'       => $userId
        ];

        $staffMemberData = (new StaffRight())->staffMemberByUserId($params);

        if (!empty($staffMemberData['data'])) {
            // 获取内购员工id
            $memUid = $staffMemberData['data'][0]['mem_uid'] ?? '';

            $staffIsBuy = (new StaffManagerService())->staffIsBuy($enterpriseId, $memUid, $qyDomain);
            if (empty($staffIsBuy)) {
                return false;
            }

            $create = [
                'enterprise_id' => $enterpriseId,
                'mem_uid'       => $memUid,
                'order_no'      => $orderNo,
                'order_money'   => $orderMoney,
                'is_my_order'   => $staffMemberData['is_my_order']
            ];

            $staffOrder = StaffOrderRelation::query()->create($create);

            $createDpData = [];
            // 写入部门数据
            foreach ($staffMemberData['dp_ids'] as $dp_id) {

                $createDpData[] = [
                    'enterprise_id'           => $enterpriseId,
                    'mem_uid'                 => $memUid,
                    'order_no'                => $orderNo,
                    'staff_order_relation_id' => $staffOrder->id,
                    'dp_id'                   => $dp_id,
                    'created_at'              => time()
                ];
            }

            // 如果有部门则插入
            if (!empty($createDpData)) {
                StaffOrderDpRelation::query()->insert($createDpData);
            } else {
                // 如果当前用户没有部门则写入空数据
                $createDpData = [
                    'enterprise_id'           => $enterpriseId,
                    'mem_uid'                 => $memUid,
                    'order_no'                => $orderNo,
                    'staff_order_relation_id' => $staffOrder->id,
                    'dp_id'                   => '',
                ];
                StaffOrderDpRelation::query()->create($createDpData);
            }
        }

        return true;
    }

    /**
     * 退款成功更新内购订单
     *
     * @param $input
     *
     * @return bool
     */
    public function updateStaffOrderRefund($input)
    {
        $order_no              = $input['order_no'];
        $settlement_refund_fee = $input['settlement_refund_fee']; //实际退款金额

        //查询订单是否存在
        $order_info = StaffOrderRelation::query()
            ->where('order_no', '=', $order_no)
            ->first();
        if (empty($order_info)) {
            return false;
        }

        //检查退款金额是否大与总支付金额
        $total_refund = $settlement_refund_fee + $order_info['refund_money'];

        if ($total_refund > $order_info['pay_money']) {
            $update = [
                'refund_money' => $order_info['pay_money'],
            ];
        } else {
            $update = [
                'refund_money' => $total_refund,
            ];
        }

        StaffOrderRelation::query()
            ->where('order_no', '=', $order_no)
            ->update($update);
    }

    /**
     * 支付成功更新内购订单
     *
     * @param $input
     *
     * @return bool
     */
    public function updateStaffOrderPay($input)
    {
        $order_no       = $input['order_no'];
        $pay_money      = $input['pay_money']; //实际支付金额
        $order_pay_time = $input['pay_at']; //实际支付时间

        //查询订单是否存在
        $order_info = StaffOrderRelation::query()
            ->where('order_no', '=', $order_no)
            ->first();

        if (empty($order_info)) {
            return false;
        }

        $update = [
            'pay_money'      => $pay_money,
            'order_pay_time' => $order_pay_time
        ];

        StaffOrderRelation::query()
            ->where('order_no', '=', $order_no)
            ->update($update);
    }

    /**
     * 获取查询时间段
     *
     * @param int $type 1=当天 2=最近7天 3=最近30天
     *
     * @return array
     */
    public function getTimeData($type = 0)
    {
        $startTime = $endTime = '';

        switch ($type) {
            case Constant::ORDER_TODAY_TYPE:
                $startTime = date('Y-m-d 00:00:00');
                break;
            case Constant::ORDER_LAST_SEVEN_DAYS:
                $startTime = date('Y-m-d', strtotime('-7 days'));
                break;
            case Constant::ORDER_LAST_THIRTY_DAYS:
                $startTime = date('Y-m-d', strtotime('-30 days'));
                break;
            case Constant::ORDER_YESTERDAY_DAYS:
                $startTime = date('Y-m-d', strtotime('-1 days'));
                $endTime = date('Y-m-d', time());
                break;
        }

        // 如果时间类型存在
        if (!empty($startTime)) {
            $startTime = strtotime($startTime);
            if(empty($endTime)){
                $endTime   = time();
            }else{
                $endTime   = strtotime($endTime);
            }

        }

        return [$startTime, $endTime];
    }

    /**
     * 获取订单时间搜索条件
     *
     * @param array $search
     *
     * @return array
     */
    public function orderTimeWhere($search = [])
    {
        $startTime = $endTime = 0;
        if (!empty($search['start_time'])) {
            $startTime = $search['start_time'];
        }

        if (!empty($search['end_time'])) {
            $endTime = $search['end_time'];
        }

        if (!empty($search['type'])) {

            list($startTime, $endTime) = $this->getTimeData($search['type']);
        }

        return [$startTime, $endTime];
    }

    /**
     * 统计当前员工业绩
     *
     * @param $params
     * @param $mem_uid
     * @param $enterprise_id
     *
     * @return array
     */
    public function getMemStatisticsInfo($params, $mem_uid, $enterprise_id)
    {
        $return_data = [];
        $start_time  = '';
        $end_time    = '';
        switch ($params['query_time']) {
            case 'today':
                $start_time = strtotime(date('Y-m-d 00:00:00'));
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'seven_days':
                $start_time = strtotime("-7 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'thirty_days':
                $start_time = strtotime("-30 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'all':
                break;
        }


        if (!empty($start_time)) {
            $sql = "SELECT t.mem_uid,sum(t.pay_money) as order_total_money,count(1) as order_total_count FROM staff_order_relation AS t WHERE t.order_no!='1' AND t.enterprise_id=" . $enterprise_id . " AND t.mem_uid='" . $mem_uid . "' AND t.pay_money>0 AND t.refund_money=0 AND t.deleted_at=0 and t.created_at>=" . $start_time . " AND t.created_at<=" . $end_time;
        } else {
            $sql = "SELECT t.mem_uid,sum(t.pay_money) as order_total_money,count(1) as order_total_count FROM staff_order_relation AS t WHERE t.order_no!='1' AND t.enterprise_id=" . $enterprise_id . " AND t.mem_uid='" . $mem_uid . "' AND t.pay_money>0 AND t.refund_money=0 AND t.deleted_at=0";
        }

        $res = DB::select($sql);
        if (!empty($res)) {
            return $res;
        }

        return $return_data;
    }

    /**
     * 统计员工业绩排行
     *
     * @param $params
     * @param $enterprise_id
     * @param $noOrderStaff
     *
     * @return array
     */
    public function getMemStatisticsRank($params, $enterprise_id, $noOrderStaff)
    {
        $return_data = [];
        $start_time  = '';
        $end_time    = '';
        switch ($params['query_time']) {
            case 'today':
                $start_time = strtotime(date('Y-m-d 00:00:00'));
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'seven_days':
                $start_time = strtotime("-7 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'thirty_days':
                $start_time = strtotime("-30 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'all':
                break;
        }
        $start = ($params['page'] - 1) * $params['size'];
        $end   = $params['size'];

        $rannum = $start + 1;

        if (!empty($start_time)) {
            $sql = "SELECT t.mem_uid,sum(t.pay_money) as order_total_money,count(1) as order_total_count FROM staff_order_relation AS t WHERE t.enterprise_id=" . $enterprise_id . " AND t.refund_money=0 AND t.deleted_at=0 AND t.created_at>=" . $start_time . " AND t.created_at<=" . $end_time . " GROUP BY t.mem_uid ORDER BY order_total_money DESC,t.id ASC limit " . $start . "," . $end;
        } else {
            $sql = "SELECT t.mem_uid,sum(t.pay_money) as order_total_money,count(1) as order_total_count FROM staff_order_relation AS t WHERE t.enterprise_id=" . $enterprise_id . " AND t.refund_money=0 AND t.deleted_at=0 GROUP BY t.mem_uid ORDER BY order_total_money DESC,t.id ASC limit " . $start . "," . $end;
        }
        $res = DB::select($sql);

        //查询
        if (!empty($res)) {
            foreach ($res as $item) {
                if (!empty($noOrderStaff[$item->mem_uid])) {
                    $item->order_total_count = $item->order_total_count - $noOrderStaff[$item->mem_uid];
                }
                $item->rownum = $rannum;
                $rannum++;
            }
            return $res;
        }
        return $return_data;
    }

    /**
     * 获取排行总条数
     *
     * @param $params
     * @param $enterprise_id
     *
     * @return int
     */
    public function getRankTotal($params, $enterprise_id)
    {
        $start_time = '';
        $end_time   = '';
        switch ($params['query_time']) {
            case 'today':
                $start_time = strtotime(date('Y-m-d 00:00:00'));
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'seven_days':
                $start_time = strtotime("-7 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'thirty_days':
                $start_time = strtotime("-30 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'all':
                break;
        }

        $total = StaffOrderRelation::query()
            ->where('enterprise_id', $enterprise_id)
            ->where('refund_money', 0)
            ->when(!empty($start_time), function ($query) use ($start_time, $end_time) {
                return $query->where('created_at', '>=', $start_time)->where('created_at', '<=', $end_time);
            })
            ->groupBy('mem_uid')
            ->get(['mem_uid']);
        return count($total);
    }

    /**
     * 统计当前员工业绩排行
     *
     * @param $params
     * @param $enterprise_id
     * @param $mem_uid
     * @param $noOrderStaff
     *
     * @return array
     */
    public function getCurrentMemRank($params, $enterprise_id, $mem_uid, $noOrderStaff)
    {
        $rannum      = 1;
        $return_data = [];
        $start_time  = '';
        $end_time    = '';
        switch ($params['query_time']) {
            case 'today':
                $start_time = strtotime(date('Y-m-d 00:00:00'));
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'seven_days':
                $start_time = strtotime("-7 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'thirty_days':
                $start_time = strtotime("-30 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'all':
                break;
        }

        if (!empty($start_time)) {
            $sql = "SELECT t.mem_uid,sum(t.pay_money) as order_total_money,count(1) as order_total_count FROM staff_order_relation AS t WHERE t.enterprise_id=" . $enterprise_id . " AND t.refund_money=0 AND t.deleted_at=0 AND t.created_at>=" . $start_time . " AND t.created_at<=" . $end_time . " GROUP BY t.mem_uid ORDER BY order_total_money DESC,t.id ASC";
        } else {
            $sql = "SELECT t.mem_uid,sum(t.pay_money) as order_total_money,count(1) as order_total_count FROM  staff_order_relation AS t WHERE t.enterprise_id=" . $enterprise_id . " AND t.refund_money=0 AND t.deleted_at=0 GROUP BY t.mem_uid ORDER BY order_total_money DESC,t.id ASC";
        }
        $res  = DB::select($sql);
        $temp = [];
        if (!empty($res)) {
            if (!empty($res)) {
                foreach ($res as $item) {
                    if (!empty($noOrderStaff[$item->mem_uid])) {
                        $item->order_total_count = $item->order_total_count - $noOrderStaff[$item->mem_uid];
                    }
                    $item->rownum = $rannum;
                    $rannum++;
                    if ($mem_uid == $item->mem_uid) {
                        $temp[0] = $item;
                        break;
                    }
                }
                return $temp;
            }
        }
        return $return_data;
    }

    /**
     * 检查员工是否有业绩订单，如果没有插入一条假数据
     *
     * @param $memUid
     * @param $enterpriseId
     * @param $qyDomain
     * @param $appIdentifier
     *
     * @return bool
     * @throws \Exception
     */
    public function checkStaffOrderInfo($memUid, $enterpriseId, $qyDomain, $appIdentifier)
    {

        $start_time = strtotime(date('Y-m-d 00:00:00'));
        $end_time   = strtotime(date('Y-m-d 23:59:59'));
        $check_res  = StaffOrderRelation::query()
            ->where('mem_uid', $memUid)
            ->where('enterprise_id', $enterpriseId)
            ->where('created_at', '>=', $start_time)
            ->where('created_at', '<=', $end_time)
            ->first();

        if (empty($check_res)) {
            $params   = ['memUid' => $memUid];
            $userInfo = (new Member())->detail($qyDomain, $appIdentifier, $params);
            // 查询用户所在部门集合
            $dpIds = [];
            if (!empty($userInfo['dpName'])) {
                $dpIds = array_column($userInfo['dpName'], 'dpId') ?? [];
            }
            $create_data = [
                'enterprise_id' => $enterpriseId,
                'mem_uid'       => $memUid,
                'order_no'      => '1',
                'is_my_order'   => 2 // 不是自己
            ];

            $staffOrder = StaffOrderRelation::query()->create($create_data);

            $createDpData = [];
            // 写入部门数据
            foreach ($dpIds as $dp_id) {

                $createDpData[] = [
                    'enterprise_id'           => $enterpriseId,
                    'mem_uid'                 => $memUid,
                    'order_no'                => '1',
                    'staff_order_relation_id' => $staffOrder->id,
                    'dp_id'                   => $dp_id,
                    'created_at'              => time()
                ];
            }

            // 如果有部门则插入
            if (!empty($createDpData)) {
                StaffOrderDpRelation::query()->insert($createDpData);
            } else {
                // 如果当前用户没有部门则写入空数据
                $createDpData = [
                    'enterprise_id'           => $enterpriseId,
                    'mem_uid'                 => $memUid,
                    'order_no'                => '1',
                    'staff_order_relation_id' => $staffOrder->id,
                    'dp_id'                   => '',
                ];
                StaffOrderDpRelation::query()->create($createDpData);
            }
        }

        return true;
    }

    /**
     * 获取包含虚假业绩员工的信息
     *
     * @param $params
     * @param $enterprise_id
     *
     * @return array
     */
    public function getNoOrderStaff($params, $enterprise_id)
    {
        $start_time = '';
        $end_time   = '';
        switch ($params['query_time']) {
            case 'today':
                $start_time = strtotime(date('Y-m-d 00:00:00'));
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'seven_days':
                $start_time = strtotime("-7 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'thirty_days':
                $start_time = strtotime("-30 day");
                $end_time   = strtotime(date('Y-m-d 23:59:59'));
                break;
            case 'all':
                break;
        }

        $check_res = StaffOrderRelation::query()
            ->where('pay_money', 0)
            ->where('enterprise_id', $enterprise_id)
            ->when(!empty($start_time), function ($query) use ($start_time, $end_time) {
                return $query->where('created_at', '>=', $start_time)->where('created_at', '<=', $end_time);
            })
            ->selectRaw("count(1) as num,mem_uid")
            ->groupBy('mem_uid')
            ->get()
            ->toArray();

        if (!empty($check_res)) {
            $check_res = array_column($check_res, 'num', 'mem_uid');
        }

        return $check_res;
    }

    /**
     * 获取用户业绩总数
     *
     * @param array $search
     * @param array $memUids
     *
     * @return array
     */
    public function getUserTotalNum($search = [], $memUids = [])
    {

        if (empty($memUids)) {
            return [];
        }

        $obj = StaffOrderRelation::query()
            ->select(DB::raw('sum(if(order_no=1,0,1)) as total_num,sum(pay_money) as total_money,mem_uid'))
            ->when(!empty($memUids), function ($query) use ($memUids) {
                return $query->where(!empty($memUids), function ($whereSql) use ($memUids) {
                    $chunkMemUids = array_chunk($memUids, 1);
                    foreach ($chunkMemUids as $key => $uid) {
                        if ($key == 0) {
                            $whereSql->whereIn('mem_uid', $uid);
                        } else {
                            $whereSql->orWhereIn('mem_uid', $uid);
                        }
                    }
                    return $whereSql;
                });
            })
            ->where('enterprise_id', $search['enterprise_id'])
            ->groupBy('mem_uid')
            ->get()->toArray();

        return array_column($obj, null, 'mem_uid') ?? [];
    }

    /**
     * 获取订单数以及订单金额
     *
     * @param array $search
     * @param array $memUids
     *
     * @return array
     */
    public function getTotalNum($search = [], $memUids = [])
    {

        if (empty($memUids)) {
            return [0, 0];
        }

        $obj = StaffOrderRelation::query()
            ->select(DB::raw('sum(if(order_no=1,0,1)) as num,sum(pay_money) as money'))
            ->when(!empty($memUids), function ($query) use ($memUids) {
                $chunkMemUids = array_chunk($memUids, 500);
                foreach ($chunkMemUids as $key => $uid) {
                    if ($key == 0) {
                        $query->whereIn('mem_uid', $uid);
                    } else {
                        $query->orWhereIn('mem_uid', $uid);
                    }
                }
                return $query;
            })
            ->when($search['settle_type'] == StaffSetting::TRADE_SETTLE, function ($query) use ($search) {
                return $query->where('order_pay_time', '>', 0)->where('refund_money', 0);
            })
            ->when($search['commission_type'] == StaffSetting::NOT_HAS_COMMISSION, function ($query) use ($search) {
                return $query->where('is_my_order', 2);
            })
            ->when(!empty($search['start_time']), function ($query) use ($search) {
                return $query->where('created_at', '>=', $search['start_time']);
            })
            ->when(!empty($search['end_time']), function ($query) use ($search) {
                return $query->where('created_at', '<', $search['end_time']);
            })
            ->where('enterprise_id', $search['enterprise_id'])
            ->get()->toArray();

        return [$obj[0]['num'] ?? 0, $obj[0]['money'] ?? 0];
    }

    /**
     * @param      $search
     *          +enterprise_id 企业ID
     *          +mem_uid 用户ID
     *          +type 类型 0=自定义时间段 1=今天 2=7天 3=30天
     *          +qy_domain 企业domain
     *          +app_identifier 企业app_identifier
     * @param int  $isAll true是否查询全部(业绩)默认查询总业绩
     *
     * @return array
     * @throws \Exception
     */
    public function getMyStaff($search, $isAll = 0)
    {
        // 获取和判断员工是不是属于允许内购得员工
        $isAuthor = (new StaffManagerService())->staffIsBuy($search['enterprise_id'], $search['mem_uid'], $search['qy_domain'], $search['app_identifier']);

        try {
            // 如果有权限则写入记录
            if ($isAuthor) {
                // 验证
                $this->checkStaffOrderInfo($search['mem_uid'], $search['enterprise_id'], $search['qy_domain'], $search['app_identifier']);
            }
        } catch (\Exception $exception) {

        }

        // 如果不是自定时间
        if ($search['type'] != 4) {
            // 获取开始时间结束时间
            list($startTime, $endTime) = $this->orderTimeWhere($search);
            $search['start_time'] = $startTime;
            $search['end_time']   = $endTime;
        }

        // 查询配置信息
        $staffSetting = StaffSetting::query()
            ->where([
                'enterprise_id' => $search['enterprise_id'],
                'type'          => StaffSetting::STAFF_TYPE,
            ])
            ->first();
        $search['commission_type'] = 0;
        $search['settle_type']     = 0;

        // 获取总数以及总订单数
        list($totalNum, $totalMoney) = $this->getTotalNum($search, [$search['mem_uid']]);

        $search['commission_type'] = $staffSetting->commission_type ?? 0;

        $money = $num = 0;
        // 是否查询业绩
        if ($isAll) {
            // 交易完成时
            if ($staffSetting->settle_type == StaffSetting::TRADE_SETTLE) {
                $search['settle_type'] = StaffSetting::TRADE_SETTLE;
                // 获取业绩数以及业绩订单数
                list($num, $money) = $this->getTotalNum($search, [$search['mem_uid']]);
            }

            // 售后维权期结束结算
            if ($staffSetting->settle_type == StaffSetting::POST_SALE_SETTLE) {

                list($num, $money) = $this->endOfCustomerServiceOrder($search);

            }
        }

        $data =
            [
                'total_num'   => intval($totalNum),
                'total_money' => intval($totalMoney),
                'num'         => intval($num) ?? 0,
                'money'       => intval($money) ?? 0,
                'start_time'  => $search['start_time'],
                'end_time'    => $search['end_time'],
                'is_author'   => $isAuthor,
                'achievement_desc' => $staffSetting->achievement_desc ?? ''
            ];

        return $data;
    }


    /**
     * 售后完成时获取订单数,订单金额
     *
     * @param $search
     *          +enterprise_id 企业ID
     *          +mem_uid 用户ID
     *          +start_time 开始时间
     *          +end_time 结束时间
     *
     * @return array
     * @throws \Exception
     */
    public function endOfCustomerServiceOrder($search)
    {

        $staffOrderList = StaffOrderRelation::query()
            ->where('mem_uid', $search['mem_uid'] ?? 0)
            ->where('order_pay_time', '>', 0)
            ->when(!empty($search['start_time']), function ($query) use ($search) {
                return $query->where('created_at', '>=', $search['start_time']);
            })
            ->when(!empty($search['end_time']), function ($query) use ($search) {
                return $query->where('created_at', '<', $search['end_time']);
            })
            ->when($search['commission_type'] == StaffSetting::NOT_HAS_COMMISSION, function ($query) use ($search) {
                return $query->where('is_my_order', 2);
            })
            ->where('refund_money', 0)
            ->where('order_no', '<>', 1)
            ->where('enterprise_id', $search['enterprise_id'])
            ->get()
            ->toArray();

        // 如果没有订单
        if (empty($staffOrderList)) {
            return [0, 0];
        }

        $orderNos  = array_column($staffOrderList, 'order_no');
        $orderList = (new StaffOrder())->listByOrderNo(['order_no' => $orderNos]);
        $orderList = array_column($orderList, null, 'order_no');

        $num = $money = 0;
        foreach ($staffOrderList as $value) {
            $orderInfo = $orderList[$value['order_no']] ?? [];
            if ($orderInfo['status'] == 100) {
                $num++;
                $money = $money + $value['pay_money'];
            }
        }

        return [$num, $money];
    }

}