<?php
/**
 * 首页详情
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2019/05/30
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\Template\Home;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Template;

class DetailController extends BaseController
{

    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only(['page', 'size', 'id']);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['template_id'] = $input['id'];

        return (new Template())->homeDetail($input);
    }
}
