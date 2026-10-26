<?php

namespace App\Http\Controllers\V1\HotWords;

use CrsSdk\MicroGoodsV1\HotWords;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    public function index()
    {
        // 接值
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        return success((new HotWords())->detail($input));
    }
}
