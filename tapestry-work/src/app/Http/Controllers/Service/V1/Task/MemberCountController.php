<?php
/**
 * 计算员工百分比
 * Created by PhpStorm.
 * User: 朱虹达
 * Date: 2018/11/16
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\Service\V1\Task;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Models\TaskStandard;
use App\Models\User;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\Task;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Redis;

class MemberCountController extends BaseController
{
    /**
     * 员工统计
     */
    const MEMBER_TOTAL = 2;

    public function index()
    {
        // 查询企业列表
        $enterprises = (new Enterprise())->lists(['limit'=>1000]);

        // 循环企业
        foreach ($enterprises['data'] as $v) {
            $total = $number = $task = 0;

            // 查询会员任务
            $list = (new Task())->taskStatistics(['enterprise_id' => $v['id'], 'task_type' => self::MEMBER_TOTAL]);

            if (!empty($list['data'])) {
                // 任务总数
                $total = array_sum(array_column($list['data'], 'task_count'));

                // 循环任务
                foreach ($list['data'] as $key => $value) {
                    $value['enterprise_id'] = $v['id'];
                    // 销售业绩
                    $number += $this->memberTask($value);
                }

                // 百分比
                $task = $total == 0 ? 0 : round($number * 100 / $total);
            }

            // 插入表
            (new \CrsSdk\Shop\Task())->taskAdd(['enterprise_id' => $v['id'], 'type' => self::MEMBER_TOTAL, 'task_number' => $task]);
        }

        return success([
            'data' => []
        ]);
    }


    /**
     * 查询每个人完成多少工作量
     * @param $data
     */
    protected function memberTask($data)
    {
        // 根据时间段计算出不同人的订单
        $count = MemberUserRelation::query()
            ->where('created_at', '>=',$data['start_time'])
            ->where('created_at','<=', $data['end_time'])
            ->whereIn('mem_uid',$data['task_member'])
            ->where(['enterprise_id' => $data['enterprise_id']])
            ->count();
        
        return $count;
    }


    // 打印sql
    function getSql($query)
    {
        $sql = str_replace("?", "'%s'", $query->toSql());
        return vsprintf($sql, $query->getBindings());
    }

}
