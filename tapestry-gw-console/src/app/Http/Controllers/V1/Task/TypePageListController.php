<?php
/**
 * 任务类型列表-分页
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/3/1
 * Time: 10:06
 */
namespace App\Http\Controllers\V1\Task;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Task;
use Illuminate\Http\Request;

class TypePageListController extends BaseController
{
    /**
     *  任务类型列表-分页
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $page = $request->post('page', 1);
        $size = $request->post('size', 10);
        $search = $request->post('search', []);

        // 服务类型判断
        if (!isset($search['serve_type']) || empty($search['serve_type'])) {
            $search['serve_type'] = $this->serveType;
        }

        $result = (new Task())->taskTypeListPage([
            'page' => $page,
            'size' => $size,
            'search' =>$search
        ]);

        return success($result);
    }
}
