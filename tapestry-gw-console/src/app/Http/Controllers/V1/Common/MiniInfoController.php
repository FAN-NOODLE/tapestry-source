<?php

namespace App\Http\Controllers\V1\Common;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\MiniProgram;

class MiniInfoController extends BaseController
{
    public function index()
    {
        $appData  = [];
        $appIdArr = (new Auth())->getShopAppid(['enterprise_id' => $this->enterpriseId]);
        if (!empty($appIdArr['data']['shop_app_id'])) {
            $appData = (new MiniProgram())->detail(['app_id' => $appIdArr['data']['shop_app_id']]);
            if (!empty($appData['data'])) {
                $appData = $appData['data'];
            }
        }
        return success($appData);
    }
}
