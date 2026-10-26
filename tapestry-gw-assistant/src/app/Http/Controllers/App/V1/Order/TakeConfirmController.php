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

class TakeConfirmController extends OrderBaseController
{
    /**
     * 确认用户自提
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index()
    {
        $id = $this->request->get('order_id');
        if (empty($id)) {
            error(654004);
        }
        $storeId = $this->request->get('store_id');

        $input['id']               = $id;
        $input['event']            = Constant::TAKE_GOODS_IN_STORE;
        $input['enterprise_id']    = $this->user['enterprise_id'];
        $input['operate_store_id'] = $storeId;
        $input['operate_id']       = $this->user['admin_id'];
        $input['operate_name']     = $this->user['mobile'];

        try {
            (new NewOrder())->changeOrderState($input);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }

        return success();
    }
}
