<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\OrderStatisticalDetailExportJob;

class OrderStatisticalDetailExportController extends BaseController
{
    /**
     * 文件名
     */
    const FILE_NAME = '订单明细列表';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'order/order-statistical-detail-export';

    /**
     * 导购顾客数据统计tag
     */
    const TAG = 'order-statistical-detail-export';

    public function index()
    {
        //查询条件校验
        $input = $this->validates();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['qy_domain']     = $this->qyDomain;
        $input['operate_name']  = $this->enterpriseAdmin['username'] ?? '';
        $input['operate_id']    = $this->enterpriseAdminId;
        $input['export_type']   = $input['export_type'] ?? 1;
        $input['tag']           = $input['tag'] ?? self::TAG;
        $input['route_url']     = self::ROUTE_URL;
        $input['version']       = $input['version'] ?? '';
        $time                   = date('YmdHis', time());

        $name = self::FILE_NAME.$time;
        (new OrderStatisticalDetailExportJob($input, $input, $name))->handle();

        return success();

    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'worker_ids'    => 'array',
            'start_time'  => 'required',
            'end_time'    => 'required',
        ];

        $message = [
            'worker_ids.array'       => '导购id必须是数组',
            'start_time.required'  => 'start_time参数缺失',
            'end_time.required'    => 'end_time参数缺失',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }

}

