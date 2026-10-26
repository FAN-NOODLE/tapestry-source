<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateAddController extends BaseController
{
    /**
     * 新增规格模板
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'template_name' => 'required',
            'tag_id'        => 'required_without:tag_name',
            'tag_name'      => 'required_without:tag_id',
            'attributes'    => 'required|array',
        ], [
            'template_name.*' => '请输入规格模版名称',
            'tag_id.*'        => '请输入标签id或标签名称',
            'tag_name.*'      => '请输入标签id或标签名称',
            'attributes.*'    => '请输入规格信息',
        ]);

        // 接受参数
        $params = $this->request->only([
            'template_name', 'tag_id', 'tag_name', 'parent_id', 'level', 'attributes'
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $template = (new Attribute())->templateAdd($params);
        } catch (\Exception $exception) {
            Log::info('规格模板新增失败' . $exception->getMessage(), $params);
            notice(1400001);
        }

        return success($template ?? []);
    }

}
