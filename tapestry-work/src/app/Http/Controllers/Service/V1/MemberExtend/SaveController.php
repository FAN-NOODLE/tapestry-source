<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:24
 */

namespace App\Http\Controllers\Service\V1\MemberExtend;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SaveController extends BaseController
{
    /**
     * 保存员工头像
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->post('enterprise_id', 0);
        // 员工id
        $memUid = $request->post('mem_uid', '');

        // 客服在微信的guide_id
        $guideId = $request->post('guide_id', '');

        $contactQrCodeId = $request->post('contact_qr_code_id', 0);


        Log::info('请求参数', $request->all());

        if ($enterpriseId && $memUid) {
            // 待更新数据
            $data = [
                'enterprise_id'   => $enterpriseId,
                'mem_uid'         => $memUid,
                'guide_id'        => $guideId,
                'contact_qr_code_id' => $contactQrCodeId
            ];

            // 更新
            MemberGuide::updateOrCreate([
                'enterprise_id' => $enterpriseId,
                'mem_uid'       => $memUid
            ], $data);

            GuideUserRelation::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('mem_uid', $memUid)
                ->update(['guide_id' => $guideId]);
        }

        return success(['data' => []]);
    }
}
