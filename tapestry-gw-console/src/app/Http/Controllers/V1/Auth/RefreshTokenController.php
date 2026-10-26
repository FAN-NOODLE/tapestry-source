<?php

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\JsonResponse;

class RefreshTokenController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function updateMiniProgramToken(): JsonResponse
    {
        $status = (new Enterprise())->refreshToken(['enterprise_id' => $this->enterpriseId]);

        return success($status);
    }
}