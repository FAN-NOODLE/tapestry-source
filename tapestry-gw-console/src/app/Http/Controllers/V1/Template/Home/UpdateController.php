<?php
/**
 * 模板编辑
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2019/05/30
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\Template\Home;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Template;
use Illuminate\Http\Request;

class UpdateController extends BaseController
{

    public function index(Request $request)
    {
        // 接受参数
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("更新首页模板");

        return (new Template())->homeSave($input);
    }
}
