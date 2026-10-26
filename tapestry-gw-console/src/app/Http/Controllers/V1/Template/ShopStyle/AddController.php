<?php

namespace App\Http\Controllers\V1\Template\ShopStyle;

use CrsSdk\Shop\Template;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    public function index(Request $request)
    {
        $input                  = $request->only(['styles', 'rgb_colors']);
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('添加企业模板风格');

        (new Template())->shopStyleAdd($input);

        return success();
    }
}
