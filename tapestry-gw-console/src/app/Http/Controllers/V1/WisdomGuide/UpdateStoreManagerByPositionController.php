<?php


namespace App\Http\Controllers\V1\WisdomGuide;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Common;
use Illuminate\Http\JsonResponse;

class UpdateStoreManagerByPositionController extends BaseController
{
    public function index(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['is_update'] = $params['is_update'] ?? 0;
        $result = Common::work()->setApi('service/v1/store-worker/update-store-manager-by-position')->post($params);
        return success($result);
    }
}