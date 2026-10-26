<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\App\V1\EnterprisePayment;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;

class RechargeController extends BaseController
{
    /**
     * 充值接口
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // todo...弃用
        return success(['data' => '']);

//        $params = $request->only(['enterprise_id', 'amount', 'voucher_url', 'income_type']);
//
//        (new EnterprisePayment())->recharge($params);
//
//        return success();
    }
}
