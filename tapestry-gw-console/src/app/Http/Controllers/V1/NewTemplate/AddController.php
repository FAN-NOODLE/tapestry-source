<?php
/**
 * 模板新增
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/06/25
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\NewTemplate;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\NewTemplate;

class AddController extends BaseController
{
    public function index(Request $request)
    {

        // 接受参数
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;
        $input['source_id'] = $this->storeIds;


        $this->writeOperatorLog("新增首页模板");

        (new NewTemplate())->Add($input);

        return success();
    }
}
