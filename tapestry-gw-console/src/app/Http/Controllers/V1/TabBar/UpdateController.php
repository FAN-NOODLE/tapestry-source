<?php

namespace App\Http\Controllers\V1\TabBar;

use CrsSdk\Shop\TabBar;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class UpdateController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("添加编辑底部按钮");

        (new TabBar())->Save($input);

        return success();
    }
}
