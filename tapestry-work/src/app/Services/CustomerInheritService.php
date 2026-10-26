<?php

namespace App\Services;


use App\Jobs\InheritTransferJob;
use App\Library\Uc\QyWechat;
use App\Micros\Common;
use App\Models\Scrm\CustomerOnJobInherit;
use App\Models\Scrm\CustomerOnJobInheritRecord;
use App\Models\Scrm\Departments;
use App\Models\Scrm\ImpressionTagRelation;
use App\Models\Scrm\MemberDepartments;
use App\Models\Scrm\Members;
use App\Models\Scrm\WorkExternalUserRelation;
use App\Models\Scrm\WorkExternalUserTag;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use Exception;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class CustomerInheritService
{
    /**
     * @param $params
     *
     * @return mixed
     * 查导购信息
     */
    public function getMember($params)
    {
        return Members::query()->where([
            'enterprise_id' => $params['enterprise_id'],
            'mobile' => $params['mobile']
        ])->select('id', 'name')->first();
    }


    public function getMemberIds($data)
    {
        $department_id = Departments::query()
            ->where(['enterprise_id' => $data['enterprise_id']])
            ->whereIn('third_id', $data['third_ids'])
            ->pluck("id")->toArray();

        return MemberDepartments::query()
            ->where(['enterprise_id' => $data['enterprise_id']])
            ->whereIn('department_id', $department_id)->pluck('member_id')->toArray();
    }

    /**
     * @param $params
     * @param $memberIds
     * @return array
     */
    public function inheritList($params, $memberIds): array
    {
        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 10;

        $tops['on'] = $this->getNum($params, CustomerOnJobInherit::STATUS_ON);
        $tops['success'] = $this->getNum($params, CustomerOnJobInherit::STATUS_SUCCESS);
        $tops['err'] = $this->getNum($params, CustomerOnJobInherit::STATUS_ERR);

        $query = CustomerOnJobInherit::query()->where([
            'enterprise_id' => $params['enterprise_id']
        ])->whereIn('member_id', $memberIds);

        if (isset($params['status'])) {
            $query->where('status', $params['status']);
        }

        if (isset($params['created_at']) && $params['created_at']) {
            $query->whereBetween('created_at', $params['created_at']);
        }
        $data = $query->orderBy('created_at', 'desc')->paginate($size, ['*'], 'page', $page);

        foreach ($data->items() as &$value) {
            $value['takeover_name'] = $value->info->takeover_name ?? "";
            $value['is_self'] = $value['member_id'] == $params['member_id'] ? 1 : 0;
            unset($value->info);
        }

        return [
            'top' => $tops,
            'list' => $data->items(),
            'total' => $data->total(),
            'page' => $data->currentPage(),
            'size' => $data->perPage()
        ];
    }

    /**
     * @param $params
     * @param $status
     * 分配数统计
     *
     * @return int
     */
    public function getNum($params, $status)
    {
        if (isset($params['status'])) {
            if ($params['status'] != $status) {
                return 0;
            }
        }

        $query = CustomerOnJobInherit::query()->where([
            'enterprise_id' => $params['enterprise_id'],
            'member_id' => $params['member_id'],
            'status' => $status
        ]);

        if (isset($params['created_at']) && $params['created_at']) {
            $query->whereBetween('created_at', $params['created_at']);
        }
        return $query->count();
    }

    /**
     * @param $params
     *
     * @return array
     */
    public function inheritInfo($params): array
    {
        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 10;

        $inherit = CustomerOnJobInherit::query()->where([
            'enterprise_id' => $params['enterprise_id'],
            'id' => $params['inherit_id']
        ])->first(['id', 'member_name', 'customers_numbers', 'created_at']);

        if (!$inherit) {
            error("任务未查询到");
        }

        $data = CustomerOnJobInheritRecord::query()->where([
            'enterprise_id' => $params['enterprise_id'],
            'inherit_id' => $params['inherit_id']
        ])->select('id', 'inherit_id', 'handover_name', 'takeover_name', 'status', 'external_name', 'errmsg', 'created_at', 'takeover_time')
            ->paginate($size, ['*'], 'page', $page);


        return [
            'top' => $inherit,
            'list' => $data->items(),
            'total' => $data->total(),
            'page' => (int)$data->currentPage(),
            'size' => (int)$data->perPage()
        ];

    }

    /**
     * @param $params
     *
     * @return bool
     */
    public function inheritTransfer($params): bool
    {
        $memberIds = [];
        $workExternalUserId = [];
        if (isset($params['is_all_checked']) && $params['is_all_checked'] == 1) { // 全选的逻辑
            $relation = Common::microScrm()->setApi("micro-scrm/v1/customer/on-inherit/select-all-list")->post($params);
        } else { // 非全选的在职接替
            //查外部联系人和导购的关系 relation_status 0导购和外部联系人是正常关系
            $relation = WorkExternalUserRelation::query()
                ->where(['enterprise_id' => $params['enterprise_id'], 'relation_status' => 0])
                ->whereIn('id', $params['relation_id'])
                ->select('enterprise_id', 'member_id', 'work_external_user_id', 'user_id')
                ->get()->toArray();
            //校验顾客数据
            if (count($params['relation_id']) != count($relation)) {
                error("选择顾客数量与查询结果不匹配");
            }
        }
        //过滤分配中的用户
        $externalUserIds = array_column($relation, 'work_external_user_id');
        $waitUserIds = CustomerOnJobInheritRecord::query()->whereIn('external_userid', $externalUserIds)
            ->where('status', CustomerOnJobInheritRecord::STATUS_WAIT)
            ->pluck('external_userid')
            ->toArray();
        foreach ($relation as $relationInfo) {
            if (!in_array($relationInfo['work_external_user_id'], $waitUserIds)) {
                $newRelation[] = $relationInfo;
            }
            $memberIds[] = $relationInfo['member_id'];
            $workExternalUserId[] = $relationInfo['work_external_user_id'];
        }

        //根据导购查导购信息(将接替的用户也查出来)
        $memberIds[] = $params['takeover_member_id'];
        $memberIds = array_unique($memberIds);
        $members = Members::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('id', $memberIds)
            ->select('id', 'userid', 'name', 'mobile')
            ->get()
            ->toArray();
        $members = collect($members)->keyBy('id')->all();

        if (empty($members[$params['takeover_member_id']])) {
            error("接替导购不存在");
        }
        //查外部联系人信息
        $external = WxUserCorpExternal::getExternalUserInfoByIds(
            $params['enterprise_id'],
            array_unique($workExternalUserId),
            ['id', 'openid', 'nick', 'name', 'avatar']
        );
        $external = collect($external)->keyBy('id')->all();
        $params['transfer_success_msg'] = str_replace("【接替者名字】", $members[$params['takeover_member_id']]['name'], $params['transfer_success_msg']);
        //写入
        $inheritID = idGenerate();
        $nowUnix = time();
        $customersNumbers = count($relation);
        DB::beginTransaction();
        try {
            $inherit = CustomerOnJobInherit::query()->insert([
                'id' => $inheritID,
                'enterprise_id' => $params['enterprise_id'],
                'member_id' => $params['member_id'],
                'member_name' => $params['member_name'],
                'errmsg' => '分配中',
                'customers_numbers' => $customersNumbers,
                'transfer_success_msg' => $params['transfer_success_msg'],
                'created_at' => $nowUnix,
                'updated_at' => $nowUnix,
            ]);

            $recordArr = [];
            $relationNums = count($newRelation);
            $ids = [];
            if ($relationNums > 1000) {
                $maxPage = (int)ceil($relationNums / 1000);
                for ($i = 1; $i <= $maxPage; $i++) {
                    $ids = array_merge($ids, uids(1000));
                }
            } else {
                $ids = uids($relationNums);
            }

            foreach ($newRelation as $key => $item) {
                $record['id'] = $ids[$key];
                $record['inherit_id'] = $inheritID;
                $record['enterprise_id'] = $params['enterprise_id'];
                $record['handover_member_id'] = $members[$item['member_id']]['id'] ?? '';
                $record['handover_userid'] = $members[$item['member_id']]['userid'] ?? '';
                $record['handover_name'] = $members[$item['member_id']]['name'] ?? '';
                $record['handover_mobile'] = $members[$item['member_id']]['mobile'] ?? '';

                $record['takeover_member_id'] = $members[$params['takeover_member_id']]['id'] ?? '';
                $record['takeover_userid'] = $members[$params['takeover_member_id']]['userid'] ?? '';
                $record['takeover_name'] = $members[$params['takeover_member_id']]['name'] ?? '';
                $record['takeover_mobile'] = $members[$params['takeover_member_id']]['mobile'] ?? '';

                $record['external_id'] = $external[$item['work_external_user_id']]['id'] ?? '';
                $record['external_userid'] = $external[$item['work_external_user_id']]['openid'] ?? '';
                $record['external_nickname'] = $external[$item['work_external_user_id']]['nick'] ?? '';
                $record['external_name'] = $external[$item['work_external_user_id']]['name'] ?? '';
                $record['external_avatar'] = $external[$item['work_external_user_id']]['avatar'] ?? '';
                $record['errmsg'] = '等待接替';
                $record['created_at'] = $nowUnix;
                $record['updated_at'] = $nowUnix;

                $recordArr[] = $record;
            }

            $recordRes = CustomerOnJobInheritRecord::query()->insert($recordArr);
            if ($inherit && $recordRes) {
                DB::commit();
                dispatch(new InheritTransferJob($inheritID, $params['transfer_success_msg'], $params['enterprise_id']));
                return true;
            } else {
                DB::rollBack();
                return false;
            }

        } catch (Exception $exception) {
            DB::rollBack();
            Log::info('在职继承业务写入数据异常' . $exception->getMessage());
            return false;
        }
    }


    /**
     * @return void
     * 分配
     */
    public function transferCustomer($inherit_id, $transfer_success_msg, $enterprise_id)
    {
        $data = CustomerOnJobInheritRecord::query()
            ->where(['inherit_id' => $inherit_id, 'status' => CustomerOnJobInheritRecord::STATUS_WAIT])
            ->get();
        $items = [];

        $collection = collect($data->toArray());
        $handover_member_id = $collection->unique('handover_member_id')->all();
        foreach ($handover_member_id as $group) {
            $info = [
                'handover_userid' => $group['handover_userid'],
                'takeover_userid' => $group['takeover_userid'],
                'external_userid' => [],
                'transfer_success_msg' => $transfer_success_msg
            ];

            //循环组装分页数据
            $i = 0;
            $external_userid = [$i => []];
            foreach ($data as $row) {
                if ($group['handover_member_id'] == $row['handover_member_id']) {
                    if (count($external_userid[$i]) >= 100) {
                        $i++;
                    }
                    $external_userid[$i][] = $row['external_userid'];
                }
            }

            foreach ($external_userid as $value) {
                $info['external_userid'] = $value;
                $items[] = $info;
            }
        }

        //调用企微接口
        foreach ($items as $value) {
            try {
                $res = (new QyWechat())->getQyResult($value, [
                    'path' => 'cgi-bin/externalcontact/transfer_customer',
                    'method' => 'POST',
                    'enterprise_id' => $enterprise_id
                ]);
                Log::info('在职继承发起企微接口返回信息' . json_encode($res));

                if ($res['errcode'] == 0) {
                    foreach ($res['customer'] as $customer) {
                        $customerMsg = ['request_callback_code' => $customer['errcode']];
                        if ($customer['errcode'] == '40129') {
                            $customerMsg['errmsg'] = '当前客户正在转接中';
                            $customerMsg['status'] = CustomerOnJobInheritRecord::STATUS_WAIT_ERROR;
                        }
                        CustomerOnJobInheritRecord::query()->where([
                            'inherit_id' => $inherit_id,
                            'external_userid' => $customer['external_userid'],
                        ])->update($customerMsg);
                    }
                }
            } catch (Exception $exception) {
                Log::info("在职继承发起企微调用异常" . $exception->getMessage());
            }
        }
    }


    /**
     * @param $params
     *
     * @return void
     * 导购信息
     */
    public function memberList($params): array
    {
        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 10;
        $arr = [];

        $department = Departments::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('third_id', $params['dp_ids'])
            ->select('id', 'name')
            ->get()
            ->toArray();

        //店铺id集合
        $department_id = collect($department)->pluck('id')->toArray();
        //查询导购自己防止导购管理范围不在自己的门店
        $memberDepartmentID = MemberDepartments::query()
            ->where(['enterprise_id' => $params['enterprise_id'], 'member_id' => $params['member_id']])
            ->value('department_id');
        if ($memberDepartmentID) {
            $department_id[] = $memberDepartmentID;
        }

        //店铺名称集合
        $storeNames = collect($department)->keyBy('id')->all();

        $sql = MemberDepartments::query()
            ->whereIn('department_id', $department_id)
            ->with('members')
            ->groupBy('member_id');
        if (isset($params['search']) && $params['search']) {
            $sql->whereHas('members', function ($query) use ($params) {
                $query->where('name', 'like', '%' . $params['search'] . '%')
                    ->orWhere('userid', 'like', '%' . $params['search'] . '%');
            });
        }
        $dpList = $sql->paginate($size, ['*'], 'page', $page);

        foreach ($dpList as $rel) {
            //用户信息
            if ($rel->members) {
                $store_name = $storeNames[$rel->department_id]['name'] ?? '';
                if (isset($arr[$rel->members->id])) {
                    $arr[$rel->members->id]['store_name'][] = $store_name;
                } else {
                    $arr[$rel->members->id] = [
                        'member_id' => $rel->members->id ?? '',
                        'name' => $rel->members->name ?? '',
                        'userid' => $rel->members->userid ?? '',
                        'avatar' => $rel->members->avatar ?? '',
                        'enterprise_id' => $rel->members->enterprise_id ?? '',
                        'store_name' => [$store_name],
                        'members' => 0,
                        'check' => true,
                    ];
                }
            }
        }

        //查询关联的人数
        $nums = WorkExternalUserRelation::query()
            ->where(['enterprise_id' => $params['enterprise_id'], 'relation_status' => 0, 'member_deleted_time' => 0, 'external_deleted_time' => 0])
            ->whereIn('member_id', collect($arr)->keys())
            ->selectRaw("count(*) as num,member_id")
            ->groupBy('member_id')
            ->get();

        //转化格式给用户数据组装人数
        $nums = collect($nums)->keyBy('member_id')->all();
        foreach ($nums as $k => $value) {
            $arr[$k]['members'] = $value['num'];
        }

        return [
            'list' => collect($arr)->values(),
            'total' => $dpList->total(),
            'page' => $dpList->currentPage(),
            'size' => $dpList->perPage()
        ];
    }

    /**
     * @return void
     * 外部联系人转移的次数
     */
    public function customerInheritRecord($params): array
    {
        $time = time() - 60 * 60 * 24 * 90;
        return CustomerOnJobInheritRecord::query()->where(['enterprise_id' => $params['enterprise_id']])
            ->whereIn('external_id', $params['external_id'])
            ->whereIn('status', [CustomerOnJobInheritRecord::STATUS_SUCCESS, CustomerOnJobInheritRecord::STATUS_WAIT, CustomerOnJobInheritRecord::STATUS_REJECT])
            ->where('created_at', '>', $time)
            ->selectRaw("count(*) as num,external_id")
            ->groupBy('external_id')
            ->get()->toArray();
    }

    /**
     * @return void
     * 外部联系人等待接替中的用户
     */
    public function customerInheritWait($params): array
    {
        $time = time() - 60 * 60 * 24;
        return CustomerOnJobInheritRecord::query()->where(['enterprise_id' => $params['enterprise_id']])
            ->whereIn('external_id', $params['external_id'])
            ->where('status', CustomerOnJobInheritRecord::STATUS_WAIT)
            ->where('created_at', '>', $time)
            ->pluck('external_id')
            ->toArray();
    }

    /**
     * @param $params
     *
     * @return array
     */
    public function memberRecord($params): array
    {
        try {
            //查到管理权限内的所有导购id
            //如果传递了导购id那就直接用导购id来查  这个条件是判断选择了导购
            $membersID = [];
            $memberList = $params['params']['member_list'] ?? '';
            if (!$memberList) {
                $members = Departments::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->whereIn('third_id', $params['third_ids'])
                    ->with('membersRel')
                    ->get();
                foreach ($members as $items) {
                    foreach ($items->membersRel as $v) {
                        $membersID[] = $v['member_id'];
                    }
                }
            } else {
                $membersID[] = $memberList;
            }

            $isKeySa = $params['params']['manage_type'] ?? '';
            if ($isKeySa) {
                $isKeySa = $isKeySa == 1 ? 1 : 0;
                //如果是专属导购只能是自己看见的
                if ($isKeySa) {
                    $membersID[] = Members::query()->where([
                        'enterprise_id' => $params['enterprise_id'],
                        'mobile' => $params['mobile']
                    ])->value('id');
                }
            }

            $query = WorkExternalUserRelation::query()
                ->with('relationToMember')
                ->where(['enterprise_id' => $params['enterprise_id'], 'relation_status' => 0])
                ->whereIn('work_external_user_id', $params['work_external_user_id'])
                ->whereIn('member_id', $membersID)
                ->select('id', 'work_external_user_id', 'member_id', 'remark', 'is_key_sa', 'created_at');

            //判断是不是专属顾问
            if ($isKeySa) {
                $isKeySa = $isKeySa == 1 ? 1 : 0;
                $query->where('is_key_sa', $isKeySa);
            }

            $data = $query->get();
            foreach ($data as $value) {
                $value->member_name = $value->relationToMember->name ?? '';
                unset($value->relationToMember);
            }

            $data = $data->toArray();

            //客户标签检索
            if ($params['params']['tag_list'] ?? '') {

                $relationId = collect($data)->pluck('id')->all();

                $userTag = WorkExternalUserTag::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->whereIn('relation_id', $relationId)
                    ->whereIn('external_tag_id', $params['params']['tag_list'])
                    ->pluck('relation_id')
                    ->toArray();

                foreach ($data as $k => $row) {
                    if (!in_array($row['id'], $userTag, true)) {
                        unset($data[$k]);
                    }
                }
            }

            //印象标签
            if ($params['params']['impression_tag_list'] ?? '') {

                $relationId = collect($data)->pluck('id')->all();

                $impressionTag = ImpressionTagRelation::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->whereIn('external_user_id', $relationId)
                    ->whereIn('tag_id', $params['params']['impression_tag_list'])
                    ->pluck('external_user_id')
                    ->toArray();

                foreach ($data as $key => $item) {
                    if (!in_array($item['id'], $impressionTag, true)) {
                        unset($data[$key]);
                    }
                }
            }

            //会员等级CLT类型
            //优质客 card_type = 10  条件 // 优质客是 clt标签 有值
            //潜在客 card_type = 20  条件 amount_total = 0 and card_level ="PROSPECT"
            $card = $params['params']['card_type'] ?? '';
            $clt = $params['params']['clt_type'] ?? '';
            if ($card || $clt) {
                $ids = collect($data)->pluck('work_external_user_id')->all();

                $cardAndClt = WxUserCorpExternal::query()->setConn($params['enterprise_id'])->whereIn('id', $ids);

                if (in_array(WxUserCorpExternal::CARD_LEVEL_HIGH, $card)) {
                    //优质客
//                    $cardAndClt->where('amount_total', '>', 0);
                    $cardAndClt->where('clt_type', '<>', '');

                } elseif (in_array(WxUserCorpExternal::CARD_LEVEL_POTENTIAL, $card)) {
                    //潜在客
                    $cardAndClt->where(['amount_total' => 0, 'card_level' => 'PROSPECT']);

                } elseif ($card) {
                    $cardAndClt->whereIn('card_level', $card);
                }

                if ($clt) {
                    $cardAndClt->whereIn('clt_type', $clt);
                }

                $cardAndCltData = $cardAndClt->pluck('id')->toArray();
                foreach ($data as $key => $item) {
                    if (!in_array($item['id'], $cardAndCltData, true)) {
                        unset($data[$key]);
                    }
                }
            }

            return $data;
        } catch (Exception $exception) {
            return [];
        }
    }

    /**
     * 查询分配顾客成功的记录
     *
     * @param $params
     *
     * @return array
     * @Time: 2022/11/17
     */
    public function customerInheritRecordForExternalIds($params): array
    {
        $startTime = strtotime(date("Y-m-d", strtotime("-7 day")));
        $data = CustomerOnJobInheritRecord::query()
            ->where(['enterprise_id' => $params['enterprise_id'], 'status' => CustomerOnJobInheritRecord::STATUS_SUCCESS])
            ->where('takeover_time', '>', $startTime)
            ->whereIn('external_id', $params['external_ids'])
            ->get(['takeover_member_id', 'takeover_userid', 'takeover_mobile', 'external_id', 'external_userid'])
            ->toArray();

        return $data;
    }
}