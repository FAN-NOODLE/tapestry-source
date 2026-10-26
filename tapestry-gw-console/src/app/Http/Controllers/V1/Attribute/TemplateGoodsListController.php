<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;

class TemplateGoodsListController extends BaseController
{
    /**
     * 规格模板商品列表
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'attribute_template_id' => 'required|integer',
            'page'                  => 'integer',
            'size'                  => 'integer',
            'search.spu_id'         => 'integer',
            'search.spu_name'       => 'string',
            'search.source_id'      => 'integer',
        ], [
            'attribute_template_id.*'     => '规格模板Id有误',
            'search.spu_id.integer'       => '商品ID有误',
            'search.spu_name.string'      => '商品名称有误',
            'search.source_id.integer'    => '渠道信息有误',
        ]);
        // 接受参数
        $params = $this->request->only(['attribute_template_id', 'search', 'size', 'page']);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Attribute())->templateGoodsList($params);

            // 匹配门店名称
            $storeList = (new Store())->lists([
                'enterprise_id' => $this->enterpriseId,
                'store_ids'     => array_column($list['list'], 'source_id'),
                'size'          => 10000,
            ])['list'];
            $storeList = array_column($storeList, 'name', 'id');
            foreach ($list['list'] as $key => $value) {
                $list['list'][$key]['store_name'] = $storeList[$value['source_id']] ?? '商城或其他';
            }
        } catch (\Exception $exception) {
            Log::info('查询规格模板商品列表失败' . $exception->getMessage(), $params);
            notice(1400004);
        }

        return success($list ?? []);
    }

}
