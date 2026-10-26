<?php


namespace App\Http\Controllers\V1\StorePerformance;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StorePerformance;

class PerformanceExportController extends BaseController

{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return success((new StorePerformance())->export($this->getParam()));
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'store_id'      => 'required|integer',
            'type'          => 'integer',
            'start_time'    => 'integer',
            'end_time'      => 'integer',
            'list_filter_type'  => 'integer',
            'list_filter_key'   => 'string',
            'route_url'     => 'required|string',
            'tag'           => 'required|string',
        ], [
            'store_id.*'            => '店铺ID不合法',
            'type.*'                => '业绩查询类型不合法',
            'start_time.*'          => '自定义查询开始时间不合法',
            'end_time.*'            => '自定义查询结束时间不合法',
            'list_filter_type.*'    => '员工查询类型不合法',
            'list_filter_key.*'     => '员工查询关键词不合法',
            'route_url.*'           => '下载路由不合法',
            'tag.*'                 => 'TAG不合法'
        ]);
        $param = $this->request->only(['store_id', 'type', 'start_time',
            'end_time', 'list_filter_type', 'list_filter_key', 'route_url', 'tag','sort_field']);
        if (!isset($param['list_filter_type']) || !$param['list_filter_type']) {
            $param['list_filter_type'] = 0;
        }
        if (!isset($param['list_filter_key']) || !$param['list_filter_key']) {
            $param['list_filter_key'] = '';
        }
        if (!isset($param['start_time']) || !$param['start_time']) {
            $param['start_time'] = 0;
        }
        if (!isset($param['end_time']) || !$param['start_time']) {
            $param['end_time'] = 0;
        }
        $param['operator_id'] = $this->enterpriseAdmin['id'];
        $param['operator_name'] = $this->enterpriseAdmin['username'];
        $param['enterprise_id'] = $this->enterpriseId;
        $param['enterprise_hash'] = $this->enterprise['hash'];
        return $param;
    }

}
