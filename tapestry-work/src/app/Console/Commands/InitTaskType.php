<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\TaskType;
use Illuminate\Console\Command;

class InitTaskType extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'initialize:task_type';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '初始化任务中心-任务类型';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     * @desc 新添加任务类型，请按顺序补充数据至$types尾部
     * @return mixed
     *
     * @throws \Exception
     */
    public function handle()
    {
        $types = [
            [
                'serve_type' => 1,
                'name' => '销售业绩',
                'describe' => '提升导购在商城小程序的销售额',
                'created_at' => time()
            ],
            [
                'serve_type' => 2,
                'name' => '客户管理',
                'describe' => '将导购在微信中的客户好友添加到企业微信中统一管理',
                'created_at' => time()
            ],
            [
                'serve_type' => 1,
                'name' => '生日祝福',
                'describe' => '对当天过生日的客户送生日祝福',
                'created_at' => time()
            ],
            [
                'serve_type' => 1,
                'name' => '礼券到期提醒',
                'describe' => '提醒礼券快到期的客户使用礼券',
                'created_at' => time()
            ],
            [
                'serve_type' => 1,
                'name' => '会员招募',
                'describe' => '邀请非会员客户注册为商城小程序会员',
                'created_at' => time()
            ]
        ];

        $message = '';
        $taskTypeInfo = TaskType::all();

        // 无数据直接入库
        if ($taskTypeInfo->isEmpty()) {
            $message = $this->initTaskType($types);
        } else {
            $alreadyType = $taskTypeInfo->pluck('name')->toArray();

            // 取库中不存在 type数据
            $initType = [];
            foreach ($types as $type) {
                if (!in_array($type['name'], $alreadyType)) {
                    $initType[] = $type;
                }
            }

            if (!empty($initType)) {
                $message = $this->initTaskType($initType);
            }
        }

        $this->info($message);
        return true;
    }

    /**
     * 初始化数据入库
     * @param $data
     * @return string
     */
    protected function initTaskType($data)
    {
        $message = '任务类型数据初始化成功';
        try {
            $res = TaskType::query()->insert($data);
            if ($res) {
                return 'The command ' . $this->signature . ' success:' . $message;
            }
        } catch (\Exception $e) {
            $message = $e->getMessage();
        }
        return 'The command ' . $this->signature . ' error:' . $message;
    }
}
