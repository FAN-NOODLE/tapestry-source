<?php

namespace App\Http\Controllers\V1\Template\PageVersion;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\TemplatePageVersion;

class DetailController extends BaseController
{
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $res = (new TemplatePageVersion())->detail($input);

        return success($res);
    }
}
