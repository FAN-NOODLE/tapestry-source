<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:00
 */

namespace App\Http\Controllers\V1\Page\Hot;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageActivityBanner;
use CrsSdk\Shop\PageHot;

class ListController extends BaseController
{

    const TYPE_HOT=0;
    /**
     * 爆品推荐列表
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/29 11:16
     * @return array
     */
    public function index()
    {
        try {
            $list = (new PageHot())->lists(['enterprise_id' => $this->enterpriseId]);
            $banner = (new PageActivityBanner())->lists(['enterprise_id' => $this->enterpriseId,'type'=>self::TYPE_HOT]);
            return success(['data'=>$list['data'],'banner'=>empty($banner['data']) ? '' : $banner['data']]);
        } catch (\Exception $e) {
            error(700101);
        }
    }
}