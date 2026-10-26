<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 15:05
 */

namespace App\Http\Controllers\V1\Task;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Task;
use Illuminate\Http\Request;
use CrsSdk\Work\Member;

class DetailController extends BaseController
{
    /**
     * 任务详情
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
        // 任务
        $task = (new Task())->detail($data);
        // 任务员工
        $taskMember = (new Task())->taskMember(['ids' => [$data['id']]]);
        $memUid = array_column($taskMember['data'], 'mem_uid');
        // 员工参数列表
        $params = [
            'pageNum' => 1,
            'pageSize' => count($memUid),
            'memUids' => $memUid
        ];
        // 调用sdk
        $condition = [
            'qy_domain' => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params' => $params
        ];

        $list = (new Member())->memberList($condition);

        // 用户列表数据格式化
        $data = $this->format($list['list']);
        // 根据员工uid获取员工姓名
        $task['data']['member'] = $data;

        return success($task);
    }

    /**
     * 格式化员工数据
     *
     * @param array $list 当前员工
     * @return array
     */
    protected function format($list)
    {
        // 初始化待返数据
        $data = [];
        $mem_uid = [];
        $mem_username = [];
        if (!$list) {
            return $data;
        }
        foreach ($list as $member) {
            $mem_uid[] = $member['memUid'];
            $mem_username[] = $member['memUsername'];
        }

        $data['mem_uid'] = $mem_uid;
        // 只展示前五个
        $data['mem_username'] = implode("、", array_slice($mem_username, 0, 5));
        $data['count'] = count($mem_uid);

        return $data;
    }
}