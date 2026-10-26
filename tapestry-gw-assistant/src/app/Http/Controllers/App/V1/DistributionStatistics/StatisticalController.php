<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/18 0008
 * Time: 17:37
 */

namespace App\Http\Controllers\App\V1\DistributionStatistics;

use CrsSdk\Work\DistributionStatistics;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\App\V1\BaseController as Controller;

class StatisticalController extends Controller
{
    /**
     * 全员分销统计
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
        Log::info('分销统计首页请求参数', [$param]);
        $result = (new DistributionStatistics())->distributionStatistic($param);

        // 转换数据格式
        $result = $this->formDataFormat($result);

        return success($result);
    }

    /**
     * 数据格式转化
     *
     * @param $data
     * @param string $key1
     * @param string $key2
     * @param string $key3
     * @param string $key4
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 19:22
     */
    public function formDataFormat($data, $key1 = 'settlement_num', $key2 = 'settlement_money')
    {
        $x = array_keys($data['detail_data']);
        $y = [
            'data1' => array_column($data['detail_data'], $key1) ?? [],
            'data2' => array_column($data['detail_data'], $key2) ?? [],
        ];

        $result = [
            'data1' => $data[$key1],
            'data2' => $data[$key2],
            'detail_data' => [
                'x' => $x,
                'y' => $y,
            ],
        ];

        return $result;
    }
}
