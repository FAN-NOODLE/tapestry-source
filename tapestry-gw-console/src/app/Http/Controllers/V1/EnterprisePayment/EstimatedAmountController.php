<?php
/**
 * 企业红包动态分布预计金额
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\EnterprisePayment;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use App\Library\RandomRedPacket;

class EstimatedAmountController extends BaseController
{

    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接值
        $params = $request->all();

        if (empty($params['number'])) {
            notice(1300001);
        }

        if (empty($params['min_amount'])) {
            notice(1300002);
        }

        if (empty($params['max_amount'])) {
            notice(1300003);
        }
        $bonus_count = $params['number'];
        $bonus_max   = $params['max_amount']; //此算法要求设置的最大值要大于平均值
        $bonus_min   = $params['min_amount'];

        $result_bonus = (new RandomRedPacket())->getBonus($bonus_count, $bonus_max, $bonus_min);

        $total_money = array_sum($result_bonus);

        $data = [
            'enterprise_payment_amount' => $total_money
        ];

        return success($data);
    }
}
