<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/11 0011
 * Time: 14:26
 */

namespace App\Http\Controllers\App\V1\Order;

use CrsSdk\Shop\Order;

class SendGoodsBatchController extends OrderBaseController
{
    /**
     * 多物流发货
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $rules   = [
            'data'                     => 'required|array',
            'data.*.order_no'          => 'required',
            'data.*.logistic_identity' => 'required',
            'data.*.logistic_no'       => 'required'
        ];
        $message = [
            'data.*'                          => '[data]数据格式不符合规范',
            'data.*.order_no.required'        => '订单编号不能为空',
            'data.*.logistic_identity.required' => '物流标示不能为空',
            'data.*.logistic_no.required'       => '物流单号不能为空'
        ];
        $this->validate($this->request, $rules, $message);

        $params           = $this->request->only(['data']);
        $params['enterprise_id']   = $this->enterpriseId;
        $params['enterprise_hash'] = $this->loginAdminEnterprise['hash'];

        try {
           (new Order())->sendGoodsBatch($params);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }

        return success([]);
    }
}
