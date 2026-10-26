<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-07-16
 * Time: 16:07
 */

namespace App\Http\Controllers\V1\AppCenter\WxMall;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxMall;

class CategoryListController extends BaseController
{
    /**
     * 好物圈物品类目列表接口
     * {domain}/app/v1/app-center/wx-mall/category-list
     *
     * @return array
     */
    public function index()
    {
        // 获取好物圈物品类目列表
        $list = (new WxMall())->wxMallCategoryList();

        return success($list);
    }
}
