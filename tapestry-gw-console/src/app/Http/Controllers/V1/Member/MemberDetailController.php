<?php

namespace App\Http\Controllers\V1\Member;

use CrsSdk\Shop\Order;
use CrsSdk\Work\Customer;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class MemberDetailController extends BaseController
{
    /**
     * 客户详情
     * @return JsonResponse
     * @throws \Exception
     */
    public function memberDetail()
    {
        $params                   = $this->request->all();
        $params['qy_domain']      = $this->qyDomain;
        $params['enterprise_id']  = $this->enterpriseId;
        $params['app_identifier'] = $this->appIdentifier;

        $data  = (new Customer)->customerDetail($params);
        $order = (new Order)->memberOrderList($params);

        $data['phone']      = $order['phone'];
        $data['order_data'] = $order['order_data'];

        return success($data);
    }
}