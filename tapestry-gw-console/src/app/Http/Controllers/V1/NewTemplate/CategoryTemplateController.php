<?php

namespace App\Http\Controllers\V1\NewTemplate;

use CrsSdk\Shop\NewTemplate;
use App\Http\Controllers\V1\BaseController;

class CategoryTemplateController extends BaseController
{
    public function categoryDetail()
    {
        $detail = (new NewTemplate())->CategoryDetail(['enterprise_id' => $this->enterpriseId]);

        return success($detail);
    }

    public function categorySave()
    {
        $detail = (new NewTemplate())->CategorySave([
            'content'       => $this->request->input('content'),
            'enterprise_id' => $this->enterpriseId
        ]);

        return success($detail);
    }
}