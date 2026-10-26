<?php

namespace App\Http\Controllers\V1\Template\Develop;

use CrsSdk\Shop\DevelopCode;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class SubmitController extends BaseController
{
    /**
     * @return JsonResponse|void
     * @throws ValidationException
     */
    public function index()
    {
        $input = $this->request->only(['develop_tag', 'develop_branch', 'develop_name']);

        $rule = [
            'develop_name' => 'required'
        ];

        $message = [
            'develop_name.required' => '开发码名称不能为空'
        ];

        $this->validate($this->request, $rule, $message);

        if (empty($input['develop_branch']) && empty($input['develop_tag'])) {
            notice(920311);
        }

        $input['enterprise_id'] = $this->enterpriseId;

        (new DevelopCode())->submit($input);

        return success();
    }
}