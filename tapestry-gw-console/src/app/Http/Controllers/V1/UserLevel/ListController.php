<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\UserLevel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\UserLevel;

class ListController extends BaseController
{
    /**
     * 会员等级列表
     *
     * @return array
     */
    public function index()
    {
        // 调用sdk
        $data = (new UserLevel())->userLevel(['enterprise_id' => $this->enterpriseId]);

        return success($data);
    }
}
