<?php
/**
 * 首页详情
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/06/25
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\NewTemplate;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\NewTemplate;

class DetailController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only(['id','version']);

        $input['enterprise_id'] = $this->enterpriseId;

        $data = (new NewTemplate())->Detail($input);

        return success($data);
    }
}
