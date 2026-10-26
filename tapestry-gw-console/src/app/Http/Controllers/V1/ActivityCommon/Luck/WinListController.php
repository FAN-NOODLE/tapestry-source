<?php


namespace App\Http\Controllers\V1\ActivityCommon\Luck;

use App\Micros\Activity\LuckMicro;
use App\Http\Controllers\V1\BaseController;

class WinListController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['is_console']    = 1;
        $res = (new LuckMicro())->winList($input);
        return success($res);
    }
}
