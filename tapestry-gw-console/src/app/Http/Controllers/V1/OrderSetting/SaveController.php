<?php
/**
 * 订单配置保存
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/07/11
 * Time: 18:20
 */
namespace App\Http\Controllers\V1\OrderSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\OrderSetting;

class SaveController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operate_store_id'] = $this->enterpriseId;
        $input['operate_id'] = $this->enterpriseAdmin['id'];
        $input['operate_name'] = $this->enterpriseAdmin['username'];
        $this->writeOperatorLog("后台保存订单配置信息");

        (new OrderSetting())->save($input);

        return success();
    }
}