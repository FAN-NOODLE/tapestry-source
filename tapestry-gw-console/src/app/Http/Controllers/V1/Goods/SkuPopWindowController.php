<?php

namespace App\Http\Controllers\V1\Goods;

use CrsSdk\Shop\Template;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SkuPopWindowController extends BaseController
{
    /**
     * 模板和猜你喜欢 SKU弹窗
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function skuPopWindow()
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        $list = (new Template())->skuPopupWindow($params);

        return success($list);
    }
}