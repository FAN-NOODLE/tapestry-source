<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterpriseNew\EnterpriseNewService;

class AddController extends BaseController
{
    /**
     * 新增企业
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();

        $enterprise = (new EnterpriseNewService())->add($input);

        return success(['data' => $enterprise]);
    }
}