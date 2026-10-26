<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Common\Constant;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Facades\Schema;

class DeleteEnterpriseDataTable extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'del:enterprise_trend_data_table';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '删除趋势图表';

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
        $enterpriseList = (new Enterprise())->listNoPage([]);

        foreach ($enterpriseList['data'] as $key => $value) {
            $table = Constant::TABLE_NAME.$value['id'];

            Schema::dropIfExists($table);
        }
    }
}
