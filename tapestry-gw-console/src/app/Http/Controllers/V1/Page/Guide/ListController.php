<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:00
 */

namespace App\Http\Controllers\V1\Page\Guide;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageGuide;

class ListController extends BaseController
{
    /**
     * 引导区列表
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/29 11:16
     */
    public function index()
    {
        try {
            $list = (new PageGuide())->lists(['enterprise_id' => $this->enterpriseId]);
            return success($list);
        } catch (\Exception $e) {
            error(700201);
        }
    }
}