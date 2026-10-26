<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use Illuminate\Support\Facades\Log;

class TemplateGroupByTagListController extends BaseController
{
    /**
     * 规格模板列表
     *
     * @return array
     */
    public function index()
    {
        // 接受参数
        $params = $this->request->only(['template_or_tag']);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['size'] = Constant::MAX_DEFAULT_SIZE;

        try {
            $list = (new Attribute())->templateList($params)['list'];
            $result = [];
            if (count($list) > 0) {
                $data = collect($list)->groupBy('tag_name')->toArray();
                foreach ($data as $key=>$tagItem) {
                    $temp['tag_name'] = $key;
                    $temp['template_list'] = [];
                    foreach ($tagItem as $templateItem) {
                        $temp['template_list'][] = $templateItem;
                    }
                    $result[] = $temp;
                }
            }
        } catch (\Exception $exception) {
            Log::info('查询规格模板列表失败' . $exception->getMessage(), $params);
            notice(1400005);
        }

        return success($result ?? []);
    }

}
