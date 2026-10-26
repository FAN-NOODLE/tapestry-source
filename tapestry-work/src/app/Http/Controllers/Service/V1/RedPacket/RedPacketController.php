<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2020/2/7                   
* +------------------------------------+
* | Time   : 10:07                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\Service\V1\RedPacket;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\CashCoupon;
use App\Models\CashCouponMember;
use App\Models\CashCouponRecord;
use App\Models\EnterpriseExtend;
use App\Models\MemberUserRelation;
use App\Services\MemberService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;

class RedPacketController extends BaseController
{
    /**
     * 红包列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function list(Request $request)
    {
        $params = $request->only(['enterprise_id', 'is_deleted', 'start_time', 'end_time', 'cash_coupon_name', 'size', 'page']);

        $onOff = EnterpriseExtend::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->value('cash_coupon_status');
        if ($onOff != EnterpriseExtend::OPEN) {
            return success([
                'cash_coupon_status' => EnterpriseExtend::CLOSE,
                'data'               => [],
                'meta'               => [
                    'total' => 0,
                    'size'  => $params['size'],
                    'page'  => $params['page']
                ]
            ]);
        }

        // 红包列表
        $list = CashCoupon::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('is_deleted', $params['is_deleted'])
            ->when(!empty($params['start_time']), function ($query) use ($params) {
                $query->where('created_at', '>', $params['start_time']);
            })
            ->when(!empty($params['end_time']), function ($query) use ($params) {
                $query->where('created_at', '<', $params['end_time']);
            })
            ->when(isset($params['cash_coupon_name']) && $params['cash_coupon_name'] != '', function ($query) use ($params) {
                $query->where('cash_coupon_name', 'like', '%' . $params['cash_coupon_name'] . '%');
            })
            ->orderByDesc('id')
            ->paginate($params['size'] ?? 10)
            ->toArray();

        // 获取红包发送数据
        if ($list['data']) {
            $cashCouponIds = array_column($list['data'], 'id');
            $sendData      = CashCouponRecord::query()
                ->whereIn('cash_coupon_id', $cashCouponIds)
                ->get()
                ->toArray();
            foreach ($list['data'] as &$value) {
                $value['sent_quantity']     = 0;
                $value['received_quantity'] = 0;
                foreach ($sendData as $val) {
                    if ($val['cash_coupon_id'] == $value['id']) {
                        switch ($val['status']) {
                            case CashCouponRecord::IS_SEND :
                                $value['sent_quantity']++;
                                break;
                            case CashCouponRecord::IS_RECEIVED :
                                $value['sent_quantity']++;
                                $value['received_quantity']++;
                                break;
                            default :
                                break;
                        }
                    }
                }
                $value['surplus_quantity'] = $value['quantity'] - $value['sent_quantity'];
            }
            unset($value);
        }

        return success([
            'cash_coupon_status' => EnterpriseExtend::OPEN,
            'data'               => $list['data'],
            'meta'               => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }

    /**
     * 红包对账列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function reconciliationList(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'id', 'mem_uid', 'start_time', 'end_time', 'status', 'size', 'page']);

        self::checkSwitchStatus($params['enterprise_id']);

        $detail = CashCoupon::query()
            ->where('id', $params['id'])
            ->first();
        if (!$detail) {
            throw new Error(20021001, '详情不存在');
        }

        $detail->sent_quantity     = 0;
        $detail->received_quantity = 0;
        $sendData                  = CashCouponRecord::query()
            ->where('cash_coupon_id', $detail->id)
            ->get()
            ->toArray();
        foreach ($sendData as $val) {
            switch ($val['status']) {
                case CashCouponRecord::IS_SEND :
                    $detail->sent_quantity++;
                    break;
                case CashCouponRecord::IS_RECEIVED :
                    $detail->sent_quantity++;
                    $detail->received_quantity++;
                    break;
                default :
                    break;
            }
        }
        $detail->surplus_quantity = $detail->quantity - $detail->sent_quantity;

        $list = CashCouponRecord::query()
            ->where('cash_coupon_id', $detail['id'])
            ->when(!empty($params['start_time']), function ($query) use ($params) {
                $query->where('send_time', '>', $params['start_time']);
            })
            ->when(!empty($params['end_time']), function ($query) use ($params) {
                $query->where('send_time', '<', $params['end_time']);
            })
            ->when(!empty($params['mem_uid']), function ($query) use ($params) {
                $query->where('mem_uid', $params['mem_uid']);
            })
            ->when(!empty($params['status']), function ($query) use ($params) {
                $query->where('status', $params['status']);
            })
            ->orderByDesc('id')
            ->paginate($params['size'] ?? 10)
            ->toArray();
        if ($list['data']) {
            $memUids   = array_unique(array_column($list['data'], 'mem_uid'));
            $ecUserIds = array_unique(array_column($list['data'], 'ec_userid'));
            // 调用uc sdk获取员工数据
            $ucParams      = [
                'pageNum'  => 1,
                'pageSize' => count($memUids),
                'memUids'  => $memUids
            ];
            $normalMembers = (new Member())->ignorePermissionsList($params['domain'], $params['identifier'], $ucParams);
            $normalMembers = array_column($normalMembers['list'], null, 'memUid');
            if (count($normalMembers) != count($memUids)) {
                // 查询离职员工
                $ucParams['eqMemStatus'] = 3;
                $deleteMembers           = (new Member())->ignorePermissionsList($params['domain'], $params['identifier'], $ucParams);
                $deleteMembers           = array_column($deleteMembers['list'], null, 'memUid');
                $normalMembers           = array_merge($normalMembers, $deleteMembers);
            }

            $customer = MemberUserRelation::query()
                ->whereIn('ec_userid', $ecUserIds)
                ->get()
                ->toArray();
            $customer = array_column($customer, null, 'ec_userid');
            foreach ($list['data'] as &$value) {
                $value['send_member']            = isset($normalMembers[$value['mem_uid']]) ? $normalMembers[$value['mem_uid']]['memUsername'] : '';
                $value['received_member']        = isset($customer[$value['ec_userid']]) ? $customer[$value['ec_userid']]['ec_name'] : '';
                $value['received_member_avatar'] = isset($customer[$value['ec_userid']]) ? $customer[$value['ec_userid']]['ec_avatar'] : '';
            }
            unset($value);
        }

        return success([
            'detail' => $detail,
            'data'   => $list['data'],
            'meta'   => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }

    /**
     * 发送人列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function members(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain'        => 'required',
            'identifier'    => 'required',
            'extend'        => 'required|array',
            'extend.id'     => 'required|integer|min:1'
        ];

        Validator::make($params, $rules)->validate();

        $cashCouponId = $params['extend']['id'];
        $enterpriseId = $params['enterprise_id'];

        self::checkSwitchStatus($enterpriseId);

        $members = CashCouponMember::query()
            ->where([
                'enterprise_id'  => $enterpriseId,
                'cash_coupon_id' => $cashCouponId
            ])->get([
                'object',
                'is_leave',
            ])->toArray();

        $leaveMemberUIds = [];
        $onJobMemberUIds = [];
        foreach ($members as $member) {
            if (CashCouponMember::MEMBER_IS_NORMAL == $member['is_leave']) {
                $onJobMemberUIds[] = $member['object'];
            } else {
                $leaveMemberUIds[] = $member['object'];
            }
        }

        $onJobMemberList = [];
        if ($onJobMemberUIds) {
            // 查询在职员工
            $onJobMemberList = MemberService::memberList(
                $params['domain'],
                $params['identifier'],
                [
                    'memUids' => $onJobMemberUIds
                ]
            );

            $onJobMemberList = array_column($onJobMemberList, null, 'memUid');
        }

        $leaveMemberList = [];
        if ($leaveMemberUIds) {
            // 查询所有离职员工
            $leaveMemberList = MemberService::memberList(
                $params['domain'],
                $params['identifier'],
                [
                    'memUids'     => $leaveMemberUIds,
                    'eqMemStatus' => MemberService::MEMBER_DELETE_REMARK
                ]
            );

            $leaveMemberList = array_column($leaveMemberList, null, 'memUid');
        }

        $memberList = [];
        foreach ($members as $member) {
            $memUid = $member['object'];
            if (CashCouponMember::MEMBER_IS_NORMAL == $member['is_leave']) {
                $memberList[] = [
                    'uid'  => $memUid,
                    'name' => $onJobMemberList[$memUid]['memUsername'] ?? '',
                ];
            } else {
                $memberList[] = [
                    'uid'  => $memUid,
                    'name' => $leaveMemberList[$memUid]['memUsername'] ?? '',
                ];
            }
        }

        return success([
            'data' => $memberList
        ]);
    }


    /**
     * 上下架
     *
     * @param Request $request
     *
     * @return array
     */
    public function sale(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain'        => 'required',
            'identifier'    => 'required',
            'extend'        => 'required|array',
            'extend.id'     => 'required|integer|min:1',
            'extend.status' => 'required|integer|in:2,3'
        ];

