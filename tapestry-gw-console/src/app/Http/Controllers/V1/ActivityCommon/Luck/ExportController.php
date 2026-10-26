<?php

namespace App\Http\Controllers\V1\ActivityCommon\Luck;

use App\Micros\Activity\LuckMicro;
use App\Http\Controllers\V1\BaseController;

class ExportController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['admin_name'] = $this->enterpriseAdmin['username'];
        $input['admin_id'] = $this->enterpriseAdmin['id'];
        $input['route_url'] = ltrim($this->request->getRequestUri(), '\/gw-console\/v1\/');
        $res = (new LuckMicro())->exportActivity($input);
        return success($res);
    }
}