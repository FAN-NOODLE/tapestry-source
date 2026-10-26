<?php

namespace App\Jobs;

use App\Common\MessageConstant;
use App\Jobs\Message\ExternalContactMessage;
use App\Jobs\Message\ExternalContactMessageAdapter;
use App\Models\StoreWorker;
use App\Services\DataStatistic\CustomerStatisticConsumeService;
use Illuminate\Support\Facades\Log;

/**
 * 外部联系人变更顾客统计
 * Class ExternalContactUpdateCustomerStatisticJob
 * @package App\Jobs
 */
class ExternalContactUpdateCustomerStatisticJob extends Job
{

    public $queue = 'external_contact_update_customer_statistic';

    /**
     * 外部联系人消息
     * @var ExternalContactMessage
     */
    private $message;

    /**
     * ExternalContactUpdateCustomerStatisticJob constructor.
     * @param ExternalContactMessage $externalContactMessage
     */
    public function __construct(ExternalContactMessage $externalContactMessage)
    {
        $this->message = $externalContactMessage;
        Log::info("ExternalContactUpdateCustomerStatisticJob product", [
            'message' => $this->message
        ]);
    }



    public function handle()
    {
        Log::info("ExternalContactUpdateCustomerStatisticJob consume", [
            'message' => $this->message,
        ]);

        // 在单方面删除时已经处理过，互删就不做处理了
        if (MessageConstant::EXTERNAL_CONTACT_DELETE == $this->message->getType()
            && $this->message->isUnilateral() == false) {
            Log::info("外部联系人删除消息已经处理过了");
            return;
        }

        $workerList = StoreWorker::query()
            ->where('enterprise_id', $this->message->getEnterpriseId())
            ->where('mem_uid', $this->message->getMemUid())
            ->get()
            ->toArray();

        Log::info("ExternalContactUpdateCustomerStatisticJob workerList", $workerList);
        if (empty($workerList)) {
            return;
        }
        if (!in_array(StoreWorker::IN_SERVICE, array_unique(array_column($workerList, 'worker_status')))) {
            Log::info("ExternalContactUpdateCustomerStatisticJob member is not servicing");
            return;
        }
        $consumeService = new CustomerStatisticConsumeService();

        $consumeService->updateCustomerStatisticByExternalContactMessage($this->message, $workerList);
        $consumeService->updateBindingCustomerStatisticByExternalContactMessage($this->message, $workerList);
    }


    public static function sendCreateMessage(array $memberUserRelation)
    {
//        $message = ExternalContactMessageAdapter::createMessage($memberUserRelation);
//        dispatch(new ExternalContactUpdateCustomerStatisticJob($message));
//        (new ExternalContactUpdateCustomerStatisticJob($message))->handle();
    }

    public static function sendDeleteMessage(array $memberUserRelation)
    {
//        $message = ExternalContactMessageAdapter::deleteMessage($memberUserRelation);
//        dispatch(new ExternalContactUpdateCustomerStatisticJob($message));
//        (new ExternalContactUpdateCustomerStatisticJob($message))->handle();
    }

}
