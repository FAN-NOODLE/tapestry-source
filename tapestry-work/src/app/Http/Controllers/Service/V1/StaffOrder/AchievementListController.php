<?php


namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\StaffOrderRelation;
use App\Models\StaffSetting;
use App\Services\CommonService;
use App\Services\MemberService;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use CrsSdk\Shop\StaffOrder;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class AchievementListController extends BaseController
{
    const MAX_LIMIT = 500; // UC数据查询返回最大条数

    /**
     * 内购员工业绩列表
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
            'is_detail'      => $searchData['is_detail'] ?? 2, // 默认不查询部门业绩统计详情
            'start_time'     => $searchData['start_time'] ?? 0,
            'end_time'       => $searchData['end_time'] ?? 0,
            'department'     => array_unique($searchData['department'] ?? []),
            'enterprise_id'  => $searchData['enterprise_id'] ?? '',
            'qy_domain'      => $searchData['qy_domain'] ?? '',
            'app_identifier' => $searchData['app_identifier'] ?? '',
            'mobile'         => $searchData['mobile'] ?? '',
            'username'       => $searchData['username'] ?? ''
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

        if ($search['is_detail'] == 1 && empty($search['department'])) {
            throw new Error(50000, '部门id必传');
        }

        $staffSetting = StaffSetting::query()
            ->where([
                'enterprise_id' => $search['enterprise_id'],
                'type'          => StaffSetting::STAFF_TYPE,
            ])
            ->first();
        if (!$staffSetting) {
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
        $staffSetting = $staffSetting->toArray();

        $paramsSearch = [];
        if ($search['mobile'] != '') {
            $paramsSearch['memMobile'] = $search['mobile'];
        }
        if ($search['username'] != '') {
            $paramsSearch['memUsername'] = $search['username'];
        }

        $memUidList = [];
        if (!empty($paramsSearch)) {
            $paramsSearch['memAll'] = 1;
            // 查询所有员工
            $normalMembers = MemberService::memberList($search['qy_domain'], $search['app_identifier'], $paramsSearch);
            $normalMembers = array_column($normalMembers, null, 'memUid');
            $memUidList = array_column($normalMembers, 'memUid');
            if (!$memUidList) {
                return success([
                    'meta' => [
                        'total' => 0,
                        'page'  => (int)$input['page'],
                        'size'  => (int)$input['size']
                    ],
                    'data' => []
                ]);
            }
        }

        // 查询所有部门
        $deptList = MemberService::departList(
            $search['qy_domain'],
            $search['app_identifier'],
            ['pageSize' => self::MAX_LIMIT]
        );

        $department = [];
        foreach ($search['department'] as $v) {
            $department = array_merge($department, CommonService::getSubs($deptList, $v, 'dpId', 'dpParentid', 'dpId'));
            array_push($department, $v);
        }
        $department = array_unique($department);

        // 查询订单列表
        $staffOrder = StaffOrderRelation::query()
            ->select(DB::raw('group_concat(staff_order_dp_relation.dp_id) as dp_id, staff_order_relation.id, staff_order_relation.mem_uid, staff_order_relation.order_no, staff_order_relation.pay_money, staff_order_relation.refund_money, staff_order_relation.order_pay_time'))
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
            ->when(!empty($memUidList), function ($query) use ($memUidList) {
                return $query->whereIn('staff_order_relation.mem_uid', $memUidList);
            })
            ->where('staff_order_dp_relation.enterprise_id', $search['enterprise_id'])
            ->where('staff_order_relation.mem_uid', '<>', '')
            ->where('staff_order_relation.deleted_at', 0)
            ->where('staff_order_dp_relation.deleted_at', 0)
            ->groupBy('staff_order_dp_relation.staff_order_relation_id')
            ->get()
            ->toArray();

        if (!isset($normalMembers) && $staffOrder) {
            // 调用uc sdk获取员工数据
            $ucParams = [
                'memAll'  => 1 // 查询全部人员
            ];
            // 查询所有员工
            $normalMembers = MemberService::memberList($search['qy_domain'], $search['app_identifier'], $ucParams);
            $normalMembers = array_column($normalMembers, null, 'memUid');
        }

        $deptList = array_column($deptList, null, 'dpId');

        $data = [];

        // 佣金结算方式：1:交易完成结算 2：售后维权期结束结算
        if ($staffSetting['settle_type'] == StaffSetting::TRADE_SETTLE) {
            foreach ($staffOrder as $val) {
                if (!isset($data[$val['mem_uid']])) {
                    $data[$val['mem_uid']] = [
                        'id'           => $val['id'],
                        'mem_username' => isset($normalMembers[$val['mem_uid']]) ? $normalMembers[$val['mem_uid']]['memUsername'] : '',
                        'mem_face'     => isset($normalMembers[$val['mem_uid']]) ? $normalMembers[$val['mem_uid']]['memFace'] : '',
                        'mem_mobile'     => isset($normalMembers[$val['mem_uid']]) ? $normalMembers[$val['mem_uid']]['memMobile'] : '',
                        'dp_names'     => [],
                        'num'          => 0,
                        'money'        => 0,
                        'total_num'    => 0,
                        'total_money'  => 0,
                    ];
                }
                if ($val['id'] < $data[$val['mem_uid']]['id']) {
                    $data[$val['mem_uid']]['id'] = $val['id'];
                }
                $dpIdList = explode(',', $val['dp_id']);
                foreach ($dpIdList as $dp) {
                    // 所属部门
                    $deptName                          = isset($deptList[$dp]) ? [$deptList[$dp]['dpName']] : [];
                    $data[$val['mem_uid']]['dp_names'] = array_unique(array_merge($data[$val['mem_uid']]['dp_names'], $deptName));
                }
                if ($val['order_no'] != 1) {
                    // 总订单数
                    $data[$val['mem_uid']]['total_num'] = $data[$val['mem_uid']]['total_num'] + 1;
                    // 总订单金额
                    $data[$val['mem_uid']]['total_money'] = $data[$val['mem_uid']]['total_money'] + $val['pay_money'];
                    if ($val['refund_money'] == 0 && $val['order_pay_time'] > 0) {
                        // 业绩订单数
                        $data[$val['mem_uid']]['num'] = $data[$val['mem_uid']]['num'] + 1;
                        // 业绩订单金额
                        $data[$val['mem_uid']]['money'] = $data[$val['mem_uid']]['money'] + $val['pay_money'];
                    }
                }
            }
        } else {
            $orderNos = array_unique(array_column($staffOrder, 'order_no'));
            // 订单号数组为空或者订单号数组里只有订单号为1的数据
            if (empty($orderNos) || (count($orderNos) == 1 && in_array('1', $orderNos))) {
                $orderList = [];
            } else {
                $orderList = (new StaffOrder())->listByOrderNo(['order_no' => $orderNos]);
                $orderList = array_column($orderList, null, 'order_no');
            }
            foreach ($staffOrder as $val) {
                if (!isset($data[$val['mem_uid']])) {
                    $data[$val['mem_uid']] = [
                        'id'           => $val['id'],
                        'mem_username' => isset($normalMembers[$val['mem_uid']]) ? $normalMembers[$val['mem_uid']]['memUsername'] : '',
                        'mem_face'     => isset($normalMembers[$val['mem_uid']]) ? $normalMembers[$val['mem_uid']]['memFace'] : '',
                        'mem_mobile'     => isset($normalMembers[$val['mem_uid']]) ? $normalMembers[$val['mem_uid']]['memMobile'] : '',
                        'dp_names'     => [],
                        'num'          => 0,
                        'money'        => 0,
                        'total_num'    => 0,
                        'total_money'  => 0,
                    ];
                }
                if ($val['id'] < $data[$val['mem_uid']]['id']) {
                    $data[$val['mem_uid']]['id'] = $val['id'];
                }

                $dpIdList = explode(',', $val['dp_id']);
                foreach ($dpIdList as $dp) {
                    // 所属部门
                    $deptName                          = isset($deptList[$dp]) ? [$deptList[$dp]['dpName']] : [];
                    $data[$val['mem_uid']]['dp_names'] = array_unique(array_merge($data[$val['mem_uid']]['dp_names'], $deptName));
                }
                if ($val['order_no'] != 1) {
                    // 总订单数
                    $data[$val['mem_uid']]['total_num'] = $data[$val['mem_uid']]['total_num'] + 1;
                    // 总订单金额
                    $data[$val['mem_uid']]['total_money'] = $data[$val['mem_uid']]['total_money'] + $val['pay_money'];
                    if (isset($orderList[$val['order_no']]) && $orderList[$val['order_no']]['status'] == 100 && $orderList[$val['order_no']]['refund_at'] == 0) {
                        // 业绩订单数
                        $data[$val['mem_uid']]['num'] = isset($data[$val['mem_uid']]['num']) ? $data[$val['mem_uid']]['num'] + 1 : 0;
                        // 业绩订单金额
                        $data[$val['mem_uid']]['money'] = isset($data[$val['mem_uid']]['money']) ? $data[$val['mem_uid']]['money'] + $val['pay_money'] : 0;
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

        $data = array_values($data);

        array_multisort(array_column($data, $filed), SORT_NUMERIC, $sort, array_column($data, 'id'), SORT_NUMERIC, SORT_ASC, $data);

        if (!isset($input['page']) || !isset($input['size'])) {
            $orderNum = 1;
            foreach ($data as $key => &$value) {
                unset($data[$key]['id']);
                $value['order_num'] = $orderNum;
                $orderNum++;
            }
            unset($value);

            return success($data);
        }

        $detail = $search['is_detail'] == 1 ? [
            'num'         => array_sum(array_column($data, 'num')),
            'money'       => array_sum(array_column($data, 'money')),
            'total_num'   => array_sum(array_column($data, 'total_num')),
            'total_money' => array_sum(array_column($data, 'total_money')),
            'dp_name'     => isset($deptList[$search['department'][0]]) ? $deptList[$search['department'][0]]['dpName'] : ''
        ] : [];

        $total = count($data);
        $start = ($input['page'] - 1) * $input['size'];
        $data  = array_slice($data, $start, $input['size']);

        $orderNum = $start + 1;
        foreach ($data as $key => &$value) {
            unset($data[$key]['id']);
            $value['order_num'] = $orderNum;
            $orderNum++;
        }
        unset($value);

        return success([
            'meta'   => [
                'total' => $total,
                'page'  => (int)$input['page'],
                'size'  => (int)$input['size']
            ],
            'data'   => $data,
            'detail' => $detail
        ]);
    }
}