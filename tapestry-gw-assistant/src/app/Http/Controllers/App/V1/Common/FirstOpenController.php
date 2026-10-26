<?php

namespace App\Http\Controllers\App\V1\Common;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class FirstOpenController extends BaseController
{

    /**
     * 首次打开
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request): JsonResponse
    {
        $params = $request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = (new GoMicroScrm())->firstOpen($params);
        return success(['is_first_open' => $result, 'type' => $params['type']]);
    }
}
