<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:43
 */

namespace App\Http\Controllers\Service\V1\MemberExtend;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberGuide;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 批量查询员工头像
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->post('enterprise_id', 0);
        // 员工id集合
        $memUids = $request->post('mem_uids', []);

        // 头像列表
        $list = MemberGuide::query()->where(['enterprise_id' => $enterpriseId])
            ->when(!empty($memUids), function ($query) use ($memUids) {
                return $query->whereIn('mem_uid', $memUids);
            })
            ->get(['enterprise_id', 'mem_uid', 'guide_id'])
            ->toArray();

        return success(['data' => $list]);
    }
}
