<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/13
 * Time: 16:11
 */

namespace App\Http\Controllers\Service\V1\Order;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberOrderStatistics;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SyncController extends BaseController {

    public function index(Request $request) {
        $orderData = $request->post('order_data');
        if (empty($orderData) || empty($orderData['orderInfo'])) {
            error(110007);
        }

        // 只记录待发货状态的订单
        if ($orderData['orderInfo']['status'] != Constant::ORDER_STATUS_WAIT_DELIVER) {
            return success();
        }
        $userId = $orderData['orderUser']['id'];
        $unionId = $orderData['orderUser']['unionId'];
        $enterpriseId = $orderData['orderInfo']['enterpriseId'];
        $createOrderAt = $orderData['orderInfo']['payAt'];

        $userOrderObj = MemberOrderStatistics::getDataByUserId($userId, $enterpriseId);
        if (empty($userOrderObj)) {
            MemberOrderStatistics::createUserOrder($userId, $unionId, 1,
                $orderData['orderInfo']['cashFee'], $orderData['orderInfo']['orderNo'], $enterpriseId, $createOrderAt);
            return success();
        }

        // 判断当前订单数据是否处理过
        if ($userOrderObj->isRecord($orderData['orderInfo']['orderNo'])) {
            return success();
        }
        if ($userOrderObj->addOrder($orderData['orderInfo']['orderNo'], $orderData['orderInfo']['cashFee'], $createOrderAt)) {
            return success();
        }
        error(110000);
    }
}