<?php
/**
 * 定时清理后台操作日志 ，仅保留近3月
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-28
 * Tome : 15:56
 */

namespace App\Console\Commands;

use App\Models\OperationLog;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoCleanOperationLog extends Command
{
    /**
     * @var string $signature
     */
    protected $signature = 'auto:clean_operation_log';

    /**
     * @var string $description
     */
    protected $description = '自动清理操作日志';


    /**
     * TaskStandard constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 清理三个月前的操作日志
     */
    public function handle()
    {
        $endDate =  strtotime("-0 year -3 month -0 day");
        Log::info('自动删除 ' . date('Y-m-d H:i:s', $endDate). '前日志 start');
        OperationLog::query()
            ->where('created_at', '<=', $endDate)
            ->forceDelete();
        Log::info('自动删除 ' . date('Y-m-d H:i:s', $endDate). '前日志 end');
    }
}
