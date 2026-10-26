<?php

namespace App\Jobs;

use App\Services\DataStatistic\CustomerStatisticConsumeService;
use Illuminate\Support\Facades\Log;

/**
 * 分享链路变更，更新顾客统计数据
 * Class ShareLinkBatchUpdateCustomerStatisticJob
 * @package App\Jobs
 */
class ShareLinkBatchUpdateCustomerStatisticJob extends Job
{

    public $queue = 'share_link_batch_update_customer_statistic';

    /**
     * @var int
     */
    private $type;

    /**
     * @var array
     */
    private $messageList;

    /**
     * ShareLinkBatchUpdateCustomerStatisticJob constructor.
     * @param int $type
     * @param array $messageList
     */
    public function __construct(int $type, array $messageList)
    {
        $this->type = $type;
        $this->messageList = $messageList;
        Log::info("ShareLinkBatchUpdateCustomerStatisticJob product", [
            'type' => $this->type,
            'messageList' => $this->messageList
        ]);
    }

    public function handle()
    {
        Log::info("ShareLinkBatchUpdateCustomerStatisticJob consume", [
            'type' => $this->type,
            'messageList' => $this->messageList,
        ]);
        (new CustomerStatisticConsumeService())->updateBindingCustomerStatisticByShareLinkBatchMessage($this->type,
            $this->messageList);
    }

}
