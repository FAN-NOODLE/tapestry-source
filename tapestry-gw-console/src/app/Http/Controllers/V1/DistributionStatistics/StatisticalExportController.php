<?php
/**
 * 异步导出分销统计
 * User: 龚国新
 * Date: 2020年2月6日
 * Time: 10:54:50
 */

namespace App\Http\Controllers\V1\DistributionStatistics;

use App\Common\Constant;
use App\Jobs\Export\DistributionStatisticalExportJob;
use App\Http\Controllers\V1\BaseController;

class StatisticalExportController extends BaseController
{
    /*
     * 文件名
     */
    const FILE_NAME = '分销统计列表';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'distribution-statistics/statistic-export';

    /**
     * 导出分销数据tag
     */
    const TAG = 'export-distribution-statistical';

    public function index()
    {
        //查询条件校验
        $input = $this->validates();

        $input['size']          = Constant::MAX_DEFAULT_SIZE;
        $input['page']          = Constant::PAGE;
        $input['enterprise_id'] = $this->enterpriseId;
        $input['qy_domain']        = $this->qyDomain;
        $input['operate_name']  = $this->enterpriseAdmin['username'] ?? '';
        $input['operate_id']    = $this->enterpriseAdminId;
        $input['export_type']   = $input['export_type'] ?? 1;
        $input['tag']           = $input['tag'] ?? self::TAG;
        $input['route_url']     = $input['route_url'] ?? self::ROUTE_URL;
        $input['version']       = $input['version'] ?? '';
        $time                   = date('YmdHis', time());
        // 导出
        (new DistributionStatisticalExportJob($input, $input, self::FILE_NAME . $time))->handle();

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
            'work_ids'    => 'array',
            'start_time'  => 'required',
            'end_time'    => 'required',
            'export_type' => 'required',
        ];

        $message = [
            'work_ids.array'       => '导购id必须是数组',
            'start_time.required'  => 'start_time参数缺失',
            'end_time.required'    => 'end_time参数缺失',
            'export_type.required' => '导出数据类型不能为空',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }

}
