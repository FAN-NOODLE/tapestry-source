<?php
/**
 * 模板新增
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2019/05/30
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\Template\Home;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Template;

class AddController extends BaseController
{

    public function index(Request $request)
    {
        // 接受参数
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("新增首页模板");

        return (new Template())->homeAdd($input);
    }
}
