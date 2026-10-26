<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateListController extends BaseController
{
    /**
     * 规格模板列表
     *
     * @return array
     */
    public function index()
    {
        // 接受参数
        $params = $this->request->only(['tag_id', 'attribute_template_id', 'template_name', 'spu_id', 'spu_name',
                                        'start_time', 'end_time', 'size', 'page']);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Attribute())->templateList($params);
        } catch (\Exception $exception) {
            Log::info('查询规格模板列表失败' . $exception->getMessage(), $params);
            notice(1400005);
        }

        return success($list ?? []);
    }

}
