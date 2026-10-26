<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use App\Services\OrderService;
use CrsSdk\Shop\Order;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class OrderStatisticalController extends BaseController
{
    /**
     * 数据罗盘订单统计
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @author  weihua <weihua@vchangyi.com>
     * @date    2020/9/7 10:15
     */
    public function index(Request $request)
    {
	    $param['start_time']    = $request->input('start_time');
	    $param['end_time']      = $request->input('end_time');
	    $param['worker_ids']    = $request->input('worker_ids');
	    $param['version']       = $request->input('version');
	    $param['enterprise_id'] = $this->enterpriseId;
        Log::info('订单统计请求参数', [$param]);
        if (!empty($param['worker_ids'])) {

            //通过导购ID数组查询合伙人信息
            $belongWorkerPartnerList = (new PartnerMember())->partnerList([
                'enterprise_id' => $this->enterpriseId,
                'worker_ids' => $param['worker_ids'],
                'page' => 1,
                'size' => 100000
            ]);
            $param['partner_ids'] = [];
            $belongWorkerPartnerList = $belongWorkerPartnerList['list'] ?? [];
            if (!empty($belongWorkerPartnerList)) {
                $param['partner_ids'] = array_column($belongWorkerPartnerList, 'id');
            }
        }
        // 获取订单统计数据
        $data = (new Order())->wisdomOrderStatistical($param);
        $orderService = new OrderService();
        $data = $orderService->transformDataFormat($data);
        return success($data);
    }

    /**
     * 按日期查询订单数据
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @author  weihua <weihua@vchangyi.com>
     * @date    2020/8/19 17:44
     */
    public function statisticalByDate(Request $request)
    {
	    $param['page']             = $request->input('page', 1);
	    $param['size']             = $request->input('size', 10);
	    $param['sort_type']        = $request->input('sort_type');
	    $param['statistical_sort'] = $request->input('statistical_sort');

        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['worker_ids']    = $request->input('worker_ids');
	    $param['version']            =  $request->input('version');
        $param['enterprise_id'] = $this->enterpriseId;
        Log::info('按日期订单统计请求 参数', [$param]);

        // 获取订单统计数据
        $data = (new Order())->wisdomOrderStatistical($param);
        foreach ($data['detail_data'] as $key => &$value) {
            $value['time'] = $key;
        }
        // 按时间倒序
        $sortTime = array_column($data['detail_data'], 'time');
        array_multisort($sortTime, SORT_DESC, $data['detail_data']);
        // 数据排序
        if (!empty($param['sort_type']) && !empty($param['statistical_sort'])) {
            switch ($param['sort_type']) {
                case 'extension_num':
                    $sortData = array_column($data['detail_data'], 'extension_num');
                    break;
                case 'extension_amount':
                    $sortData = array_column($data['detail_data'], 'extension_amount');
                    break;
                case 'service_num':
                    $sortData = array_column($data['detail_data'], 'service_num');
                    break;
                case 'service_amount':
                    $sortData = array_column($data['detail_data'], 'service_amount');
                    break;
            }
            if ($param['statistical_sort'] == 'desc') {
                array_multisort($sortData, SORT_DESC, $data['detail_data']);
            } else {
                array_multisort($sortData, SORT_ASC, $data['detail_data']);
            }
        }
        // 总数据
        $total = count($data['detail_data']);
        // 总页数
        $endPage = intval(ceil($total / $param['size']));
        // 偏移量
        $start = ($param['page'] - 1) * $param['size'];

        $result = [
	        'list'     => array_values(array_slice($data['detail_data'], $start, $param['size'])),
	        'page'     => intval($param['page']),
	        'size'     => intval($param['size']),
	        'total'    => $total,
	        'end_page' => $endPage,
        ];

        return success($result);
    }

    /**
     * 按导购查询订单统计
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/20 10:48
     */
    public function statisticalByGuide(Request $request)
    {
        $param['page']             = $request->input('page', 1);
        $param['size']             = $request->input('size', 10);
        $param['sort_type']        = $request->input('sort_type');
        $param['statistical_sort'] = $request->input('statistical_sort');

        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['worker_ids']    = $request->input('worker_ids');
	    $param['version']            =  $request->input('version');
        $param['enterprise_id'] = $this->enterpriseId;
        $param['domain']        = $this->qyDomain;
        Log::info('按导购订单统计请求参数', [$param]);

        // 获取订单统计数据
//        $data = (new Order())->wisdomOrderStatisticalByGuide($param);

        return success([]);
    }
}
