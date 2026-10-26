<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateUpdateController extends BaseController
{
    /**
     * 编辑规格模板
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'attribute_template_id'  => 'required|integer',
            'attributes'             => 'array',
            'del_attribute_id'       => 'array',
            'del_attribute_item_ids' => 'array',
            'template_name'          => 'sometimes',
        ], [
            'attribute_template_id.*' => '属性模板id不能为空',
        ]);

        // 接受参数
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;

        if (($params['template_name'] ?? '') && empty(trim($params['template_name']))) {
            notice(100000, '请输入规格模版名称');
        }

        try {
            (new Attribute())->templateUpdate($params);
        } catch (\Exception $exception) {
            Log::info('规格模板编辑失败' . $exception->getMessage(), $params);
            notice(1400008);
        }

        return success();
    }
}
