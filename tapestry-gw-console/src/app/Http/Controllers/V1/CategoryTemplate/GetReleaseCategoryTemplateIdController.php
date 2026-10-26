<?php


namespace App\Http\Controllers\V1\CategoryTemplate;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CategoryTemplate;

class GetReleaseCategoryTemplateIdController extends BaseController
{

    public function index()
    {
        $input = [
            'enterprise_id' => $this->enterpriseId
        ];

        return success((new CategoryTemplate())->getReleaseTemplateId($input));
    }

}