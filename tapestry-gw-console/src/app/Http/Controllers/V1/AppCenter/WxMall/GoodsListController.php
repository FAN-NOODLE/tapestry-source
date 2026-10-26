<?php
/**
 * 商品列表
 * User: 何岳龙
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\WxMall;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxMall;

class GoodsListController extends BaseController
{
    public function index()
    {
        $page = request()->post('page', 1);
        $size = request()->post('size', 10);
        $search = request()->post('search', []);

        return (new WxMall())->goodsList([
            'enterprise_id' => $this->enterpriseId,
            'page' => $page,
            'size' => $size,
            'search' => $search
        ]);
    }
}
