<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use App\Models\User;
use Illuminate\Support\Facades\DB;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class CustomerTransferController extends BaseController
{
    /**
     * 转交类型:1部分转交
     */
    const TRANSFER_PART = 1;

    /**
     * 转交类型:3全部转交待添加
     */
    const TRANSFER_TOTAL_WAITING_ADD = 2;

    /**
     * 转交类型:3全部转交已添加
     */
    const TRANSFER_TOTAL_ALREADY_ADD = 3;


    /**
     * 转交客户到目标导购
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $domain   = $request->input('qy_domain', '');
        $identify = $request->input('app_identifier', '');
        $params   = $request->input('params', []);

        if (empty($domain)) {
            throw new Error(10000, '企业域名不能为空');
        }

        if (empty($identify)) {
            throw new Error(20000, '企业标识不能为空');
        }

        $rules = [
            'enterprise_id'  => 'required',
            'from_mem_uid'   => 'required',
            'target_mem_uid' => 'required',
            'transfer_type'  => 'required',
        ];

        $message = [
            'enterprise_id.required'  => '当前企业[enterprise_id]不能为空',
            'from_mem_uid.required'   => '当前导购[from_mem_uid]不能为空',
            'target_mem_uid.required' => '目标导购[target_mem_uid]不能为空',
            'transfer_type.required'  => '转交类型[transfer_type]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $enterpriseId    = $params['enterprise_id'];
        $fromMemUid      = $params['from_mem_uid'];
        $targetMemUid    = $params['target_mem_uid'];
        $transferType    = $params['transfer_type'];
        $unionIds        = $params['union_ids'] ?? [];
        $isThanLimitUser = false;

        $membersList = MemberGuide::query()
            ->where([
                'enterprise_id' => $enterpriseId
            ])
            ->whereIn('mem_uid', [$fromMemUid, $targetMemUid])
            ->get([
                'mem_uid',
                'guide_id',
                'limit_user_count',
            ])
            ->toArray();

        $membersList = array_column($membersList, null, 'mem_uid');

        // 1. 查询转交用户 union_id 集合 // todo
        $fromUsers = User::query()
            ->where([
                'user.enterprise_id' => $enterpriseId
            ])
            ->when($transferType, function ($query) use ($enterpriseId, $transferType, $unionIds, $fromMemUid) {
                // 部分转交用户
                if (self::TRANSFER_PART == $transferType) {
                    $query->whereIn('union_id', $unionIds);
                }
                // 全部转交已添加
                if (self::TRANSFER_TOTAL_ALREADY_ADD == $transferType) {
                    $query->join('member_user_relation', function ($join) use ($enterpriseId, $fromMemUid) {
                        $join->on([
                            'member_user_relation.union_id' => 'user.union_id'
                        ])->where([
                            'member_user_relation.is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
                            'member_user_relation.delete_type'         => MemberUserRelation::NORMAL_STATUS,
                            'member_user_relation.enterprise_id'       => $enterpriseId,
                            'member_user_relation.mem_uid'             => $fromMemUid,
                            'member_user_relation.deleted_at'          => 0,
                        ]);
                    });
                }
                // 全部转交未添加
                if (self::TRANSFER_TOTAL_WAITING_ADD == $transferType) {
                    $query->join('guide_user_relation', function ($join) use ($enterpriseId, $fromMemUid) {
                        $join->on([
                            'guide_user_relation.union_id' => 'user.union_id'
                        ])->where([
                            'guide_user_relation.enterprise_id' => $enterpriseId,
                            'guide_user_relation.mem_uid'       => $fromMemUid,
                            'guide_user_relation.deleted_at'    => 0,
                        ]);
                    })->whereNotIn('guide_user_relation.union_id', function ($query) use ($fromMemUid, $enterpriseId) {
                        $query->select('member_user_relation.union_id')
                            ->from('member_user_relation')
                            ->where([
                                'member_user_relation.is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
                                'member_user_relation.delete_type'         => MemberUserRelation::NORMAL_STATUS,
                                'member_user_relation.enterprise_id'       => $enterpriseId,
                                'member_user_relation.mem_uid'             => $fromMemUid,
                                'member_user_relation.deleted_at'          => 0
                            ]);
                    });
                }
            })
            ->get([
                'user.id',
                'user.union_id'
            ])
            ->toArray();

        if (!$fromUsers) {
            throw new Error(30000, '客户非来源导购客户');
        }

        $unionIdArr = array_column($fromUsers, 'union_id');

        $exists = GuideUserRelation::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'mem_uid'       => $targetMemUid
            ])
            ->whereIn('union_id', $unionIdArr)
            ->exists();

        if ($exists) {
            throw new Error(40000, '存在客户已被分配到目标导购');
        }

        // 2. 相关数据校验
        $assignedCount = GuideUserRelation::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'mem_uid'       => $targetMemUid,
                'guide_id'      => $membersList[$targetMemUid]['guide_id']
            ])
            ->whereIn('guide_user_relation.union_id', function ($query) use ($enterpriseId) {
                $query->select('user.union_id')
                    ->from('user')
                    ->where([
                        'user.enterprise_id' => $enterpriseId,
                        'user.deleted_at'    => 0
                    ]);
            })
            ->count();

        $count = self::TRANSFER_PART == $transferType ? count($unionIdArr) : 1;
        if ($assignedCount + $count > $membersList[$targetMemUid]['limit_user_count']) {
            throw new Error(50000, '目标导购客户分配超过上限');
        }

        // 3. 获取目标导购和用户的关系, 转交数量$key + 1 = userCount, 初始化unionIdArr(只删除转交用户)
        $insert     = [];
        $unionIdArr = [];
        foreach ($fromUsers as $key => $user) {
            if ($key + 1 <= $membersList[$targetMemUid]['limit_user_count'] - $assignedCount) {
                $insert[]     = [
                    'enterprise_id' => $enterpriseId,
                    'mem_uid'       => $targetMemUid,
                    'user_id'       => $user['id'],
                    'union_id'      => $user['union_id'],
                    'guide_id'      => $membersList[$targetMemUid]['guide_id'],
                    'relation_at'   => time(),
                    'created_at'    => time(),
                    'updated_at'    => time(),
                ];
                $unionIdArr[] = $user['union_id'];
            } else {
                $isThanLimitUser = true;
                continue;
            }
        }

        // 4. 分配关系数据处理,外部联系人关系处理
        DB::beginTransaction();
        try {
            // 解绑来源导购和目标客户的外部联系人绑定关系
            MemberUserRelation::query()
                ->where([
                    'enterprise_id' => $enterpriseId,
                    'mem_uid'       => $fromMemUid,
                ])
                ->whereIn('union_id', $unionIdArr)
                ->update([
                    'is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_NO
                ]);
            // 建立目标导购和目标客户的外部联系人的绑定关系
            MemberUserRelation::query()
                ->where([
                    'enterprise_id' => $enterpriseId,
                    'mem_uid'       => $targetMemUid,
                ])
                ->whereIn('union_id', $unionIdArr)
                ->update([
                    'is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES
                ]);
            // 删除来源导购和转交用户的关系
            GuideUserRelation::query()
                ->where([
                    'enterprise_id' => $enterpriseId,
                    'mem_uid'       => $fromMemUid,
                    'guide_id'      => $membersList[$fromMemUid]['guide_id'],
                ])
                ->whereIn('union_id', $unionIdArr)
                ->delete();

            // 5. 新增目标导购和用户的关系
            GuideUserRelation::query()->insert($insert);
            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error(60000, '转交客户异常,请稍后重试');
        }

        return success(['data' => [
            'number' => $isThanLimitUser ? count($unionIdArr) : 0,
        ]]);
    }

}
