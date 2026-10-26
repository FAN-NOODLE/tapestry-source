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

class DeleteController extends BaseController
{
    /**
     * 删除任务
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $data = $request->post();
        // 必填字段验证
        if (!isset($data['id']) && empty($data['id'])) {
            error(600003);
        }

        // 查询任务详情
        $data = (new Task())->detail($data);

        if (isset($data['data']['task_name'])) {
            $this->writeOperatorLog('删除任务:'.$data['data']['task_name']);
        }

        // 删除
        (new Task())->delete($data['data']);

        return success([
            'data' => [],
            'meta' => []
        ]);
    }
}
