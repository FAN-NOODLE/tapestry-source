<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateTagListController extends BaseController
{
    /**
     * 规格模板列表
     *
     * @return array
     */
    public function index()
    {
        // 接受参数
        $params = $this->request->only(['size', 'page']);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Attribute())->templateTagList($params);
        } catch (\Exception $exception) {
            Log::info('查询规格标签模板列表失败' . $exception->getMessage(), $params);
            notice(1400007);
        }

        return success($list ?? []);
    }

}
