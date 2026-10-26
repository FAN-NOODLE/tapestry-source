<?php

namespace App\Http\Controllers\V1\AppCenter;

use CrsSdk\Shop\Cheap;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class CouponSpuListController extends BaseController
{
    /**
     * 未参加活动商品列表接口
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index()
    {
        $input = $this->request->input();

        $this->validate($this->request, [
            'page' => 'required|integer',
            'size' => 'required|integer'
        ]);

        // 请求参数
        $input['enterprise_id'] = $this->enterpriseId;

        // 调用SDK
        $list = (new Cheap())->couponSpuList($input);

        return success($list);
    }
}
