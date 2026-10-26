<?php

namespace App\Http\Controllers\V1\TabBar;

use CrsSdk\Shop\TabBar;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->only(['id']);

        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new TabBar())->Detail($input)['data'];

        return success(['list' => $list]);
    }
}
