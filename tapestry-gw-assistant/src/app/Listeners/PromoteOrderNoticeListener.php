<?php

namespace App\Listeners;

use App\Common\Constant;
use App\Events\StoreOrderStatus;
use App\Library\Message\FactoryMessage;
use CrsSdk\Shop\Store;
use CrsSdk\Work\StoreWorker;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class PromoteOrderNoticeListener
{
    use InteractsWithQueue;
    /**
     *
     *
     * @param StoreOrderStatus $storeOrderStatus
     * @return void
     */
    public function handle(StoreOrderStatus $storeOrderStatus)
    {
        // 获取的数据是
        $queueDataOfOrder = $storeOrderStatus->getQueueData()['orderInfo'];
        $queueOrderStatusTrace = $storeOrderStatus->getQueueData()['orderStatusTrace'];
        // 订单状态
        $status = $queueDataOfOrder['status'];
        // 订单编号
        $orderNo = $queueDataOfOrder['orderNo'];
        // 订单所属推广导购ID
        $storeWorkerId = $queueDataOfOrder['storeStaffId'];
        Log::error("[PromoteOrderNoticeListener - {$orderNo}-推广订单发短信开始]");
        $jsonDecode=json_encode($storeOrderStatus->getQueueData());
        Log::error("[PromoteOrderNoticeListener - {$orderNo}-jsonDecode参数{$jsonDecode}]");
        Log::error("[PromoteOrderNoticeListener - {$orderNo}]推送短信开始");

        try {

            if (!in_array($status, Constant::PROMOTE_ORDER_NOTICE_STATUS)) {
                Log::info("[PromoteOrderNoticeListener - {$orderNo}]订单状态({$status})不需要推送消息");
                return;
            }

            if (!in_array(intval($queueOrderStatusTrace['from']),Constant::PROMOTE_ORDER_STATUS_TRACE)) {
                Log::info("[PromoteOrderNoticeListener - {$orderNo}]流转状态不是发消息订单");
                return;
            }
            $storeWork = new StoreWorker();
            $store            = app(Store::class);
            $smsMessageObject = (new FactoryMessage())->getMessageHandlers('easySms');
            // 推广订单消息
            if (!empty($storeWorkerId)) {
                // 查询推广人是否是导购
                $promoteStaffDetail = $storeWork->detail([
                    'worker_id'     => $storeWorkerId,
                    'enterprise_id' => $queueDataOfOrder['enterpriseId']
                ]);
                // 推广人是否在职
                if ($promoteStaffDetail['worker_status'] == 2) {
                    Log::error("[PromoteOrderNoticeListener]导购{$promoteStaffDetail['name']}已经离职");
                    return;
                }
                // 是否可以接收短信
                if ($promoteStaffDetail['is_send_sms'] == 2) {
                    Log::error("[PromoteOrderNoticeListener]导购{$promoteStaffDetail['name']}已拒收短信提醒");
                    return;
                }
                // 短信参数
                $smsMessageData['subject'] = $promoteStaffDetail['mobile'];
                // 用户设置了接收短信手机号
                if (!empty($promoteStaffDetail['sms_mobile'])) {
                    $smsMessageData['subject'] = $promoteStaffDetail['sms_mobile'];
                }
                $smsMessageData['data']['orderNo'] = $orderNo;
                switch ($promoteStaffDetail['type']) {
                    case Constant::MEMBER_GUIDE_TYPE_STORE:
                        $storeDetail                         = $store->detail([
                            'enterprise_id' => $queueDataOfOrder['enterpriseId'],
                            'id'            => $promoteStaffDetail['store_id']
                        ]);
                        $smsMessageData['templateCode']      = Constant::SMS_TEMPLATE['promote_store'];
                        $smsMessageData['data']['storeName'] = $storeDetail['short_name'];
                        break;
                    case Constant::MEMBER_GUIDE_TYPE_COMMUNITY:
                    case Constant::MEMBER_GUIDE_TYPE_WORK:
                        $smsMessageData['templateCode'] = Constant::SMS_TEMPLATE['promote_work'];
                        break;
                }

                // 给推广人发送消息
                $smsMessageObject->send($smsMessageData);
                Log::info("[PromoteOrderNoticeListener-{$orderNo}]给推广人发送短信成功", $smsMessageData);

                // 如果此导购是门店导购且非店长，给他的店长推送短信提醒
                if (Constant::MEMBER_GUIDE_TYPE_STORE == $promoteStaffDetail['type'] && $promoteStaffDetail['is_boss'] != 1) {
                    $storeManager = $storeWork->getList([
                        'enterprise_id' => $queueDataOfOrder['enterpriseId'],
                        'store_id'      => $promoteStaffDetail['store_id'],
                        'is_boss'       => 1,
                        'page'          => 1,
                        'size'          => 1
                    ])['list'];

                    if (!empty($storeManager)) {
                        $storeManager = current($storeManager);
                        if (2 == $storeManager['is_send_sms']) {
                            Log::error("[OrderChangeNoticeListener]导购的店长{$storeManager['name']}已拒收短信");
                            return;
                        }

                        $smsMessageData['subject'] = $storeManager['mobile'];
                        if (!empty($storeManager['sms_mobile'])) {
                            $smsMessageData['subject'] = $storeManager['sms_mobile'];
                        }

                        $smsMessageObject->send($smsMessageData);
                        Log::info("[PromoteOrderNoticeListener-{$orderNo}]给推广人店长发送短信成功", $smsMessageData);
                    }
                }
            }
            return;
        } catch (\Exception $exception) {
            Log::error("[PromoteOrderNoticeListener - {$orderNo}]推送短信失败，错误信息{$exception->getMessage()}");
        }
    }
}
