<?php
/**
 * 微信优惠卷列表
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-23
 * Time: 17:48
 */
namespace App\Http\Controllers\V1\AppCenter\MerchantCoupon;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MerchantCoupon;

class ListController extends BaseController
{
    public function index()
    {
        // 接受查询参数
        $input = $this->request->only([
            'page',
            'size',
            'search'
        ]);

        $input['search']        = empty($input['search']) ? '' : json_decode($input['search'], true);
        $input['enterprise_id'] = $this->enterpriseId;
        $list = (new MerchantCoupon())->list($input);

        return success($list);
    }
}
