<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\GuideManage;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberGuide;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class GuideBaseLIstController extends BaseController
{

    /**
     * 导购员工列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->input('params', []);
        $page   = $params['page'] ?? 1;
        $size   = $params['size'] ?? 10;
        $rules  = [
            'enterprise_id' => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $memberSql = MemberGuide::query()
            ->where([
                'enterprise_id' => $params['enterprise_id'],
            ])
//            ->where('guide_id', '<>', '')
            ->when(!empty($params['mem_uid']), function ($query) use ($params) {
                return $query->where('mem_uid', $params['mem_uid']);
            })
            ->when(!empty($params['mem_uids']), function ($query) use ($params) {
                return $query->whereIn('mem_uid', $params['mem_uids']);
            });

        $total      = $memberSql->count();
        $memberList = $memberSql->forpage($page, $size)
            ->get()
            ->toArray();

        $data = [
            'data' => $memberList,
            'meta' => [
                'total' => $total ?? 0,
                'size'  => intval($page),
                'page'  => intval($size)
            ],
        ];
        return success($data);
    }
}
