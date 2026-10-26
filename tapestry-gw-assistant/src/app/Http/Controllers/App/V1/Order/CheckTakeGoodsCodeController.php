<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/11 0011
 * Time: 14:26
 */

namespace App\Http\Controllers\App\V1\Order;

use App\Common\Constant;
use CrsSdk\Shop\NewOrder;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use CrsSdk\Shop\Order;

class CheckTakeGoodsCodeController extends OrderBaseController
{
    /**
     * 订单所对应的收货码的redis key
     */
    const  ORDER_CODE_KEY = "ORDER_TAKE_GOODS_CODE";

    /**
     * 确认用户自提
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index(Request $request,NewOrder $newOrder)
    {
        // 验证输入参数
        $rules = [
            'order_id'        => 'required',
            'store_id'        => 'required',
            'take_goods_code' => 'required'
        ];
        $message = [
            'order_id.required'        => '订单ID不能为空',
            'store_id.required'        => '店铺ID不能为空',
            'take_goods_code.required' => '提货码不能为空'
        ];
        $this->validate($request, $rules, $message);

        // 接收参数
        $id                 = $request->post('order_id');
        $storeId            = $request->post('store_id');
        $inputTakeGoodsCode = $request->post('take_goods_code');

        //获取订单详情中的收货码
        $orderTakesGoodsCode = Redis::get(self::ORDER_CODE_KEY . "_{$id}");
        if (!$orderTakesGoodsCode) {
            $orderDetail         = (new Order())->orderDetail(['order_id' => $id]);

            $orderTakesGoodsCode = $orderDetail['take_goods_code'];
            Redis::setex(self::ORDER_CODE_KEY . "_{$id}",3600, $orderTakesGoodsCode);
        }

        //判断收货码是否一致
        if(strtolower($orderTakesGoodsCode) != strtolower($inputTakeGoodsCode)){
            error(654016);
        }


        $input['id']               = $id;
        $input['event']            = Constant::TAKE_GOODS_IN_STORE;
        $input['enterprise_id']    = $this->enterpriseId;
        $input['operate_store_id'] = $storeId;
        $input['operate_id']       = $this->adminId;
        $input['operate_name']     = $this->mobile;


        try {
            return success($newOrder->changeOrderState($input));
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }

        return success();
    }
}
