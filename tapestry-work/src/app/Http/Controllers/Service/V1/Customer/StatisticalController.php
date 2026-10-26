<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/7
 * Time: 11:15
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class StatisticalController extends BaseController
{
    /**
     * 统计客户数
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->input();
        
        // 客户数
        $count = 0;
        $where = [
            ['enterprise_id', '=', $input['enterprise_id']],
            ['is_exclusive', '=', $input['is_exclusive']],
            ['created_at', '>=', $input['created_at']]
        ];

        // 分门店查询
        if (isset($where['stores'])) {
            $count = MemberUserRelation::query()
                ->where($where)
                ->when($input['type'] === 1, function ($query) {
                    // 自然流量
                    $query->where('store_id', '=', '');
                })
                ->when($input['type'] === 2, function ($query) use ($input) {
                    // 门店
                    $query->whereIn('store_id', $input['stores']);
                })
                ->when($input['type'] === 3, function ($query) use ($input) {
                    // 自然流量+门店
                    $query->whereIn('store_id', $input['stores'])->orWhere('store_id', '=', '');
                })
                ->count();
        } else {
            $count = MemberUserRelation::query()
                ->where($where)
                ->count();
        }

        return success(['data' => ['count' => $count]]);
    }
}
