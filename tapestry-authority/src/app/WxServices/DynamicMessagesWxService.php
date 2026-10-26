<?php

namespace App\WxServices;

use App\Traits\RequestRecord;
use App\Traits\OpenPlatformTrait;

class DynamicMessagesWxService
{
    use OpenPlatformTrait, RequestRecord;

    public function createActivityId($enterpriseId)
    {
        try {
            $miniProgram = $this->miniProgram($enterpriseId);

            $result = $miniProgram->activity_message->createActivityId();

            $this->wxRecord($result, [
                'url'  => 'cgi-bin/message/wxopen/activityid/create',
                'file' => __FILE__,
                'line' => __LINE__,
                'data' => [
                    'eid'    => $enterpriseId,
                    'result' => $result
                ],
            ]);
        } catch (\Exception $e) {
            $this->wxRecord(['errmsg' => $e->getMessage()], [
                'url'  => 'cgi-bin/message/wxopen/activityid/create',
                'file' => __FILE__,
                'line' => __LINE__,
                'data' => [
                    'eid' => $enterpriseId,
                ],
            ], $e);

            notice(800003, [], ['code:' . $e->getCode(), 'message:' . $e->getMessage()]);
        }

        return success($result ?? []);
    }

    public function updateActivityMessage($enterpriseId, $activityId, $state, $params)
    {
        try {
            $miniProgram = $this->miniProgram($enterpriseId);

            $result = $miniProgram->activity_message->updateMessage(
                $activityId,
                $state,
                $params
            );
            $this->wxRecord($result, [
                'url'  => 'cgi-bin/message/wxopen/updatablemsg/send',
                'file' => __FILE__,
                'line' => __LINE__,
                'data' => [
                    'eid'    => $enterpriseId,
                    'result' => $result
                ],
            ]);

        } catch (\Exception $e) {
            $this->wxRecord(['errmsg' => $e->getMessage()], [
                'url'  => 'cgi-bin/message/wxopen/updatablemsg/send',
                'file' => __FILE__,
                'line' => __LINE__,
                'data' => [
                    'eid' => $enterpriseId,
                ],
            ], $e);

            notice(800004, [], ['code:' . $e->getCode(), 'message:' . $e->getMessage()]);
        }

        return $result ?? [];
    }
}
