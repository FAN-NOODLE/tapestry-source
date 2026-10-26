<?php

namespace App\Jobs;

use App\Jobs\Message\ShareLinkMessage;
use App\Services\DataStatistic\CustomerStatisticConsumeService;
use Illuminate\Support\Facades\Log;

/**
 * 分享链路变更，更新顾客统计数据
 * Class ShareLinkUpdateCustomerStatisticJob
 * @package App\Jobs
 */
class ShareLinkUpdateCustomerStatisticJob extends Job
{

    public $queue = 'share_link_update_customer_statistic';

    /**
     * @var ShareLinkMessage
     */
    private $message;

    /**
     * ShareLinkUpdateCustomerStatisticJob constructor.
     * @param ShareLinkMessage $message
     */
    public function __construct(ShareLinkMessage $message)
    {
        $this->message = $message;
        Log::info("ShareLinkUpdateCustomerStatisticJob product", [
            'message' => $this->message
        ]);
    }

    public function handle()
    {
        Log::info("ShareLinkUpdateCustomerStatisticJob consume", [
            'message' => $this->message,
        ]);
        (new CustomerStatisticConsumeService())->updateBindingCustomerStatisticByShareLinkMessage($this->message);
    }
}
