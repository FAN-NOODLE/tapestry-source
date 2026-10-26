<?php


namespace App\Http\Controllers\V1\Goods;


use App\Http\Controllers\V1\BaseController;
use App\Services\GoodsBatchChangeService;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Redis;
use Illuminate\Validation\ValidationException;

class GoodsBatchChangeResultController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $redis = Redis::connection();
        $params = $this->getParams();
        $service = new GoodsBatchChangeService($this->enterpriseId, $params['change_type']);
        $result = $redis->get($service->getResultKey());
        if ($result) {
            $array = json_decode($result, true);
            if ($array['is_complete'] === $service::STATUS_COMPLETE) {
                $redis->del($service->getProgressSystemKey());
            }
            return success($array);
        }
        return success($service->makeResultData($service::STATUS_NO_PROCESS));
    }

    /**
     * @return array
     * @throws ValidationException
     */
    protected function getParams(): array
    {
        $inStr = GoodsBatchChangeService::TYPE_COLLECT_STR;
        return $this->validate(request(), [
            'change_type'   => 'required|integer|in:' . $inStr,
        ], [
            'change_type.*' => '更改类型不合法',
        ]);
    }
}