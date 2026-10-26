<?php
/**
 * 首页指标达成率数据统计
 * User: zhuhongda
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Home;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Task;
use CrsSdk\Shop\Task as ShopTask;

class TaskController extends BaseController
{
    public function index()
    {

        // 定义默认数据
        $data = [
            'sales_task' => 0, 'coupons_task' => 0, 'member_task' => 0, 'customer_task' => 0
        ];

        $params = [
            'enterprise_id' => $this->enterpriseId
        ];

        // 非超管按照门店筛选
        $roleInfo = $this->getStores($this->enterpriseAdmin['role']);
        if (!$roleInfo['is_super_admin']) {
            $params['stores'] = $roleInfo['stores'];
            $params['type'] = $roleInfo['type'];
            
            // 调用SDK查询任务完成数据
            $list = (new Task())->targetCount($params)['data'];

            // 遍历计算百分比
            foreach ($list as $key => $value) {
                $percent = $value['complete_number'] > 0 ? round($value['complete_number'] / $value['total'] * 100, 2) : 0;
                switch ($key) {
                    case 'sales':
                        $data['sales_task'] = $percent;
                        break;
                    case 'customer':
                        $data['customer_task'] = $percent;
                        break;
                    case 'coupon':
                        $data['coupons_task'] = $percent;
                        break;
                    case 'member':
                        $data['member_task'] = $percent;
                        break;
                }
            }
        } else {
            $list = (new ShopTask())->taskList($params);

            if (!empty($list['data'])) {
                foreach ($list['data'] as $value) {
                    switch ($value['type']) {
                        case 1:
                            $data['sales_task'] = $value['task_number'];
                            break;
                        case 2:
                            $data['customer_task'] = $value['task_number'];
                            break;
                        case 4:
                            $data['coupons_task'] = $value['task_number'];
                            break;
                        case 5:
                            $data['member_task'] = $value['task_number'];
                            break;
                    }
                }
            }
        }

        return success($data);
    }
}
