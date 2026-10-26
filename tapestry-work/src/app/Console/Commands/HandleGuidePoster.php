<?php

namespace App\Console\Commands;

use App\Jobs\StoreWorkerGenerateQrCodeJob;
use App\Models\Scrm\Members;
use CrsSdk\Boss\Enterprise;
use Exception;
use Illuminate\Console\Command;

class HandleGuidePoster extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'handle_guide_poster';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '手动处理海报为空的导购,生成海报';


    public function __construct()
    {
        parent::__construct();
    }

    /**
     * @return bool
     * @throws Exception
     */
    public function handle(): bool
    {
        // 设置内存限制为无限制
        ini_set('memory_limit', '-1');
        
        // 查询企业列表
        $enterprises = (new Enterprise())->lists()['data'];

        foreach ($enterprises as $enterprise) {
            $members = Members::query()
                ->where('status', 1)
                ->where('enterprise_id', $enterprise['id'])
                ->where('qr_poster_code', '')
                ->get();
            if ($members->count() == 0) {
                continue;
            }
            $members = $members->toArray();

            $domain       = $enterprise['domain'];
            $enterpriseId = $enterprise['id'];

            $workerIds = array_column($members, 'id');

            // 异步处理生成导购海报
            $data = [
                'worker_ids'    => $workerIds,
                'enterprise_id' => $enterpriseId,
                'store_id'      => 0,
                'domain'        => $domain,
                'operator_id'   => 0,
                'operator_name' => '系统',
            ];

            dispatch(new StoreWorkerGenerateQrCodeJob($data));
        }

        return true;
    }
}
