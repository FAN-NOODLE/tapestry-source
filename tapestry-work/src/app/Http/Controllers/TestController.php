<?php

namespace App\Http\Controllers;

use App\Jobs\ExternalContactUpdateCustomerStatisticJob;
use App\Jobs\Message\ShareLinkMessage;
use App\Jobs\ShareLinkUpdateCustomerStatisticJob;
use App\Models\MemberUserRelation;
use App\Services\DataStatistic\DailyMemberCustomerStatisticService;
use App\Services\SourceServices;
use CrsSdk\Console\Media;

class TestController extends Controller
{
    /**
     * 测试
     *
     * @throws \Exception
     */
    public function index()
    {
//        $data = (new Media())->list(['type' => '1']);
//
//        (new SourceServices())->sysSourceMediaId($data['data']);

//        $this->testShareLinkMessage();
        $this->testExternal();

    }

    public function testShareLinkMessage()
    {
        $msgStr = <<<MSG
        {
            "type":1,
            "enterpriseId":114,
            "userId":135450491992192,
            "currWorkerId":129069103160768,
            "currWorkerMemUid":"2CD9903A7F000001421861C69453221A",
            "oldWorkerId":null,
            "oldWorkerMemUid":null,
            "createdAt":1612409936,
            "updatedAt":1612409936
        }
MSG;
        $msgArr = json_decode($msgStr, true);

        $message = new ShareLinkMessage();
        $message->setType($msgArr['type']);
        $message->setEnterpriseId($msgArr['enterpriseId']);
        $message->setUserId($msgArr['userId']);
        $message->setCurrWorkerId($msgArr['currWorkerId']);
        $message->setCurrWorkerMemUid($msgArr['currWorkerMemUid']);
        $message->setOldWorkerId($msgArr['oldWorkerId']);
        $message->setOldWorkerMemUid($msgArr['oldWorkerMemUid']);
        $message->setCreatedAt($msgArr['createdAt']);
        $message->setUpdatedAt($msgArr['updatedAt']);

        dispatch(new ShareLinkUpdateCustomerStatisticJob($message));
    }

    public function testExternal()
    {
        $releation = MemberUserRelation::query()->find(2140284)->toArray();

        ExternalContactUpdateCustomerStatisticJob::sendCreateMessage($releation);
//        $this->dispatch(new ExternalContactUpdateCustomerStatisticJob::($releation))
    }
}

