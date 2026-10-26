<?php
/**
 * 订单删除
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Home;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use CrsSdk\Shop\User;
use CrsSdk\Work\Customer;
use Exception;
use Illuminate\Http\Request;

class ChartController extends BaseController
{
    /**
     * 销售
     */
    const SALES = 1;

    public function index(Request $request)
    {
        // 列表
        try {
            $input = $request->only(['date_type', 'type']);

            // 企业id
            $input['enterprise_id'] = $this->enterpriseId;

            if (!isset($input['date_type'])) {
                $input['date_type'] = Constant::WEEK;
            }

            if (!isset($input['type'])) {
                $input['type'] = self::SALES;
            }

            // 临时变量存type，防止下面分门店查会覆盖掉
            $type = $input['type'];

            // 非超管按照门店筛选
            $role = $this->getStores($this->enterpriseAdmin['role']);
            if (!$role['is_super_admin']) {
                $input['stores'] = $role['stores'];
                $input['type'] = $role['type'];
            }

            // 调用SDK
            switch ($type) {
                case 1:
                    // 销售
                    return success((new Order())->salesChart($input));
                case 2:
                    // 会员
                    return success((new User())->userChart($input));
                    break;
                case 3:
                    // 客户
                    return success((new Customer())->chart($input)['data']);
                    break;
                case 4:
                    // 下单客户
                    return success((new Order())->userChart($input));
            }

        } catch (Exception $e) {
            error(451009);
        }

        return success();
    }
}
