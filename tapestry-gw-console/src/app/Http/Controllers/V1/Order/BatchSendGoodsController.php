<?php

namespace App\Http\Controllers\V1\Order;

use App\Common\CacheKey;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Exception;

class BatchSendGoodsController extends BaseController
{
    /**
     * 批量发货接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws Exception
     */
    public function index()
    {
        $this->validate($this->request, [
            'list'                     => 'present|array',
            'list.*.order_no'          => 'required|string',
            'list.*.logistic_identity' => 'required|string',
            'list.*.logistic_no'       => 'required|string',
        ]);
        $list = $this->request->input('list');

        $orderList = [];
        foreach ($list as $key => $value) {
            $orderList[] = [
                'order_no'          => $value['order_no'],
                'logistic_identity' => $value['logistic_identity'],
                'logistic_no'       => $value['logistic_no']
            ];
        }

        // 记录最后一次使用的物流公司标识
        if ($list !== []) {
            $lastLogisticKey = $list[count($list) - 1]['logistic_identity'];
            $key = sprintf(CacheKey::LAST_LOGISTIC_KEY, $this->enterpriseId);
            $redis = redis();
            $redis->set($key, $lastLogisticKey);
        }

        $data = (new Order())->sendGoodsBatch([
            'enterprise_id'   => $this->enterprise['id'],
            'enterprise_hash' => $this->enterprise['hash'],
            'data'            => $orderList
        ]);

        return success($data);
    }
}
