<?php
/**
 * Created by PhpStorm.
 * User: HuiYang
 * Time: 2020年11月18日19:34:38
 */

namespace App\Http\Controllers\Service\V1\Order;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CrsOrder;
use App\Services\GuideManagerService;
use CrsSdk\Shop\Order;

class OrderStatisticalClassificationController extends BaseController
{
    /**
     * 订单分类统计
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/19 14:13
     */
    public function index()
    {
        $input        = $this->getParam();

        $returnResult = [
            'staff_order_num' => 0,
            'partner_order_num' => 0
        ];
        $data = CrsOrder::query()
            ->select(DB::raw('count("partner_id") as num'), 'partner_id')
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('status', '>=', 30)// 已完成支付
            ->where('pay_at', '>=', $input['start_time'])
            ->where('pay_at', '<=', $input['end_time'])
            ->where('store_staff_id', '<>', 0)
            ->where('store_staff_id', $input['admin_id'])
            ->groupBy('partner_id')
            ->get();
        if (!empty($data)) {
            $data = $data->toArray();
            $groupData = array_column($data, 'partner_id');
            if (isset($groupData[0])) {
                $returnResult['staff_order_num'] = $groupData[0]['num'];
                unset($groupData[0]);
            }
            if (!empty($groupData)) {
                $returnResult['partner_order_num'] = array_sum(array_column($groupData, 'num'));
            }
        }

        return success($returnResult);

    }

    /**
     * 验证参数
     *
     * @return array|null|string
     * @author  HuiYang
     * @date    2020年11月18日19:37:3
     */
    private function getParam()
    {
        $input    = $this->request->input();
        $rules    = [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'admin_id'      => 'required',
            'start_time'    => 'required',
            'end_time'      => 'required',
        ];
        $messages = [
            'enterprise_id.*'     => '企业id参数不合法',
            'store_id.*'          => '门店id参数不合法',
            'admin_id.required'   => 'admin_id参数缺失',
            'start_time.required' => 'start_time参数缺失',
            'end_time.required'   => 'end_time参数缺失',
        ];
        $this->validate($this->request, $rules, $messages);

        return $input;
    }
}