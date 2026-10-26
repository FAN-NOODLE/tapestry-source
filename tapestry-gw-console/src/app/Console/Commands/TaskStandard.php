<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2018/12/4
 * Time: 19:40
 */
namespace App\Console\Commands;

use CrsSdk\Shop\Task;
use Illuminate\Console\Command;

class TaskStandard extends Command
{
    /**
     * @var string
     */
    protected $signature = 'task:task_standard_number';

    /**
     * @var string
     */
    protected $description = '管理后台首页统计';

    /**
     * TaskStandard constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 命令行入口方法
     *
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:05 AM
     */
    public function handle()
    {
        // 销售达标率
        $this->salesTask();
        // 礼券任务数
        $this->couponsTask();
        // 会员招募任务数
        $this->memberTask();
        // 客户标记任务数
        $this->customerTask();
    }

    /**
     * 销售达标率
     */
    protected function salesTask()
    {
        (new Task())->taskStandard();
    }

    /**
     * 礼券任务数
     */
    protected function couponsTask()
    {
        (new Task())->coupon();
    }

    /**
     * 会员招募任务数
     */
    protected function memberTask()
    {
        (new Task())->taskUser();
    }

    /**
     * 客户标记任务数
     */
    protected function customerTask()
    {
        (new \CrsSdk\Work\Task())->memberCount();
    }
}
