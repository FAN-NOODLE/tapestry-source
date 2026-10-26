<?php
/**
 * 任务类型列表
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 17:52
 */
namespace App\Http\Controllers\Service\V1\Task;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\TaskType;
use Illuminate\Http\Request;

class TaskTypeListController extends BaseController
{
    public function index(Request $request)
    {
        $page = $request->post('page', 1);
        $size = $request->post('size', 10);
        $search = $request->post('search', []);

        // 检测服务类型，若服务类型为shop+crs,去除条件
        if (isset($search['serve_type'])
            && !empty($search['serve_type'])
            && $search['serve_type'] == 0
        ) {
            unset($search['serve_type']);
        }

        $orm = $this->getData($page, $size, $search);

        $result = [
            'data' => $orm->items(),
            'meta' => [
                'page' => $page,
                'size' => $size,
                'total' => $orm->total() ?? 0
            ]
        ];
        return success($result);
    }

    /**
     * 获取数据分页模型
     * @param $page
     * @param $size
     * @param $search
     * @return \Illuminate\Contracts\Pagination\LengthAwarePaginator
     */
    public function getData($page, $size, $search)
    {
        $orm = TaskType::query()
                ->when(isset($search['name']) && !empty($search['name']), function ($query) use ($search) {
                    return $query->{'where'}('name', 'like', '%'.$search['name'].'%');
                })
                ->when(isset($search['serve_type']) && !empty($search['serve_type']), function ($query) use ($search) {
                    return $query->{'where'}('serve_type', '=', $search['serve_type']);
                })
                ->paginate($size, ['id', 'name', 'describe', 'serve_type'], '', $page);
        return $orm;
    }
}
