<?php

namespace App\Http\Controllers\V1\Template\Version;

use CrsSdk\Shop\Version;
use App\Http\Controllers\V1\BaseController;

class ReleaseController extends BaseController
{
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("确认发布");

        return (new Version())->release($input);
    }
}
