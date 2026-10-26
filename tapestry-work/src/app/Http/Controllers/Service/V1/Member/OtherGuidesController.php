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
use App\Library\Uc\Member;
use App\Models\MemberExtend;
use App\Models\MemberGuide;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class OtherGuidesController extends BaseController
{
    const DEF_MEM_ROLE = '导购';

    /**
     * 其余导购列表
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
            'enterprise_id' => 'required',
            'mem_uid'       => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'mem_uid.required'       => '导购[mem_uid]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $enterpriseId = $params['enterprise_id'];

        // 其他导购列表
        $memberList = (new MemberGuide())
            ->leftJoin('guide_user_relation', function ($join) use ($enterpriseId) {
                $join->on([
                    'guide_user_relation.mem_uid' => 'member_guide.mem_uid'
                ])->where([
                    'guide_user_relation.deleted_at' => 0
                ])->whereIn('guide_user_relation.union_id', function ($query) use ($enterpriseId) {
                    $query->select('user.union_id')
                        ->from('user')
                        ->where([
                            'enterprise_id' => $enterpriseId,
                            'deleted_at'    => 0
                        ]);
                });
            })
            ->selectRaw('
                member_guide.mem_uid,
                member_guide.limit_user_count,
                count( DISTINCT guide_user_relation.union_id ) as already_assigned_num'
            )
            ->where([
                'member_guide.enterprise_id' => $enterpriseId,
            ])
            ->where([
                ['member_guide.mem_uid', '<>', $params['mem_uid']]
            ])
            ->groupBy('member_guide.mem_uid')
            ->get()
            ->toArray();

        $memUidArr  = array_column($memberList, 'mem_uid');
        $memberList = array_column($memberList, null, 'mem_uid');

        $members = [];
        if ($memUidArr) {
            $where = [
                'pageNum'  => $params['page'] ?? 1,
                'pageSize' => $params['size'] ?? 20,
                'memUids'  => $memUidArr,
            ];

            $result = (new Member())->list($domain, $identify, $where);

            $sortArr = [];
            foreach ($result['list'] ?? [] as $member) {
                $memUid    = $member['memUid'];
                $members[] = [
                    'mem_uid'              => $memUid,
                    'mem_name'             => $member['memUsername'],
                    'mem_avatar'           => $member['memFace'] ?? '',
                    'dp_id'                => implode(' ', array_column($member['dpName'] ?? [], 'dpId')),
                    'dp_name'              => implode(' ', array_column($member['dpName'] ?? [], 'dpName')),
                    'mem_role'             => $member['memRole'],
                    'mem_gender'           => $member['memGender'],
                    'mem_mobile'           => $member['memMobile'],
                    'limit_user_count'     => $memberList[$memUid]['limit_user_count'] ?? 0,
                    'already_assigned_num' => $memberList[$memUid]['already_assigned_num'] ?? 0,
                ];
                $sortArr[] = $memberList[$memUid]['already_assigned_num'] ?? 0;
            }

            $sortArr = array_multisort($sortArr, SORT_DESC, $members);
        }

        return success([
            'data' => $members,
            'meta' => [
                'total' => $result['total'] ?? 0,
                'size'  => $result['pageSize'] ?? 10,
                'page'  => $result['pageNum'] ?? 1
            ],
        ]);
    }
}
