<?php
/**
 * 订单状态变更
 *
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-07-15
 * Time: 18:32
 */

namespace App\Listeners;

use App\Common\Constant;
use App\Events\OrderServiceNotice;
use App\Events\StoreOrderStatus;
use App\Library\Message\FactoryMessage;
use CrsSdk\Shop\Order;
use CrsSdk\Shop\Store;
use CrsSdk\Work\StoreWorker;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Support\Facades\Log;

class ServiceOrderListener
{
    public function handle(StoreOrderStatus $storeOrderStatus)
    {
        $orderInfo = $storeOrderStatus->getQueueData()['orderInfo'];
        // 订单编号
        $orderNo = $orderInfo['order_no'];
        // 订单状态
        $status = $orderInfo['status'];
        Log::error("[ServiceOrderListener - {$orderNo}]推送短信开始");
        try {
            // 查询门店信息
            /** @var $store Store **/
            $store = app(Store::class);
            $storeDetail = $store->detail([
                'enterprise_id' => $orderInfo['enterpriseId'],
                'id' => $orderInfo['serviceStoreId']
            ]);
            $workers    = (new StoreWorker())->getList([
                'enterprise_id' => $orderInfo['enterpriseId'],
                'store_id'      => $orderInfo['service_store_id'],
                'is_page'       => false
            ]);
            $workers = array_column($workers['list'], null, 'id');
            $smsMessage = [
                'templateCode' => Constant::SMS_TEMPLATE['service_store'],
                'data'         => [
                    'storeName'   => $storeDetail['short_name'],
                    'orderStatus' => Constant::ORDER_STATUS_TEXT[$status],
                    'orderNo'     => $orderNo
                ],
            ];
            // 发送短息
            $smsObject = (new FactoryMessage())->getMessageHandlers('easySms');
            // 还没有人接单，给门店下所有的人员发信息
            if ($orderInfo['serviceStaffId'] == 0 && $status == 40) {
                foreach ($workers['list'] as $worker) {
                    if ($worker['worker_status'] == Constant::WORKER_OFFLINE) {
                        continue;
                    }

                    if ($worker['is_send_sms'] == 2) {
                        continue;
                    }

                    $smsMessage['subject'] = $worker['mobile'];
                    if (empty($smsMessage['subject'])) {
                        $smsMessage['subject'] = $worker['sms_mobile'];
                    }

                    $smsObject->send($smsMessage);
                    Log::info("[ServiceOrderListener]导购{$worker['name']}发送短信成功", $smsMessage);
                }
            } else if (in_array($status, Constant::SERVICE_ORDER_NOTICE_STATUS) && !empty($orderInfo['serviceStaffId'])) {
                foreach ($workers as $worker) {
                    if ($worker['id'] != $orderInfo['serviceStaffId'] && $worker['is_boss'] != Constant::IS_BOSS_YES) {
                        continue;
                    }

                    if ($worker['worker_status'] == Constant::WORKER_OFFLINE) {
                        Log::error("[ServiceOrderListener]导购{$worker['name']}已离职");
                        return;
                    }

                    if ($worker['is_send_sms'] == 2) {
                        Log::error("[ServiceOrderListener]导购{$worker['name']}拒收短信");
                        return;
                    }
                    $smsMessage['subject'] = $worker['sms_mobile'];
                    if (empty($smsMessage['subject'])) {
                        $smsMessage['subject'] = $worker['mobile'];
                    }
                    $smsObject->send($smsMessage);
                    Log::info("[ServiceOrderListener]导购{$worker['name']}发送短信成功", $smsMessage);
                }
            }
            return;
        } catch (\Exception $exception) {
            Log::error("[ServiceOrderListener - {$orderNo}]推送短信失败，错误信息{$exception->getMessage()}");
        }
    }
}
