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
use CrsSdk\Work\DistributionSettlement;

class DetailController extends BaseController
{
    private $settlementSdk;

    public function __construct()
    {
        parent::__construct();
        $this->settlementSdk = new DistributionSettlement();
    }

    public function index()
    {
        return success($this->settlementSdk->detail($this->getParam()));

    }

    private function getParam()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;
        return $param;
    }
}