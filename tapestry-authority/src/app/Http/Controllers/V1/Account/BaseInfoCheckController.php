<?php

namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseService;

class BaseInfoCheckController extends BaseController
{
    public function index()
    {
        $params = $this->request->all();

        if (empty($params['enterprise_id'])) {
            notice(700014);
        }

        $config = EnterpriseService::baseConfigCheck($params['enterprise_id']);

        return success($config);
    }
}