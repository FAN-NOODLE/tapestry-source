<?php

namespace App\Http\Controllers\V1\Template\ShopStyle;

use CrsSdk\Shop\Template;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class InfoController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->only(['version']);

        $input['enterprise_id'] = $this->enterpriseId;

        $info = (new Template())->shopStyleInfo($input)['data'];

        return success($info);
    }
}
