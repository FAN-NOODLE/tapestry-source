<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/3/19
 * Time: 16:51
 */

namespace App\Http\Controllers\App\V1\DataList;


use App\Http\Controllers\App\V1\BaseController;
use App\Http\Controllers\App\V1\Task\My\ListController;
use App\Common\Constant;
use CrsSdk\Shop\Order;

class BlockController extends BaseController
{
    /**
     * 是否有用户信息 1：是，2：否 2
     */
    const ORDER_USER_FALSE = 2;

    /**
     * 应用块接口
     *
     * @return array
     */
    public function index()
    {
        // 获取订单数
        $num = 0;
        $orderNum = 0;
        $data = new ListController();
        $data = $data->index();

        if (!empty($data['data'])) {
            foreach ($data['data'] as $value) {
                if ($value['status'] == 1) {
                    $num ++;
                }
            }
        }
        // 开通商城
        if ($this->storeStatus == Constant::STORE_OPEN) {
            // 订单待付款数
            $input['search']['enterprise_id'] = $this->enterpriseId;
            $input['search']['is_user'] = self::ORDER_USER_FALSE;
            $input['search']['mem_uid'] = $this->memUid;
            $input['search']['status'] = 1;
            $input['size'] = $input['size'] ?? Constant::DEFAULT_SIZE;
            $input['page'] = $input['page'] ?? Constant::DEFAULT_PAGE;

            $list = (new Order())->orderList($input);

            $orderNum = $list['meta']['total'];
        }
       $list = [

           [
               'id' => 1,
               'name' => '业绩订单',
               'info' => '订单数据一目了然',
               'img' => 'https://oss-crs.vchangyi.com/wechat/performance_orders.png',
               'icon' => $orderNum
           ],
           [
               'id' => 2,
               'name' => '我的任务',
               'info' => '派发任务实时提醒',
               'img' => 'https://oss-crs.vchangyi.com/wechat/task.png',
               'icon' => $num
           ],
           [
               'id' => 3,
               'name' => '添加客户',
               'info' => '批量转移老客户',
               'img' => 'https://oss-crs.vchangyi.com/wechat/customer.png',
           ],
           [
               'id' => 4,
               'name' => '数据统计',
               'info' => '业绩数据统计',
               'img' => 'https://oss-crs.vchangyi.com/wechat/data_statistics.png',
           ]
       ];
        // 商城未开启
        if ($this->storeStatus == Constant::STORE_CLOSE) {
            array_shift($list);
        }
        return success(['data' => $list]);
    }
}
