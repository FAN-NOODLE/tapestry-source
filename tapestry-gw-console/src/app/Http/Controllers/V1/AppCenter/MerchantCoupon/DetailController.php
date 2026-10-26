<?php
/**
 * 微信优惠卷详情
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-23
 * Time: 17:48
 */

namespace App\Http\Controllers\V1\AppCenter\MerchantCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MerchantCoupon;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input                  = $request->only(['id']);
        $input['enterprise_id'] = $this->enterpriseId;

        $detail = (new MerchantCoupon())->detail($input);

        return success($detail);
    }
}
