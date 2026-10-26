<?php

namespace App\Http\Controllers\V1\Member;

use CrsSdk\Shop\Order;
use CrsSdk\Work\Customer;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class MemberListController extends BaseController
{
    /**
     * 客户列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function memberList()
    {
        $all = $this->request->all();
        $all['qy_domain'] = $this->qyDomain;
        $all['app_identifier'] = $this->appIdentifier;
        $all['enterprise_id'] = $this->enterpriseId;

        // work获取列表数据
        $list = (new Customer)->customerList($all);

        $ids = $list['ids'];

        if (empty($ids)) {
            $order['data'] = [];
        } else {
            // 获取客户结果对应的订单数据
            $order = (new Order)->calcUserOrder([
                'ids'           => array_filter($ids),
                'enterprise_id' => $this->enterpriseId,
            ]);
        }

        $data = $list['list'];

        // 数据矫正
        foreach ($data['data'] as &$datum) {
            foreach ($order['data'] as $item) {
                if ($item['union_id'] == $datum['union_id']) {
                    $datum['phone'] = $item['phone'] ?? '';
                    $datum['order_total'] = $item['order_total'] ?? 0;
                    $datum['order_price_total'] = $item['order_price_total'] ?? 0;
                }
            }
        }

        return success([
            'list'  => $data['data'],
            'page'  => $data['meta']['page'],
            'size'  => $data['meta']['size'],
            'total' => $data['meta']['total'],
        ]);
    }
}