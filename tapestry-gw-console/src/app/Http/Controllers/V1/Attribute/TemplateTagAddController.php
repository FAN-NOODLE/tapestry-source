<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateTagAddController extends BaseController
{
    /**
     * 规格模板列表
     *
     * @return array
     */
    public function index()
    {
        // 接受参数
        $params = $this->request->only(['tag_name']);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Attribute())->templateTagAdd($params);
        } catch (\Exception $exception) {
            Log::info('新增规格标签模板失败' . $exception->getMessage(), $params);
            notice(1400009);
        }

        return success($list ?? []);
    }

}
