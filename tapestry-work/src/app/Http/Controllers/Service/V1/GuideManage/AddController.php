<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:29
 */

namespace App\Http\Controllers\Service\V1\GuideManage;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use Illuminate\Http\Request;
use App\Exceptions\Error;

class AddController extends BaseController
{
    const DEF_MEM_ROLE = '导购';

    /**
     * 新增导购
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->post('enterprise_id', 0);

        // 员工id集合
        $memUids = $request->post('mem_uids', []);
        if (!is_array($memUids) || empty($memUids)) {
            throw new Error(100000, '导购集合数据格式不正确!');
        }
        $memUids = array_unique($memUids);


        $allMenUid = [];
        $list      = MemberGuide::query()
            ->withTrashed()
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('mem_uid', $memUids)
            ->get(['mem_uid'])
            ->toArray();

        if (!empty($list)) {
            $allMenUid = array_unique(array_column($list, 'mem_uid'));
        }

        $guideData = [];
        foreach ($memUids as $memUid) {
            if (in_array($memUid, $allMenUid)) {
                $update = [
                    'enterprise_id'      => $enterpriseId,
                    'mem_uid'            => $memUid,
                    'deleted_at'         => 0,
                    'trade_guide_status' => 2
                ];

                $info = MemberGuide::query()
                    ->withTrashed()
                    ->where([
                        'enterprise_id' => $enterpriseId,
                        'mem_uid'       => $memUid,
                    ])->first();
                
                // 更新状态
                MemberGuide::query()
                    ->withTrashed()
                    ->where('id', $info->id)->update($update);
                continue;
            }

            $guideData[] = [
                'enterprise_id' => $enterpriseId,
                'mem_uid'       => $memUid,
                'created_at'    => time()
            ];
        }


        if (!empty($guideData)) {
            MemberGuide::query()->insert($guideData);
        }

        $data = ['data' => []];

        return success($data);
    }
}
