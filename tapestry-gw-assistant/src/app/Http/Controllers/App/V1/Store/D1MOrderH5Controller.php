<?php
/**
 * Created by PHPStorm.
 * author: baizhongyang
 * Date: 2020/8/20
 * Time: 17:04
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;

class D1MOrderH5Controller extends BaseController
{

    /**
     * 获取我的订单跳转h5地址
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->only(['type', 'openid']);

        $params = [
            'mobile' => $this->mobile,
            'enterprise_id' => $this->enterpriseId,
            'type' => $input['type'],
            'openid' => $input['openid'] ?? ''
        ];
        $res = Common::work()->setApi('service/v1/store/d1m-order-h5-url')->post($params);

        return success($res);

    }

}
