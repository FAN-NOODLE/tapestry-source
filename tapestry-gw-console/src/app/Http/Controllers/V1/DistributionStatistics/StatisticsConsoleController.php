<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-02
 * Time: 16:15
 */

namespace App\Http\Controllers\V1\DistributionStatistics;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionStatistics;
use Illuminate\Http\Request;

class StatisticsConsoleController extends BaseController
{
    /**
     * 分销统计导购首页
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params                  = $request->only(['worker_ids', 'start_time', 'end_time']);
        $params['enterprise_id'] = $this->enterpriseId;
        $data                    = (new DistributionStatistics())->statisticConsole($params);

        $data = $this->formDataFormat($data);

        return success($data);
    }


    /**
     * 分销统计数据概览格式转换
     *
     * @param        $data
     * @param string $key1
     * @param string $key2
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 20:53
     */
    protected function formDataFormat($data, $key1 = 'settlement_num', $key2 = 'settlement_money')
    {
        $x = array_keys($data['detail_data']);
        $y = [
            'data1' => array_column($data['detail_data'], $key1) ?? [],
            'data2' => array_column($data['detail_data'], $key2) ?? [],
        ];

        $result = [
            'data1'       => $data[$key1],
            'data2'       => $data[$key2],
            'detail_data' => [
                'x' => $x,
                'y' => $y,
            ],
        ];

        return $result;
    }
}
