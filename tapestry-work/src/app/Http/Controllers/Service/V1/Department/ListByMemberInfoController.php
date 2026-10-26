<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:27
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Scrm\Members;
use Illuminate\Http\Request;

class ListByMemberInfoController extends BaseController
{
    /**
     * 根据员工信息查询部门列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id', '');
        $params = $request->only(['enterprise_id', 'mobile', 'userid', 'name']);

        $list = Members::query()->with('departmentsRel.departments')
            ->where('enterprise_id', $enterpriseId)
            ->when(!empty($params['mobile']), function ($query) use ($params) {
                return $query->whereIn('mobile', is_array($params['mobile']) ? $params['mobile'] : [$params['mobile']]);
            })
            ->when(!empty($params['userid']), function ($query) use ($params) {
                return $query->whereIn('userid', is_array($params['userid']) ? $params['userid'] : [$params['userid']]);
            })
            ->when(!empty($params['name']), function ($query) use ($params) {
                return $query->where('name', $params['name']);
            })
            ->get(['id'])
            ->toArray();
        $returnArr = [];
        foreach ($list as $val) {
            if (!empty($val['departments_rel'])) {
                foreach ($val as $rel) {
                    $returnArr[] = $rel['departments'];
                }
            }
        }
        return success($returnArr);
    }


}
