<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:00
 */

namespace App\Http\Controllers\V1\Page\Banner;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageBanner;

class ListController extends BaseController
{
    /**
     * 轮播图列表
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/28 11:16
     */
    public function index()
    {
        try {
            $list = (new PageBanner())->lists(['enterprise_id' => $this->enterpriseId]);
            return success($list);
        } catch (\Exception $e) {
            error(700001);
        }
    }
}