<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-07-15
 * Time: 18:32
 */

namespace App\Listeners;

use App\Common\Constant;
use App\Events\OrderServiceNotice;
use App\Library\Message\FactoryMessage;
use CrsSdk\Shop\Order;
use CrsSdk\Shop\Store;
use CrsSdk\Work\StoreWorker;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Support\Facades\Log;

class ServiceOrderNoticeListener implements ShouldQueue
{
    public function handle(OrderServiceNotice $orderServiceNotice)
    {
        $orderNo        = $orderServiceNotice->getOrderNo();
        $serviceStaffId = $orderServiceNotice->getWorkerId();
        $enterpriseId   = $orderServiceNotice->getEnterpriseId();
        Log::error("[ServiceOrderNoticeListener - {$orderNo}]推送短信开始");
        try {
            // 存在服务店员
            if (!empty($serviceStaffId)) {
                $storeWork = new StoreWorker();
                // 查询导购信息
                $serviceStaffDetail = $storeWork->detail([
                    'enterprise_id' => $enterpriseId,
                    'worker_id'     => $serviceStaffId
                ]);

                if (empty($serviceStaffDetail)) {
                    Log::error("[ServiceOrderNoticeListener]没有查询导购信息");
                    return;
                }

                // 导购是否在职
                if (Constant::WORKER_OFFLINE == $serviceStaffDetail['worker_status']) {
                    Log::error("[ServiceOrderNoticeListener]导购{$serviceStaffDetail['name']}已离职");
                    return;
                }
                // 是否开启短信接收
                if (2 == $serviceStaffDetail['is_send_sms']) {
                    Log::error("[ServiceOrderNoticeListener]导购{$serviceStaffDetail['name']}已拒绝接收短信通知");
                    return;
                }
                // 是否门店导购
                if (in_array($serviceStaffDetail['type'], [Constant::MEMBER_GUIDE_TYPE_COMMUNITY, Constant::MEMBER_GUIDE_TYPE_WORK])) {
                    Log::error("[ServiceOrderNoticeListener]导购{$serviceStaffDetail['name']}非门店导购，不接收服务订单通知");
                    return;
                }
                // 查询订单详情
                $orderDetail = (new Order())->orderDetail([
                    'order_no'      => $orderNo,
                    'enterprise_id' => $enterpriseId
                ]);

                $smsMessageData['subject']             = $serviceStaffDetail['mobile'];
                $smsMessageData['data']['orderStatus'] = Constant::ORDER_STATUS_TEXT[$orderDetail['status']];
                $smsMessageData['data']['orderNo']     = $orderNo;
                $smsMessageData['templateCode'] = Constant::SMS_TEMPLATE['service_store'];
                if (!empty($serviceStaffDetail['sms_mobile'])) {
                    $smsMessageData['subject'] = $serviceStaffDetail['sms_mobile'];
                }
                // 查询门店信息
                /** @var $store Store * */
                $store                               = app(Store::class);
                $storeData                           = $store->detail([
                    'enterprise_id' => $enterpriseId,
                    'id'            => $serviceStaffDetail['store_id']
                ]);
                $smsMessageData['data']['storeName'] = $storeData['short_name'];
                // 发送短信
                $smsObject = (new FactoryMessage())->getMessageHandlers('easySms');
                $smsObject->send($smsMessageData);
                Log::info("[ServiceOrderNoticeListener-{$orderNo}]给服务人发送短信成功", $smsMessageData);
                // 如果当前导购非店长，需要给他的店长推送消息
                if ($serviceStaffDetail['is_boss'] == 2) {
                    $serviceStoreManager = $storeWork->getList([
                        'enterprise_id' => $enterpriseId,
                        'store_id'      => $serviceStaffDetail['store_id'],
                        'is_boss'       => 1,
                        'page'          => 1,
                        'size'          => 1
                    ])['list'];

                    if (!empty($serviceStoreManager)) {
                        $serviceStoreManager = current($serviceStoreManager);
                        if (2 == $serviceStoreManager['is_send_sms']) {
                            Log::error("[ServiceOrderNoticeListener]导购的店长{$serviceStoreManager['name']}已拒收短信");
                            return;
                        }

                        $smsMessageData['subject'] = $serviceStoreManager['mobile'];
                        if (!empty($serviceStoreManager['sms_mobile'])) {
                            $smsMessageData['subject'] = $serviceStoreManager['sms_mobile'];
                        }
                        $smsObject->send($smsMessageData);
                        Log::info("[ServiceOrderNoticeListener-{$orderNo}]给服务人店长发送短信成功", $smsMessageData);
                    }
                }
            }
            return;
        } catch (\Exception $exception) {
            Log::error("[ServiceOrderNoticeListener - {$orderNo}]推送短信失败，错误信息{$exception->getMessage()}");
        }
    }
}
