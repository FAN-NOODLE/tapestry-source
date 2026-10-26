<?php

namespace App\Http\Controllers\V1\Template\ShopStyle;

use CrsSdk\Shop\Template;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class SetController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only(['style_id']);

        // 参数验证规则
        $rules = ['style_id' => 'required|integer'];

        // 提示语
        $message = ['style_id.*' => '请选择样式'];

        $this->validate($request, $rules, $message);

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("设置企业模板风格");

        (new Template())->shopStyleSet($input);

        return success();
    }
}
