<?php
/**
 * Created by PhpStorm.
 * User: baizhongyang
 * Date: 2020/8/3
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use CrsSdk\Shop\CouponUser;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;

class CouponListController extends BaseController
{
    /**
     * 客户消费优惠券列表
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params                  = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $params['qy_hash']       = $this->loginAdminEnterprise['hash'];
        $params['type']          = 1; // 状态未使用
        $data                    = (new CouponUser())->couponListFromGuide($params);

        return success($data);
    }
}