<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\MemberExtend;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberGuide;
use App\Models\GuideUserRelation;
use App\Exceptions\Error;
use Illuminate\Http\Request;

class ChangeAllocationStatusController extends BaseController
{
    /**
     * 更新员支付即导购分配状态
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $tradeGuideStatus = $request->post('trade_guide_status');
        $memUid           = $request->post('mem_uid');
        $enterpriseId     = $request->post('enterprise_id');

        if (empty($memUid)) {
            throw new Error(110005, 'memUid不能为空');
        }

        if (empty($enterpriseId)) {
            throw new Error(110006, '企业id不能为空');
        }

        $memberExtend = MemberGuide::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->first();

        if (!$memberExtend) {
            throw new Error(110007, '导购不存在');
        }

        if ($tradeGuideStatus) {
            if (!in_array($tradeGuideStatus, [MemberGuide::OPEN, MemberGuide::CLOSE])) {
                throw new Error(110010, '更新支付即导购分配状态不正确');
            }

            MemberGuide::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('mem_uid', $memUid)
                ->update(['trade_guide_status' => $tradeGuideStatus]);
        }

        return success();
    }
}
