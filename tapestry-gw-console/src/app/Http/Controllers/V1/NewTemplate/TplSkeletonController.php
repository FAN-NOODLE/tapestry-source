<?php
/**
 * 获取模板骨架
 * Created by PhpStorm.
 * User: 侯英才
 * Date: 2020/03/20
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\NewTemplate;

use App\Http\Controllers\V1\BaseController;

class TplSkeletonController extends BaseController
{
    public function index()
    {
        $num = $this->request->input('num', 0);
        if (empty($num)) {
            error(453005);
        }

        $data = config('store_template');

        return success([$data[$num - 1]]);
    }
}
