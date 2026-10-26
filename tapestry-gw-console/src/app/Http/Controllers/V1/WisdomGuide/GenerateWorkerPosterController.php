<?php


namespace App\Http\Controllers\V1\WisdomGuide;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StoreWorker;
use Exception;
use Illuminate\Http\JsonResponse;

class GenerateWorkerPosterController extends BaseController
{
    /**
     * 手动生成导购海报
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $input                  = $this->request->only(['id']);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['domain']        = $this->qyDomain;
        $input['operator_id']   = $this->enterpriseAdminId;
        $input['operator_name'] = $this->enterpriseAdmin['username'];

        $result = (new StoreWorker())->generateWorkerPoster($input);

        return success($result);
    }
}
