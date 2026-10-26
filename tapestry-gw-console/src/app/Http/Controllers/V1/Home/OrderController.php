<?php
/**
 * 首页上半部分数据统计
 * User: zhuhongda
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Home;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;

class OrderController extends BaseController
{
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only(['date_type']);

        // 企业id
        $input['enterprise_id'] = $this->enterpriseId;

        // 默认今日
        if (!isset($input['date_type'])) {
            $input['date_type'] = Constant::DAY;
        }

        // 非超管按照门店筛选
        $roleInfo = $this->getStores($this->enterpriseAdmin['role']);
        if (!$roleInfo['is_super_admin']) {
            $input['stores'] = $roleInfo['stores'];
            $input['type'] = $roleInfo['type'];
        }

        // 查询销售额、订单数、会员数、客户数
        $sales = (new Order())->orderStatistical($input);

        return success([
            'final_amount'    => $sales['data']['final_amount'],
            'order_number'    => $sales['data']['order_number'],
            'user_deal_count' => $sales['data']['user_deal_count'],
            'profits'         => $sales['data']['profits'],
            'refund'          => $sales['data']['refund']
        ]);
    }
}
