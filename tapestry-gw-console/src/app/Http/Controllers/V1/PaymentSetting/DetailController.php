<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\PaymentSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Payment;

class DetailController extends BaseController
{
    /**
     * @return array
     */
    public function index()
    {
        $data = (new Payment())->detail(['enterprise_id'=>$this->enterpriseId]);

        return success($data['data']);
    }
}
