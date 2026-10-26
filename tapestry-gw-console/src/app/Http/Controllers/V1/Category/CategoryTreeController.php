<?php

namespace App\Http\Controllers\V1\Category;

use CrsSdk\MicroGoodsV1\Category;
use App\Http\Controllers\V1\BaseController;

class CategoryTreeController extends BaseController
{
    /**
     * 商城版
     */
    const SHOP = 1;

    /**
     * 展示空分类
     */
    const SHOW_EMPTY = 1;

    public function categoryTree()
    {
        $category = (new Category())->recursionList([
            'type'       => self::SHOP,
            'source_id'  => $this->enterpriseId,
            'show_empty' => self::SHOW_EMPTY,
        ]);

        return success(['list' => $category['list']]);
    }
}