        Validator::make($params, $rules)->validate();

        $cashCouponId = $params['extend']['id'];
        $enterpriseId = $params['enterprise_id'];

        self::checkSwitchStatus($enterpriseId);

        $cashCoupon = CashCoupon::query()->where(['is_deleted' => CashCoupon::NOT_DELETED])->find($cashCouponId);

        if (!$cashCoupon) {
            throw new Error(10000, '红包不存在');
        }

        CashCoupon::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'id'            => $cashCouponId,
            ])
            ->update([
                'is_sale' => $params['extend']['status']
            ]);

        return success(['data' => []]);
    }

    /**
     * 创建红包
     *
     * @param Request $request
     *
     * @return array]
     */
    public function create(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id'           => 'required|integer|min:1',
            'domain'                  => 'required',
            'identifier'              => 'required',
            'extend.cash_coupon_name' => 'required',
            'extend.amount'           => 'required|integer|min:1',
            'extend.quantity'         => 'required|integer|min:1',
            'extend.send_name'        => 'required',
            'extend.wishing'          => 'required',
            'extend.remark'           => 'required',
            'extend.members'          => 'required|array',
            'extend.members.*.value'  => 'required',
            'extend.members.*.label'  => 'required',
        ];

        Validator::make($params, $rules)->validate();

        $enterpriseId = $params['enterprise_id'];
        $memberList   = $params['extend']['members'];

        self::checkSwitchStatus($enterpriseId);

        $params['extend']['enterprise_id'] = $enterpriseId;
        $params['extend']['number']        = self::getNumberNo($enterpriseId);

        $insert = [];
        foreach ($memberList as $member) {
            $insert[] = [
                'enterprise_id' => $enterpriseId,
                'object'        => $member['value'],
                'created_at'    => time(),
                'updated_at'    => time(),
            ];
        }

        try {
            DB::beginTransaction();
            // 红包数据写入
            $create = CashCoupon::query()->create($params['extend']);

            // 外键赋值
            foreach ($insert as &$val) {
                $val['cash_coupon_id'] = $create['id'];
            }

            // 发送人写入
            CashCouponMember::query()->insert($insert);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error($e->getCode() ?: 10000, $e->getMessage() ?: '红包创建异常');
        }

        $create['members'] = $params['extend']['members'];

        return success([
            'data' => $create
        ]);
    }

    /**
     * 红包详情
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain'        => 'required',
            'identifier'    => 'required',
            'extend'        => 'required|array',
            'extend.id'     => 'required|integer|min:1'
        ];

        Validator::make($params, $rules)->validate();

        $cashCouponId = $params['extend']['id'];
        $enterpriseId = $params['enterprise_id'];

        self::checkSwitchStatus($enterpriseId);

        // 刪除的也可查看
        $cashCoupon = CashCoupon::query()->find($cashCouponId);

        if (!$cashCoupon) {
            throw new Error(10000, '红包不存在');
        }

        // 红包有效发送人查询
        $members = CashCouponMember::query()
            ->where([
                'enterprise_id'  => $enterpriseId,
                'cash_coupon_id' => $cashCouponId
            ])->get([
                'object',
                'is_leave',
            ])->toArray();

        $leaveMemberUIds = [];
        $onJobMemberUIds = [];
        foreach ($members as $member) {
            if (CashCouponMember::MEMBER_IS_NORMAL == $member['is_leave']) {
                $onJobMemberUIds[] = $member['object'];
            } else {
                $leaveMemberUIds[] = $member['object'];
            }
        }

        // 在职员工查询
        $members = [];
        if ($onJobMemberUIds) {
            $onJobMemberList = MemberService::memberList(
                $params['domain'],
                $params['identifier'],
                [
                    'memUids' => $onJobMemberUIds
                ]
            );

            // 在职员工集合
            $onJobMemberList = array_column($onJobMemberList, null, 'memUid');

            foreach ($onJobMemberUIds as $memUid) {
                $members[] = [
                    'value'    => $memUid,
                    'label'    => $onJobMemberList[$memUid]['memUsername'] ?? '',
                    'is_leave' => CashCouponMember::MEMBER_IS_NORMAL
                ];
            }
        }

        // 离职员工查询
        if ($leaveMemberUIds) {
            $leaveMemberList = MemberService::memberList(
                $params['domain'],
                $params['identifier'],
                [
                    'memUids'     => $leaveMemberUIds,
                    'eqMemStatus' => MemberService::MEMBER_DELETE_REMARK
                ]
            );

            // 离职员工集合
            $leaveMemberList = array_column($leaveMemberList, null, 'memUid');

            foreach ($leaveMemberUIds as $memUid) {
                $members[] = [
                    'value'    => $memUid,
                    'label'    => $leaveMemberList[$memUid]['memUsername'] ?? '',
                    'is_leave' => CashCouponMember::MEMBER_IS_LEAVE
                ];
            }
        }

        $cashCoupon['members'] = $members;

        return [
            'data' => $cashCoupon
        ];
    }

    /**
     * 更新红包
     *
     * @param Request $request
     *
     * @return array
     */
    public function update(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id'           => 'required|integer|min:1',
            'domain'                  => 'required',
            'identifier'              => 'required',
            'extend'                  => 'required|array',
            'extend.id'               => 'required|integer|min:1',
            'extend.cash_coupon_name' => 'required',
            'extend.amount'           => 'required|integer|min:1',
            'extend.quantity'         => 'required|integer|min:1',
            'extend.send_name'        => 'required',
            'extend.wishing'          => 'required',
            'extend.remark'           => 'required',
            'extend.members'          => 'required|array',
            'extend.members.*.value'  => 'required',
            'extend.members.*.label'  => 'required',
        ];

        Validator::make($params, $rules)->validate();

        $cashCouponId = $params['extend']['id'];
        $enterpriseId = $params['enterprise_id'];

        self::checkSwitchStatus($enterpriseId);

        $cashCoupon = CashCoupon::query()->where(['is_deleted' => CashCoupon::NOT_DELETED])->find($cashCouponId);

        if (!$cashCoupon) {
            throw new Error(10000, '红包不存在');
        }

        if (in_array($cashCoupon['is_sale'], [CashCoupon::IS_ON_SALE, CashCoupon::IS_OFF_SALE])) {
            throw new Error(10000, '只有未上架红包才能编辑');
        }

        // 查询本地记录
        $localMemberUIds = CashCouponMember::query()
            ->where([
                'enterprise_id'  => $enterpriseId,
                'cash_coupon_id' => $cashCouponId
            ])
            ->pluck('object', 'id')
            ->toArray();

        $insert = [];
        $idsArr = [];
        foreach ($params['extend']['members'] as $member) {
            $memberUid = $member['value'];
            if (in_array($memberUid, $localMemberUIds)) {
                $id = array_search($memberUid, $localMemberUIds);
                // 更新数据
                $idsArr[] = $id;
            } else {
                $insert[] = [
                    'enterprise_id'  => $enterpriseId,
                    'cash_coupon_id' => $cashCouponId,
                    'object'         => $memberUid,
                    'created_at'     => time(),
                    'updated_at'     => time(),
                ];
            }
        }
        try {
            DB::beginTransaction();

            // 红包数据赋值
            $cashCoupon['cash_coupon_name'] = $params['extend']['cash_coupon_name'];
            $cashCoupon['amount']           = $params['extend']['amount'];
            $cashCoupon['quantity']         = $params['extend']['quantity'];
            $cashCoupon['send_name']        = $params['extend']['send_name'];
            $cashCoupon['wishing']          = $params['extend']['wishing'];
            $cashCoupon['remark']           = $params['extend']['remark'];

            // 更新红包数据
            $cashCoupon->save();

            // 删除不需要的
            CashCouponMember::query()
                ->where([
                    'enterprise_id'  => $enterpriseId,
                    'cash_coupon_id' => $cashCouponId,
                ])
                ->whereNotIn('id', $idsArr)
                ->delete();

            // 插入新增的
            $insert && CashCouponMember::query()->insert($insert);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error(50000, '红包更新异常');
        }

        $cashCoupon['members'] = $params['extend']['members'];

        return [
            'data' => $cashCoupon
        ];
    }

    /**
     * 删除红包
     *
     * @param Request $request
     *
     * @return array
     */
    public function delete(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain'        => 'required',
            'extend'        => 'required|array',
            'extend.id'     => 'required|integer|min:1'
        ];

        Validator::make($params, $rules)->validate();

        $cashCouponId = $params['extend']['id'];
        $enterpriseId = $params['enterprise_id'];

        self::checkSwitchStatus($enterpriseId);

        $cashCoupon = CashCoupon::query()->where(['is_deleted' => CashCoupon::NOT_DELETED])->find($cashCouponId);

        if (!$cashCoupon) {
            throw new Error(10000, '红包不存在');
        }

        // 删除红包
        CashCoupon::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'id'            => $cashCouponId,
            ])
            ->update([
                'is_deleted' => CashCoupon::DELETED
            ]);

        return success(['data' => []]);
    }

    /**
     * 获取编号
     *
     * @param $enterpriseId
     *
     * @return string
     */
    public function getNumberNo($enterpriseId)
    {
        $number = CashCoupon::query()
            ->where([
                'enterprise_id' => $enterpriseId
            ])
            ->orderByDesc('number')
            ->withTrashed()
            ->value('number');

        $number = intval(substr($number, 1));

        $number = 'H' . sprintf("%04d", ++$number);

        return $number;
    }

    /**
     * 检查红包开关
     *
     * @param $enterpriseId
     */
    public function checkSwitchStatus($enterpriseId)
    {
        $extend = EnterpriseExtend::query()
            ->where([
                'enterprise_id'      => $enterpriseId,
                'cash_coupon_status' => EnterpriseExtend::OPEN
            ])->exists();

        if (!$extend) {
            throw new Error(52000, '请先开启应用“红包管理”!');
        }
    }
}