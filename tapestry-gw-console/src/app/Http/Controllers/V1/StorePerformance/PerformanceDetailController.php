<?php

namespace App\Http\Controllers\V1\StorePerformance;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StorePerformance;

class PerformanceDetailController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return success((new StorePerformance())->detail($this->getParam()));
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'store_id'      => 'required|integer',
            'type'          => 'required|integer',
            'start_time'    => 'integer',
            'end_time'      => 'integer',
            'worker_id'     => 'required|integer',
            'page'          => 'integer',
            'size'          => 'integer'
        ], [
            'store_id.*'            => '店铺ID不合法',
            'type.*'                => '业绩查询类型不合法',
            'start_time.*'          => '自定义查询开始时间不合法',
            'end_time.*'            => '自定义查询结束时间不合法',
            'worker_id.*'           => '员工ID不合法',
            'page.*'                => '分页不合法',
            'size.*'                => '分页数量不合法',
        ]);
        $param = $this->request->only(['store_id', 'type', 'start_time',
            'end_time', 'worker_id', 'page', 'size','order_type']);
        if (!isset($param['start_time']) || !$param['start_time']) {
            $param['start_time'] = 0;
        }
        if (!isset($param['end_time']) || !$param['start_time']) {
            $param['end_time'] = 0;
        }
        if ($param['size'] > 50) {
            $param['size'] = 50;
        }
        if ($param['size'] < 1) {
            $param['size'] = 10;
        }
        $param['enterprise_id'] = $this->enterpriseId;
        return $param;
    }
}
