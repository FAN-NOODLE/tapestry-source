<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\JsonResponse;

class CloseOrOpenDepartmentController extends BaseController
{

    /**
     * @return JsonResponse
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     * @throws Exception
     */
    public function index(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        Sdk::work()->setApi('service/v1/store/close-or-open-department')
            ->post($params);
        return success([]);
    }
}
