<?php


namespace App\Http\Controllers\V1\ActivityCommon\Luck;

use App\Micros\Activity\LuckMicro;
use App\Http\Controllers\V1\BaseController;

class JoinUserListController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $res = (new LuckMicro())->joinUserList($input);
        return success($res);
    }
}
