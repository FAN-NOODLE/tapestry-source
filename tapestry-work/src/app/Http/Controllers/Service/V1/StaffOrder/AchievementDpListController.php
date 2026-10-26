<?php


namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StaffOrderRelation;
use App\Models\StaffSetting;
use App\Services\CommonService;
use App\Services\MemberService;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use CrsSdk\Shop\StaffOrder;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class AchievementDpListController extends BaseController
{
    const MAX_LIMIT = 500; // UC数据查询返回最大条数

    /**
     * 内购部门业绩列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $input      = $request->only(['page', 'size', 'search']);
        $searchData = $input['search'] ?? [];
        $search     = [
            'type'           => $searchData['type'] ?? 0,
            'order_type'     => $searchData['order_type'] ?? 1,
            'order_desc'     => $searchData['order_desc'] ?? 1,
            'start_time'     => $searchData['start_time'] ?? 0,
            'end_time'       => $searchData['end_time'] ?? 0,
            'department'     => array_unique($searchData['department'] ?? []),
            'enterprise_id'  => $searchData['enterprise_id'] ?? '',
            'qy_domain'      => $searchData['qy_domain'] ?? '',
            'app_identifier' => $searchData['app_identifier'] ?? '',
            'is_all'         => $searchData['is_all'] ?? 0
        ];

        if ($search['type'] != Constant::ORDER_APPOINT_DAYS) {
            list($search['start_time'], $search['end_time']) = (new StaffOrderService())->orderTimeWhere($search);
        }

        if ($search['start_time'] > time()) {
            throw new Error(50000, '开始时间不能大于当前时间');
        }
        if ($search['start_time'] > $search['end_time']) {
            throw new Error(50000, '开始时间不能大于结束时间');
        }

        $time = $search['end_time'] - $search['start_time'];
        if ($time > 3600 * 24 * 365 * 2) {
            throw new Error(50000, '时间搜索范围最大支持两年');
        }

        if (empty($search['department']) && $search['is_all'] == 0) {
            if (!isset($input['page']) || !isset($input['size'])) {
                return success([]);
            }
            return success([
                'meta' => [
                    'total' => 0,
                    'page'  => (int)$input['page'],
                    'size'  => (int)$input['size']
                ],
                'data' => []
            ]);
        }

        $staffSetting = StaffSetting::query()
            ->where([
                'enterprise_id' => $search['enterprise_id'],
                'type'          => StaffSetting::STAFF_TYPE,
            ])
            ->first();
        if (!$staffSetting) {
            return success([
                'meta' => [
                    'total' => 0,
                    'page'  => (int)$input['page'],
                    'size'  => (int)$input['size']
                ],
                'data' => []
            ]);
        }
        $staffSetting = $staffSetting->toArray();

        $departmentRecursion = $department = $deptList = [];
        if ($search['is_all'] == 0) {
            // 查询所有部门
            $deptList = MemberService::departList(
                $search['qy_domain'],
                $search['app_identifier'],
                ['pageSize' => self::MAX_LIMIT]
            );

            foreach ($search['department'] as $v) {
                $departmentRecursion[$v] = CommonService::getSubs($deptList, $v, 'dpId', 'dpParentid', 'dpId');
                array_push($departmentRecursion[$v], $v);
                $department = array_merge($department, $departmentRecursion[$v]);
            }
            $department = array_unique($department);
            $deptList   = array_column($deptList, null, 'dpId');
        }

        // 查询订单列表
        $staffOrder = StaffOrderRelation::query()
            ->select(DB::raw('staff_order_dp_relation.dp_id, staff_order_relation.id, staff_order_relation.mem_uid, staff_order_relation.order_no, staff_order_relation.pay_money, staff_order_relation.refund_money, staff_order_relation.order_pay_time'))
            ->join('staff_order_dp_relation', function ($join) {
                $join->on([
                    ['staff_order_relation.id', '=', 'staff_order_dp_relation.staff_order_relation_id']
                ]);
            })
            ->when($staffSetting['commission_type'] == StaffSetting::NOT_HAS_COMMISSION, function ($query) use ($search) {
                return $query->where('staff_order_relation.is_my_order', StaffOrderRelation::IS_MY_ORDER_NO);
            })
            ->when(!empty($department), function ($query) use ($department) {
                return $query->whereIn('staff_order_dp_relation.dp_id', $department);
            })
            ->when(!empty($search['start_time']), function ($query) use ($search) {
                return $query->where('staff_order_relation.created_at', '>', $search['start_time']);
            })
            ->when(!empty($search['end_time']), function ($query) use ($search) {
                return $query->where('staff_order_relation.created_at', '<', $search['end_time']);
            })
            ->where('staff_order_dp_relation.enterprise_id', $search['enterprise_id'])
            ->where('staff_order_relation.mem_uid', '<>', '')
            ->where('staff_order_relation.order_no', '<>', 1)
            ->where('staff_order_relation.deleted_at', 0)
            ->where('staff_order_dp_relation.deleted_at', 0)
            ->get()
            ->toArray();

        $data = [];

        // 查询全部
        if ($search['is_all'] == 1) {
            $data[]   = [
                'dp_id'       => 0,
                'dp_name'     => '全部',
                'total_num'   => 0,
                'total_money' => 0,
                'num'         => 0,
                'money'       => 0
            ];
            $orderTmp = [];
            // 佣金结算方式：1:交易完成结算 2：售后维权期结束结算
            if ($staffSetting['settle_type'] == StaffSetting::TRADE_SETTLE) {
                foreach ($staffOrder as $val) {
                    if (in_array($val['order_no'], $orderTmp)) {
                        continue;
                    }
                    array_push($orderTmp, $val['order_no']);
                    // 总订单数
                    $data[0]['total_num'] = $data[0]['total_num'] + 1;
                    // 总订单金额
                    $data[0]['total_money'] = $data[0]['total_money'] + $val['pay_money'];
                    if ($val['refund_money'] == 0 && $val['order_pay_time'] > 0) {
                        // 业绩订单数
                        $data[0]['num'] = $data[0]['num'] + 1;
                        // 业绩订单金额
                        $data[0]['money'] = $data[0]['money'] + $val['pay_money'];
                    }
                }
            } else {
                $orderNos  = array_column($staffOrder, 'order_no');
                $orderList = (new StaffOrder())->listByOrderNo(['order_no' => $orderNos]);
                $orderList = array_column($orderList, null, 'order_no');
                foreach ($staffOrder as $val) {
                    if (in_array($val['order_no'], $orderTmp)) {
                        continue;
                    }
                    array_push($orderTmp, $val['order_no']);
                    // 总订单数
                    $data[0]['total_num']   = $data[0]['total_num'] + 1; // 总订单金额
                    $data[0]['total_money'] = $data[0]['total_money'] + $val['pay_money'];
                    if (isset($orderList[$val['order_no']]) && $orderList[$val['order_no']]['status'] == 100 && $orderList[$val['order_no']]['refund_at'] == 0) {
                        // 业绩订单数
                        $data[0]['num'] = $data[0]['num'] + 1;
                        // 业绩订单金额
                        $data[0]['money'] = $data[0]['money'] + $val['pay_money'];
                    }
                }
            }
        } else {
            foreach ($departmentRecursion as $key => $depart) {
                $data[$key] = [
                    'dp_id'       => $key,
                    'dp_name'     => isset($deptList[$key]) ? $deptList[$key]['dpName'] : '',
                    'total_num'   => 0,
                    'total_money' => 0,
                    'num'         => 0,
                    'money'       => 0
                ];

                $orderTmp = [];
                // 佣金结算方式：1:交易完成结算 2：售后维权期结束结算
                if ($staffSetting['settle_type'] == StaffSetting::TRADE_SETTLE) {
                    foreach ($staffOrder as $val) {
                        if (in_array($val['dp_id'], $depart)) {
                            if (in_array($val['order_no'], $orderTmp)) {
                                continue;
                            }
                            array_push($orderTmp, $val['order_no']);
                            // 总订单数
                            $data[$key]['total_num'] = $data[$key]['total_num'] + 1;
                            // 总订单金额
                            $data[$key]['total_money'] = $data[$key]['total_money'] + $val['pay_money'];
                            if ($val['refund_money'] == 0 && $val['order_pay_time'] > 0) {
                                // 业绩订单数
                                $data[$key]['num'] = $data[$key]['num'] + 1;
                                // 业绩订单金额
                                $data[$key]['money'] = $data[$key]['money'] + $val['pay_money'];
                            }
                        }
                    }
                } else {
                    $orderNos = array_column($staffOrder, 'order_no');
                    if (!$orderNos) {
                        continue;
                    }
                    $orderList = (new StaffOrder())->listByOrderNo(['order_no' => $orderNos]);
                    $orderList = array_column($orderList, null, 'order_no');
                    foreach ($staffOrder as $val) {
                        if (in_array($val['dp_id'], $depart)) {
                            if (in_array($val['order_no'], $orderTmp)) {
                                continue;
                            }
                            array_push($orderTmp, $val['order_no']);
                            // 总订单数
                            $data[$key]['total_num']   = $data[$key]['total_num'] + 1; // 总订单金额
                            $data[$key]['total_money'] = $data[$key]['total_money'] + $val['pay_money'];
                            if (isset($orderList[$val['order_no']]) && $orderList[$val['order_no']]['status'] == 100 && $orderList[$val['order_no']]['refund_at'] == 0) {
                                // 业绩订单数
                                $data[$key]['num'] = $data[$key]['num'] + 1;
                                // 业绩订单金额
                                $data[$key]['money'] = $data[$key]['money'] + $val['pay_money'];
                            }
                        }
                    }
                }
            }
        }

        // 业绩订单额排序
        $filed = 'money';
        $sort  = SORT_DESC;
        // 业绩订单数排序
        if ($search['order_type'] == 2) {
            $filed = 'num';
        }
        // 1 升序 2 降序
        if ($search['order_desc'] == 1) {
            $sort = SORT_ASC;
        }

        $data  = array_values($data);
        $total = count($data);

        array_multisort(array_column($data, $filed), $sort, $data);

        if (!isset($input['page']) || !isset($input['size'])) {
            $orderNum = 1;
            foreach ($data as &$value) {
                $value['order_num'] = $orderNum;
                $orderNum++;
            }
            unset($value);

            return success($data);
        }

        $start = ($input['page'] - 1) * $input['size'];
        $data  = array_slice($data, $start, $input['size']);

        $orderNum = $start + 1;
        foreach ($data as &$value) {
            $value['order_num'] = $orderNum;
            $orderNum++;
        }
        unset($value);

        return success([
            'meta' => [
                'total' => $total,
                'page'  => (int)$input['page'],
                'size'  => (int)$input['size']
            ],
            'data' => $data
        ]);
    }

}