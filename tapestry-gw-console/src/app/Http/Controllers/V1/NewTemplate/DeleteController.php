<?php
/**
 * 删除模板
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/06/25
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\NewTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\NewTemplate;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    public function index(Request $request)
    {

        // 接受参数
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;

        (new NewTemplate())->Delete($input);

        return success();
    }
}
