<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\ContactQrCode;
use App\Models\EnterpriseExtend;
use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;
use App\Models\ContactQrCodeParams;

class GuideController extends BaseController
{

    /**
     * 支付即导购列表
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
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $enterpriseId = $params['enterprise_id'];

        // 先查询导购列表 基本信息
        $memberList = MemberGuide::query()
            ->where([
                'enterprise_id' => $enterpriseId,
            ])
            ->get([
                'mem_uid',
                'contact_qr_code_id',
                'allocation_status',
                'trade_guide_status',
                'limit_user_count',
            ])
            ->toArray();

        $memberList = array_column($memberList, null, 'mem_uid');

        $memberUIds = array_column($memberList, 'mem_uid');
        $contactQrCodeIdList = array_filter(array_column($memberList, 'contact_qr_code_id'));

        $guideCountMap = [];
        if (!empty($memberUIds) && !empty($contactQrCodeIdList)) {
            $guideCountUserList = MemberUserRelation::query()
                ->select(DB::raw('mem_uid, count(*) as count'))
                ->where('enterprise_id', $enterpriseId)
                ->whereIn('mem_uid', $memberUIds)
                ->whereIn('source_type', $contactQrCodeIdList)
                ->groupBy('mem_uid')
                ->get()
                ->toArray();
            $guideCountMap = array_column($guideCountUserList, 'count', 'mem_uid');
        }

        $members = [];
        if ($memberUIds) {
            $where = [
                'pageNum'     => $params['page'] ?? 1,
                'pageSize'    => $params['size'] ?? 20,
                'memUids'     => $memberUIds,
            ];

            // 拼接请求参数
            if (!empty($params['mem_mobile'])) {
                $where['memMobile'] = $params['mem_mobile'];
            }
            if (!empty($params['mem_user_name'])) {
                $where['memUsername'] = $params['mem_user_name'];
            }
            if (!empty($params['dp_id'])) {
                $where['dpId'] = $params['dp_id'];
            }

            $result = (new Member())->ignorePermissionsList($domain, $identify, $where);

            $sortArr = [];
            foreach ($result['list'] ?? [] as $member) {
                $memUid    = $member['memUid'];
                $members[] = [
                    'mem_uid'              => $memUid,
                    'mem_name'             => $member['memUsername'],
                    'mem_avatar'           => $member['memFace'] ?? '',
                    'dp_id'                => implode('；', array_column($member['dpName'] ?? [], 'dpId')),
                    'dp_name'              => implode('；', array_column($member['dpName'] ?? [], 'dpName')),
                    'mem_role'             => $member['memRole'],
                    'mem_gender'           => $member['memGender'],
                    'mem_mobile'           => $member['memMobile'],
                    'trade_guide_status'   => $memberList[$memUid]['trade_guide_status'] ?? 0,
                    'already_relation_num' => $guideCountMap[$memUid] ?? 0,
                ];
                $sortArr[] = $guideCountMap[$memUid] ?? 0;
            }

            $sortArr && array_multisort($sortArr, SORT_DESC, $members);
        }

        $data = [
            'list'  => $members,
        ];

        return success([
            'data' => $data,
            'meta' => [
                'total' => $result['total'] ?? 0,
                'size'  => $result['pageSize'] ?? 20,
                'page'  => $result['pageNum'] ?? 1
            ],
        ]);
    }
}
