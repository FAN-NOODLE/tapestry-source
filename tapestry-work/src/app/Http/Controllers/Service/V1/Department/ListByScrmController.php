<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:27
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Scrm\Departments;
use Exception;
use Illuminate\Http\Request;

class ListByScrmController extends BaseController
{
    /**
     * 部门列表
     *
     * @param Request $request
     *
     * @return array
     * @throws Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id', '');
        $dpIds        = $request->post('dp_ids', '');
        $name         = $request->post('name', '');
        $result       = Departments::query()
            ->when($name, function ($query) use ($name) {
                $query->where('name', 'like', '%' . $name . '%');
            })
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('third_id', $dpIds)
            ->get(['id', 'third_id', 'name', 'third_parent_id'])
            ->toArray();

        return success($result);
    }


}
