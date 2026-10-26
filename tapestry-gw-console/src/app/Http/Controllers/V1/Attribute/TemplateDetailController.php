<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateDetailController extends BaseController
{
    /**
     * 规格模板详情
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'attribute_template_id' => 'required|integer',
        ], [
            'attribute_template_id.*' => '属性模板id不能为空',
        ]);

        // 接受参数
        $params = $this->request->only(['attribute_template_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $template = (new Attribute())->templateDetail($params);
        } catch (\Exception $exception) {
            Log::info('查询规格模板详情失败' . $exception->getMessage(), $params);
            notice(1400003);
        }

        return success($template ?? []);
    }

}
