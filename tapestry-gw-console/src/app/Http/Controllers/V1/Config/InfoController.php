<?php

namespace App\Http\Controllers\V1\Config;

use CrsSdk\Boss\EnterpriseConfig;
use App\Http\Controllers\V1\BaseController;

class InfoController extends BaseController
{
    public function index()
    {
        $params                 = $this->request->only(['type']);
        $where['enterprise_id'] = $this->enterpriseId;

        if (!empty($params['type'])) {
            $where['type'] = $params['type'];
        }

        $info   = (new EnterpriseConfig())->info($where);
        $result = isset($info['data']) ? $info['data'] : $info;

        return success(['list' => $result]);
    }
}
