<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:20
 */

namespace App\Services;


use App\Models\TaskMemberRecord;

class TaskMemberRecordService
{
    /**
     * 记录模型
     *
     * @var TaskMemberRecord
     */
    protected $taskMemberRecord;

    /**
     * TaskMemberRecordService constructor.
     * @param TaskMemberRecord $taskMemberRecord
     */
    public function __construct(TaskMemberRecord $taskMemberRecord)
    {
        // 记录模型赋值
        $this->taskMemberRecord = $taskMemberRecord;
    }

    /**
     * 礼劵详情
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/30 11:16
     * @param array $where 查询条件
     * @return mixed
     */
    public function detail($where)
    {
        $query = $this->taskMemberRecord->query();
        // 条件组合
        if (!empty($where)) {
            foreach ($where as $key => $value) {
                if(!empty($value)){
                    $query = $query->where($key, $value);
                }
            }
        }

        try {
            return $query->first();
        } catch (\Exception $e) {
            error(110000);
        }
    }
}