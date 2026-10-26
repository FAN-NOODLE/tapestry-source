<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-20
 * Time: 17:20
 */

namespace App\Http\Controllers\V1\BaPerformance;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class BatchPosterController extends BaseController
{
    /**
     * 批量生成海报
     *
     * @param Request $request
     * @param StoreWorker $storeWorker
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request, StoreWorker $storeWorker)
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain']        = $this->enterprise['domain'];
        $params['operation_id']  = $this->enterpriseAdminId;
        $params['operator_name'] = $this->enterpriseAdmin['username'];
        try {
            $storeWorker->makePosterBatch($params);
        } catch (Error $error) {
            Log::error("调用Shop:StoreWorker->makePosterBatch错误，错误信息：{$error->getMessage()}, 错误code：{$error->getCode()}", $params);
            notice($error->getCode(), $error->getMessage());
        }
        return success();
    }
}
