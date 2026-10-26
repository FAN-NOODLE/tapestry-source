<?php

namespace App\Http\Controllers\V1\OrderSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\OrderSetting;
use Illuminate\Http\JsonResponse;
use Exception;

class DetailsController extends BaseController
{
    /**
     * 订单配置详情
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $params['enterprise_id'] = $this->enterpriseId;
        // 获取配置信息
        $data = (new OrderSetting())->details($params);

        return success($data);
    }

}
