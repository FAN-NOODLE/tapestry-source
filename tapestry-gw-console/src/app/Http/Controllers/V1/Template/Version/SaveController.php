<?php

namespace App\Http\Controllers\V1\Template\Version;

use CrsSdk\Shop\Version;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class SaveController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only([
            'index_template_code',
            'category_template_code',
            'store_index_template_code',
            'my_template_code'
        ]);

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("选择摸板编辑保存");

        (new Version())->save($input);

        return success();
    }
}