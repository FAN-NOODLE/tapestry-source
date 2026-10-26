<?php


namespace App\Http\Controllers\V1\ActivityCommon\Luck;

use App\Micros\Activity\LuckMicro;
use App\Http\Controllers\V1\BaseController;

class InitRankListController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['type']          = 1;
        $res = (new LuckMicro())->initRankList($input);
        return success($res);
    }
}