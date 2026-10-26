<?php
/**
 * 模板列表
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2019/05/30
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\Template\Home;

use CrsSdk\Shop\Template;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;

class ListController extends BaseController
{

    public function index(Request $request)
    {
        $input = $request->only(['page', 'size']);
        $input['enterprise_id'] = $this->enterpriseId;

        return (new Template())->homeList($input);
    }
}
