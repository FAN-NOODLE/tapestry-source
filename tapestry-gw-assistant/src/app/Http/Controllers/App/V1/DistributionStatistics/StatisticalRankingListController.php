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

class StatisticalRankingListController extends Controller
{
    /**
     * 全员分销统计排行
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/18 11:54
     */
    public function index(Request $request)
    {
        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['type']          = $request->input('type', 1);
        $param['store_id']      = $request->input('store_id', 0);
        $param['enterprise_id'] = $this->enterpriseId;
        $param['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $param['admin_id']      = $this->loginAdminEnterprise['admin_id'];
        $param['mobile']        = $this->mobile;
        $param['role']          = $this->getRole($param['mobile']);
        Log::info('分销统计排行榜请求参数', [$param]);

        $result = (new DistributionStatistics())->statisticRanking($param);

        return success(['list' => $result]);
    }
}
