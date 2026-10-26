<?php

namespace App\Http\Controllers\V1\Page\Footer;

use CrsSdk\Shop\Footer;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    public function index()
    {
        $list = (new Footer())->detail(['enterprise_id' => $this->enterpriseId]);

        return success($list);
    }
}