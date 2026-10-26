<?php


namespace App\Http\Controllers\V1\Goods;


use App\Http\Controllers\V1\BaseController;
use App\Services\GoodsBatchChangeService;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Redis;
use Illuminate\Validation\ValidationException;

class GoodsBatchChangeSelectController extends BaseController
{
    public const HAS_RUNNING = 1;
    public const NO_RUNNING = 2;
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $redis = Redis::connection();
        $params = $this->getParams();
        $service = new GoodsBatchChangeService($this->enterpriseId, $params['change_type']);
        $result = $redis->get($service->getProgressSystemKey());
        if ($result) {
            return success(['has_running' => self::HAS_RUNNING]);
        }
        return success(['has_running' => self::NO_RUNNING]);
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