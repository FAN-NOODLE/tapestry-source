<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateMoveController extends BaseController
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
            'attribute_template_id' => 'required|integer',
            'tag_id'                => 'required_without:tag_name',
            'tag_name'              => 'required_without:tag_id',
        ], [
            'attribute_template_id.*' => '属性模板id不能为空',
            'tag_id.*'                => '标签id和标签名称不能为空',
            'tag_name.*'              => '标签id和标签名称不能为空',
        ]);

        // 接受参数
        $params = $this->request->only(['attribute_template_id', 'tag_id', 'tag_name']);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            (new Attribute())->templateMove($params);
        } catch (\Exception $exception) {
            Log::info('规格模板移动失败' . $exception->getMessage(), $params);
            notice(1400006);
        }

        return success();
    }
}
