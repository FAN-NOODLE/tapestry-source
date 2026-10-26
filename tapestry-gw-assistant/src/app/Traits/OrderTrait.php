<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/11 0011
 * Time: 14:26
 */

namespace App\Traits;


Trait OrderTrait
{
    /**
     * 检查订单id是否为空并返回订单id
     */
    public function checkOrderId()
    {
        if(! ($order_id = request()->input('order_no'))) {
            error('654004');
        }
        return $order_id;
    }

    /**
     * @param $input 提交到状态机的数据
     */
    public function addUserInfo(array &$input)
    {
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operate_store_id'] = $this->storeId;
        $input['operate_id'] = $this->adminId;
        $input['operate_name'] = $this->adminName;

    }

}