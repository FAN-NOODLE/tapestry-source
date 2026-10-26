<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use App\Models\EnterpriseExtend;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class RelativeMemberController extends BaseController
{
    /**
     * 根据客户信息查找所属导购(员工)
     *
     * @param Request $request
     *
     * @return array|mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $this->validate($request, [
            'union_id'      => 'required',
            'enterprise_id' => 'required',
        ], [
            'union_id.required'      => '客户唯一id不能为空',
            'enterprise_id.required' => '企业id不能为空',
        ]);
        // 客户union_id，企业id
        $params = $request->only('union_id', 'enterprise_id');

        $resultData = [];

        // 专属导购信息
        $memberUserRelation = MemberUserRelation::query()
            ->leftJoin('member_guide', 'member_user_relation.mem_uid', 'member_guide.mem_uid')
            ->where('member_user_relation.union_id', $params['union_id'])
            ->where('member_user_relation.enterprise_id', $params['enterprise_id'])
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_guide.trade_guide_status', MemberGuide::OPEN)
            ->where('member_guide.guide_id', '<>', '')
            ->where('member_user_relation.deleted_at', 0)
            ->where('member_guide.deleted_at', 0)
            ->first(['member_guide.mem_uid', 'member_guide.guide_id']);

        if ($memberUserRelation) {
            $resultData = [
                'mem_uid'  => $memberUserRelation->mem_uid,
                'guide_id' => $memberUserRelation->guide_id,
            ];
        } else {
            // 查询分配状态为开启导购
            $tradeGuidesList = MemberGuide::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('trade_guide_status', MemberGuide::OPEN)
                ->get()
                ->toArray();

            Log::debug('支付即导购集合: ' . json_encode($tradeGuidesList));

            if ($tradeGuidesList) {
                // 随机支付即导购
                $tradeGuide = $tradeGuidesList[array_rand($tradeGuidesList ,1)];

                Log::debug('支付即导购: ' . json_encode($tradeGuide));

                $resultData = [
                    'mem_uid'  => $tradeGuide['mem_uid'],
                    'guide_id' => $tradeGuide['guide_id'],
                ];
                if (!$tradeGuide['guide_id']) {
                    $resultData = [];
                }
            }
        }

        return success(['data' => $resultData ?: []]);
    }
}
