<?php


namespace App\Http\Controllers\App\V1\DistributionGoods;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\MicroGoodsV1\Category;

class CategoryListController extends BaseController
{

    public function index()
    {
        $list = (new Category())->recursionList([
            'source_id' => $this->enterpriseId,
            'show_empty' => 1,  // 是否展示无商品的空分类 0否 1是
            'type' => 1,  // 分类类型 1商城 2门店
        ])['list'] ?? [];

        foreach ($list as &$value) {
            $children = [];
            foreach ($value['children'] as $val) {
                $children = array_merge($children, $val['children'] ?? []);
            }
            $value['children'] = $children;
        }
        unset($value);

        return success($list);
    }
}