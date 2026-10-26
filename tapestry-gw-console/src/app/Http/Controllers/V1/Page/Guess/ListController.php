<?php

namespace App\Http\Controllers\V1\Page\Guess;

use CrsSdk\Shop\PageHot;
use Illuminate\Http\JsonResponse;
use CrsSdk\Shop\PageActivityBanner;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /*
     * 类型：猜你喜欢
     */
    const TYPE_GUESS = 2;

    /**
     * 猜你喜欢 - 列表
     * @return JsonResponse
     */
    public function index()
    {
        // 商品列表
        $list = (new PageHot())->lists([
            'type'          => self::TYPE_GUESS,
            'enterprise_id' => $this->enterpriseId,
        ]);

        // banner图
        $banner = (new PageActivityBanner())->lists([
            'type'          => self::TYPE_GUESS,
            'enterprise_id' => $this->enterpriseId,
        ]);

        return success([
            'list'   => $list,
            'banner' => $banner,
        ]);
    }
}