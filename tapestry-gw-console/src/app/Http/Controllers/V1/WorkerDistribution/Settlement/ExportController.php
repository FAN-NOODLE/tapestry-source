<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/16 0016
 * Time: 17:23
 */

namespace App\Http\Controllers\V1\WorkerDistribution\Settlement;


use App\Http\Controllers\V1\BaseController;
use App\Http\Controllers\V1\VoidLoginProxy;
use App\Jobs\Export\SettlementExportJob;

class ExportController extends BaseController
{
    public function index()
    {
        $all                   = request()->all();
        $all['qy_domain']      = $this->qyDomain;
        $all['app_identifier'] = $this->appIdentifier;
        $all['enterprise_id']  = $this->enterpriseId;
        $all['operate_id']     = $this->enterpriseAdminId;
        $all['operate_name']   = $this->enterpriseAdmin['username']??"";
        $all['tag']            = 'settlement_list';
        $all['route_url']      = 'worker-distribution/settlement/export';

        $downloadName = "全员分销-结算管理_" . date("YmdHis");;
        try {
            (new SettlementExportJob($all, $all, $downloadName))->handle();
        } catch (\Exception $error) {
            notice($error->getCode(), $error->getMessage());
        }
        return success();
    }

}