<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/11 0011
 * Time: 14:26
 */

namespace App\Http\Controllers\App\V1\Order;

use CrsSdk\Shop\Order;

class QueryOrderLogisticController extends OrderBaseController
{
    /**
     * 查询物流列表
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index()
    {
        $params = $this->request->input();

        if (empty($params['order_no'])) {
            error(654015);
        }

        $params['enterprise_hash'] = $this->loginAdminEnterprise['hash'];

        try {
            $data = (new Order())->queryOrderLogistic($params);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }

        return success($data['list']);
    }
}
