<?php


namespace App\Http\Controllers\V1\StorePerformance;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StorePerformance;

class PerformanceFetchController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return success((new StorePerformance())->fetch($this->getParam()));
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
            'list_filter_type'  => 'integer',
            'list_filter_key'   => 'string',
            'page'          => 'integer',
            'size'          => 'integer',
        ], [
            'store_id.*'            => '店铺ID不合法',
            'type.*'                => '业绩查询类型不合法',
            'start_time'            => '自定义查询开始时间不合法',
            'end_time'              => '自定义查询结束时间不合法',
            'list_filter_type.*'    => '员工查询类型不合法',
            'list_filter_key.*'     => '员工查询关键词不合法',
            'page.*'                => '分页传递错误',
            'size.*'                => '分页数传递错误',
        ]);
        $param = $this->request->only(['store_id', 'type', 'start_time', 'end_time', 'list_filter_type', 'list_filter_key', 'page', 'size','sort_field']);
        if (!isset($param['start_time']) || !$param['start_time']) {
            $param['start_time'] = 0;
        }
        if (!isset($param['end_time']) || !$param['start_time']) {
            $param['end_time'] = 0;
        }
        if (!isset($param['start_time'])) {
            $param['start_time'] = 0;
        }
        if (!isset($param['end_time'])) {
            $param['end_time'] = 0;
        }

        $param['enterprise_id'] = $this->enterpriseId;
        return $param;
    }
}

