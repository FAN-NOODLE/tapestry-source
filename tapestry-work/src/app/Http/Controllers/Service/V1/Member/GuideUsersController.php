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
use App\Models\MemberUserRelation;
use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class GuideUsersController extends BaseController
{
    // 已添加
    const ALREADY_RELATION = 2;

    // 未添加
    const WAITING_RELATION = 1;

    /**
     * 导购客户列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $domain   = $request->input('qy_domain', '');
        $identify = $request->input('app_identifier', '');
        $limit    = $request->input('size', 20);
        $params   = $request->input('params', []);

        if (empty($domain)) {
            throw new Error(10000, '企业域名不能为空');
        }

        if (empty($identify)) {
            throw new Error(20000, '企业标识不能为空');
        }

        $rules = [
            'enterprise_id' => 'required',
            'mem_uid'       => 'required',
            'status'        => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'mem_uid.required'       => '导购[mem_uid]不能为空',
            'status.required'        => '客户[status]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $enterpriseId = $params['enterprise_id'];
        $status       = $params['status'];
        $memUid       = $params['mem_uid'];

        // 本地导购用户列表 // todo
        $usersList = User::query()
            ->when($status, function ($query) use ($enterpriseId, $memUid, $status) {
                // 待添加
                if (self::WAITING_RELATION == $status) {
                    $query->join('guide_user_relation', function ($join) use ($enterpriseId, $memUid) {
                        $join->on([
                            'guide_user_relation.union_id' => 'user.union_id'
                        ])->where([
                            'guide_user_relation.enterprise_id' => $enterpriseId,
                            'guide_user_relation.mem_uid'       => $memUid,
                            'guide_user_relation.deleted_at'    => 0,
                        ]);
                    })->whereNotIn('guide_user_relation.union_id', function ($query) use ($memUid, $enterpriseId) {
                        $query->select('member_user_relation.union_id')
                            ->from('member_user_relation')
                            ->where([
                                'is_binding_customer'              => MemberUserRelation::BINDING_CUSTOMER_YES,
                                'member_user_relation.delete_type' => MemberUserRelation::NORMAL_STATUS,
                                'enterprise_id'                    => $enterpriseId,
                                'mem_uid'                          => $memUid,
                                'deleted_at'                       => 0
                            ]);
                    });
                }
                // 已添加
                if (self::ALREADY_RELATION == $status) {
                    $query->join('member_user_relation', function ($join) use ($enterpriseId, $memUid) {
                        $join->on([
                            'member_user_relation.union_id' => 'user.union_id'
                        ])->where([
                            'member_user_relation.is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
                            'member_user_relation.delete_type'         => MemberUserRelation::NORMAL_STATUS,
                            'member_user_relation.enterprise_id'       => $enterpriseId,
                            'member_user_relation.mem_uid'             => $memUid,
                            'member_user_relation.deleted_at'          => 0,
                        ]);
                    });
                }
            })
            ->where([
                'user.enterprise_id' => $enterpriseId
            ])
            ->select([
                'user.id',
                'user.union_id',
                'user.type',
                'user.nick_name',
                'user.real_name',
                'user.avatar',
                'user.mobile',
            ])
            ->paginate($limit)
            ->toArray();

        return success([
            'data' => $usersList['data'],
            'meta' => [
                'total' => $usersList['total'] ?? 0,
                'size'  => $usersList['per_page'] ?? 20,
                'page'  => $usersList['current_page'] ?? 1
            ],
        ]);
    }
}
