<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 15:04
 */

namespace App\Http\Controllers\V1\Task;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Task;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 任务列表
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $data['page'] = $request->post('page', 1);
        $data['size'] = $request->post('size', 10);
        $data['search'] = $request->post('search', []);

        // 服务类型判断
        if (!isset($data['search']['serve_type']) || empty($data['search']['serve_type'])) {
            $data['search']['serve_type'] = $this->serveType;
        }

        $data['enterprise_id'] = $this->enterpriseId;
        $data['store'] =  $this->getStores($this->enterpriseAdmin['role'], false) ?? [];

        $list = (new Task())->taskList($data);

        return success($list);
    }
}