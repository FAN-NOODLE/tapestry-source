<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/22 17:25
 */

namespace App\Http\Controllers\V1\Setting\OperationLog;

use App\Http\Controllers\V1\BaseController;
use App\Models\OperationLog;
use Illuminate\Http\Request;

class ListController extends BaseController
{

    /**
     * 角色列表接口（分页、搜索）
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $page   = (int)$request->input('page', 1);
        $size   = (int)$request->input('size', 10);
        $search = json_decode($request->input('search', '{}'), true);

        $where = [
            [
                'enterprise_id',
                '=',
                $this->enterpriseId
            ]
        ];
        // 操作人查询
        if (isset($search['admin_name']) && !empty($search['admin_name'])) {
            $where[] = [
                'admin_name',
                'LIKE',
                $search['admin_name'] . '%'
            ];
        }
        // 时间段查询
        if (isset($search['starttime']) && !empty($search['starttime']) && isset($search['endtime']) && !empty($search['endtime'])) {
            $where[] = [
                'created_at',
                '>',
                $search['starttime']
            ];
            $where[] = [
                'created_at',
                '<',
                $search['endtime']
            ];
        }
        // 查询模块
        if (isset($search['module']) && !empty($search['module'])) {
            $where[] = [
                'module',
                'like',
                '%' . $search['module'] . '%'
            ];
        }
        $data = OperationLog::query()->select([
            'admin_name',
            'module',
            'page',
            'content',
            'ip',
            'http_method',
            'user_agent',
            'created_at'
        ])->where($where)->orderBy('id', 'DESC')->forpage($page, $size)->get()->toArray();

        $total = OperationLog::query()->where($where)->count();

        return success([
            'list'  => $data,
            'total' => $total,
            'page'  => $page,
            'size'  => $size
        ]);
    }

}
