<?php

namespace App\Http\Controllers\V1\Logistics;

use App\Common\CacheKey;
use CrsSdk\Open\Logistics;
use App\Http\Controllers\V1\BaseController;

class KeyListController extends BaseController
{
    /**
     * 支持的快递公司编码接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $data = (new Logistics())->keyList();

        $key = sprintf(CacheKey::LAST_LOGISTIC_KEY, $this->enterpriseId);
        $redis = redis();
        $lastLogisticKey = $redis->get($key);

        $list = [];
        foreach ($data['data'] as $key => $value) {
            $list[] = [
                'key' => $key,
                'value' => $value,
                'is_last_use' => ($lastLogisticKey === $key)
            ];
        }

        return success(['list' => $list]);
    }
}
