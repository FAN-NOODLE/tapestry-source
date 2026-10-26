<?php

namespace App\Http\Controllers\V1\Template\ShopStyle;

use CrsSdk\Shop\Template;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class DeleteController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only(['style_ids']);
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('删除企业模板风格');

        (new Template())->shopStyleDelete($input);

        return success();
    }
}
