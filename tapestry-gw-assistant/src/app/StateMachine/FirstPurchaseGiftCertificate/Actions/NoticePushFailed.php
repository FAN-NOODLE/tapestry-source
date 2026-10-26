<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/8/13
 * Time: 23:10
 */

namespace App\StateMachine\FirstPurchaseGiftCertificate\Actions;

use App\StateMachine\Action;
use App\Traits\GiftCertificateMachineTools;

class NoticePushFailed implements Action
{
    use GiftCertificateMachineTools;

    public function handle($from, $event, $to, $data, $stateObj, $eventDesc)
    {
        //todo 执行业务逻辑
        $objectId = $data['id'];
        $update  = [
            'status' => $to,
            'result_response_time' => time(),
            'failed_message'=>$data['failed_message']
        ];
        $this->updateGiftCertificate($objectId, $update);
        $this->recordGiftCertificateTrace($objectId, $from, $to, $event, $eventDesc,$update);

        $res_msg = [
            'data' => [],
            'meta' => []
        ];
        return $res_msg;
    }
}