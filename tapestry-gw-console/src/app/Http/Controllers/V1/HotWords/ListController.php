<?php

namespace App\Http\Controllers\V1\HotWords;

use CrsSdk\MicroGoodsV1\HotWords;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        return success((new HotWords())->lists($input));
    }
}
