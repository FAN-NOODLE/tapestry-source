<?php
/**
 * Created by PhpStorm.
 * User: lianghuan
 * Date: 2018/12/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\AppManage;

class ListController extends BaseController
{
    /**
     * 新人优惠列表
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/25 11:16
     * @return
     */
    public function index()
    {
        try {
            $list = (new AppManage())->lists(['enterprise_id' => $this->enterpriseId]);
            return success($list);
        } catch (\Exception $e) {
            error(740001);
        }
    }
}
