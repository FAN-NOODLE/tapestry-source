<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-07-31
 * Time: 13:16
 */

namespace Tests\App\V1\OrderStatusNotice;

use App\Events\StoreOrderStatus;

class OrderStatusNoticeTest extends \TestCase
{

    private $queueData = '{
    "meta": {
        "source": "crs"
    },
    "orderInfo": {
        "id": 102174798588608,
        "orderNo": "1102174798224064",
        "enterpriseId": 99045793733312,
        "status": 20,
        "type": 14,
        "totalAmount": 1,
        "preferentialPrice": 0,
        "activityDiscountAmount": 0,
        "finalAmount": 1,
        "paymentChannel": 1,
        "payAt": 0,
        "refundAt": 0,
        "createdAt": 1596162038,
        "deliveryType": 1,
        "deliveryName": "王牌",
        "deliveryMobile": "13467154433",
        "deliveryAddress": "是可以",
        "remarks": "",
        "thirdNo": "0",
        "sendAt": 0,
        "signAt": 0,
        "cashFee": 0,
        "couponFee": 0,
        "otherFee": 0,
        "couponCount": 0,
        "province": "山西省",
        "city": "晋城市",
        "area": "泽州县",
        "logisticIdentity": "",
        "logisticNo": "",
        "isLottery": 0,
        "logisticsPrice": 0,
        "lotteryId": 0,
        "storeId": 99045793733312,
        "storeStaffId": 0,
        "serviceStoreId": 99045793733312,
        "serviceStaffId": 0,
        "totalRefundFee": 0,
        "goodsCount": 1
    },
    "orderRefund": {
        "refundApplyId": 0,
        "refundStatus": 0,
        "serviceStoreId": 0,
        "refundApplyType": 0,
        "refundFee": 0,
        "currentRefundFee": 0,
        "currentOrderStatus": 0,
        "refundReasonId": 0,
        "reason": "",
        "orderGoodsStatus": 0,
        "isManual": 0,
        "refundThirdId": "",
        "outRefundNo": "",
        "createdAt": 0,
        "updatedAt": 0,
        "logisticIdentity": "",
        "logisticsNo": ""
    },
    "action": "send_message"
}';

    public function test_run()
    {
        $mockNotice = json_encode($this->queueData, true);

        event(new StoreOrderStatus($mockNotice, microtime(true)));
    }
}
