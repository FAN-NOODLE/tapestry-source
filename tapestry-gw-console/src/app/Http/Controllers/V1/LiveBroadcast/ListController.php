<?php

namespace App\Http\Controllers\V1\LiveBroadcast;

use CrsSdk\MicroActivity\Live;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('直播列表');
        //活动微服务
        $return = (new Live())->lists($input);

        return success($return);
    }
}
