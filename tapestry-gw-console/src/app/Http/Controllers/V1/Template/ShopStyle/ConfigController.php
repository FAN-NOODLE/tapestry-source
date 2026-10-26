<?php

namespace App\Http\Controllers\V1\Template\ShopStyle;

use CrsSdk\Shop\Template;
use App\Http\Controllers\V1\BaseController;

class ConfigController extends BaseController
{
    public function index()
    {
        $info = (new Template())->shopStyleConfig(['enterprise_id' => $this->enterpriseId]);

        return success(['list' => $info['data']]);
    }
}
