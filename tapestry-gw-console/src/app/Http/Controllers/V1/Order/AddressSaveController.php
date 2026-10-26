<?php
/**
 * Created by PhpStorm.
 * User: zq
 * Date: 19-4-8
 * Time: 下午6:20
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class AddressSaveController extends BaseController
{
    /**
     * 更新收货地址
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        return (new Order())->orderAddressSave($input);
    }
}
