<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Models\Enterprise;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class ListNoPageController extends BaseController
{
    /**
     * 企业列表（不分页）
     *
     * @date 2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->input();

        $enterprise = Enterprise::query()
            ->when(isset($params['enterprise_ids']), function ($query) use ($params) {
                return $query->whereIn('id', $params['enterprise_ids']);
            })->when(isset($params['name']), function ($query) use ($params) {
                return $query->where('name', 'like', '%' . $params['name'] . '%');
            })->get()
            ->toArray();

        return success(
            ['data' => (array)$enterprise]
        );
    }
}
