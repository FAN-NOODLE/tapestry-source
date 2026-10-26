<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/12/12
 * Time: 18:31
 */

namespace App\Console\Commands;

use App\Services\QyWxMediaService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class UpdateSourceName extends Command
{
    /**
     * @var string
     */
    protected $signature = 'update_source_name';

    /**
     * @var string
     */
    protected $description = '更新无名称素材数据';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 定时任务执行逻辑
     */
    public function handle()
    {
        Log::info('数据同步开始---');
        try {
            (new QyWxMediaService())->updateSourceName();;
        } catch (\Exception $e) {
            Log::info('数据更新失败: '. $e->getMessage());
        }
        echo "Success---数据更新完成";
        return true;
    }
}
