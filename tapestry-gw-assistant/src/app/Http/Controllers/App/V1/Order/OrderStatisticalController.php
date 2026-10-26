<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/18 0008
 * Time: 17:37
 */

namespace App\Http\Controllers\App\V1\Order;

use CrsSdk\Work\OrderStatistical;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class OrderStatisticalController extends OrderBaseController
{
    /**
     * 订单统计
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/18 11:54
     */
    public function index(Request $request)
    {
        $param['start_time']     = $request->input('start_time');
        $param['end_time']       = $request->input('end_time');
        $param['store_id']       = $request->input('store_id', 0);
        $param['worker_ids']     = $request->input('worker_ids');
        $param['store_ids']      = $request->input('store_ids');
        $param['department_ids'] = $request->input('department_ids');
        $param['enterprise_id']  = $this->enterpriseId;
        $param['qy_domain']      = $this->loginAdminEnterprise['domain'];
        $param['admin_id']       =  $this->loginAdminEnterprise['admin_id'];
        $param['mobile']         = $this->mobile;
        Log::info('订单统计请求参数', [$param]);

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

        $result = (new OrderStatistical())->orderStatistical($param);
        // 转换数据格式
        $result = $this->transformDataFormat($result);

        return success($result);


    }
}
