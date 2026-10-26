<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/23
 * Time: 18:25
 */

namespace App\Models;

class Task extends Base
{
    /*
     * 任务类型：1:销售业绩，2:客户管理，3:生日祝福，4:礼券到期提醒，5:会员招募
     */
    const TASK_TYPE = [
        1 => '销售业绩',
        2 => '客户管理',
        3 => '生日祝福',
        4 => '礼券到期提醒',
        5 => '会员招募',
    ];

    /**
     * 销售业绩
     */
    const  SALES_PERFORMANCE = 1;

    /**
     * 客户管理
     */
    const CUSTOMER_MANAGEMENT = 2;

    /**
     * 生日祝福
     */
    const BIRTHDAY_WISHES = 3;

    /**
     * 礼券到期提醒
     */
    const COUPONS_EXPIRATION = 4;

    /**
     * 会员招募
     */
    const MEMBER_RECRUITMENT = 5;

    /**
     * 礼劵提醒状态
     */
    const MEMBER_COUPON = 100;

    /*
     * 任务周期：(1:一天，2:一周，3:一个月，4:三个月，5:六个月，6:一年)
     */
    const TASK_CYCLE = [
        1 => '一天',
        2 => '一周',
        3 => '一个月',
        4 => '三个月',
        5 => '六个月',
        6 => '一年',
    ];

    /**
     * 一天
     */
    const ONE_DAY_TYPE = 1;

    /**
     * 一周
     */
    const ONE_WEEK_TYPE = 2;

    /**
     * 一个月
     */
    const ONE_MONTHS_TYPE = 3;

    /**
     * 三个月
     */
    const THREE_MONTHS_TYPE = 4;

    /**
     * 六个月
     */
    const SIX_MONTHS_TYPE = 5;

    /**
     * 一年
     */
    const ONE_YEAR_TYPE = 6;

    /**
     * 任务相关状态
     */
    const TASK_STATUS_START = 1;
    const TASK_STATUS_HAND = 2;
    const TASK_STATUS_END = 3;

    /**
     * 是否循环： （1:循环，2：不循环）
     */
    const TASK_LOOP_STATUS = 1;
    const TASK_NO_CYCLE_STATUS = 2;

    /**
     * 提醒类型： （1:开始，2：结束）
     */
    const REMIND_TYPE_START = 1;
    const REMIND_TYPE_END = 2;

    /**
     * 是否提醒： （1:提醒）
     */
    const IS_REMIND = 1;

    /**
     * 表名
     * @var string
     */
    protected $table = 'task';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'task_name',
        'task_type',
        'serve_type',
        'store_id',
        'start_time',
        'end_time',
        'task_cycle',
        'is_loop',
        'task_index',
        'is_start_task',
        'is_end_task',
        'task_start_time',
        'task_end_time'
    ];

    public function taskMemberRelations()
    {
        return $this->hasMany(TaskMemberRelation::class, 'task_id', 'id');
    }
}
