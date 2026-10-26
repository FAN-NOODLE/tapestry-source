<?php
/**
 * 模板编辑
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/06/25
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\NewTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\NewTemplate;
use Illuminate\Http\Request;

class UpdateController extends BaseController
{

    public function index(Request $request)
    {

        // 接受参数
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;
        $input['source_id'] = $this->storeIds;;

        $this->writeOperatorLog("更新首页模板");

        return success((new NewTemplate())->Save($input));
    }
}
