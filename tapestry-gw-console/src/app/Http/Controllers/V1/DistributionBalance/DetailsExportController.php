<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/8/4
 * Time: 10:01 下午
 */

namespace App\Http\Controllers\V1\DistributionBalance;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\DistributionBalanceDetailsExportJob;
use Illuminate\Http\Request;

class DetailsExportController extends BaseController
{


    /**
     * @param Request $request
     * @return int
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $all                   = $request->all();
        $all['size']          = Constant::MAX_DEFAULT_SIZE;
        $all['page']          = Constant::PAGE;
        $all['domain']          = $this->qyDomain;
        $all['qy_domain']      = $this->qyDomain;
        $all['app_identifier'] = $this->appIdentifier;
        $all['enterprise_id']  = $this->enterpriseId;
        $all['operate_id']     = $this->enterpriseAdminId;
        $all['operate_name']   = $this->enterpriseAdmin['username'];
        $all['tag']            = 'distribution_balance_details';
        $all['route_url']      = 'distribution-balance/details-export';

        $downloadName = "全员分销-余额管理导购详情_" . date("YmdHis");
        try {
            (new DistributionBalanceDetailsExportJob($all, $all, $downloadName))->handle();
        } catch (\Exception $error) {
            notice($error->getCode(), $error->getMessage());
        }
        return success();
    }

}