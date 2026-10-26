<?php

namespace App\Http\Controllers\App\V1\EnterpriseNew;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\MicroAuthority\Permissions;

class AdditionalTreeController extends BaseController
{
    /**
     * 获取选配附加功能
     * @return array
     */
    public function additionalTree()
    {
        $version = $this->request->input('version');

        $tree = (new Permissions())->getAdditional(['version' => $version]);

        return success(['data' => $tree]);
    }
}