<?php

namespace App\Http\Controllers\V1\Setting;

use Illuminate\Http\JsonResponse;
use CrsSdk\MicroAuthority\Account;
use App\Http\Controllers\V1\BaseController;

class BaseConfigController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $info = (new Account())->BaseInfoCheck(['enterprise_id' => $this->enterpriseId]);

        return success($info);
    }
}