<?php
/**
 * Created by PhpStorm.
 * @author Calm
 * @time 2020-09-24 18:01:37
 */


namespace App\Http\Controllers\V1\WorkerDistribution\Settlement;

use App\Http\Controllers\V1\BaseController;
use App\Http\Controllers\V1\VoidLoginProxy;
use CrsSdk\Work\DistributionSettlement;

class ListController extends BaseController
{
    private $settlementSdk;

    public function __construct()
    {
        parent::__construct();
        $this->settlementSdk = new DistributionSettlement();
    }

    public function index()
    {
        return success($this->settlementSdk->list($this->getParam()));
    }

    private function getParam()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;
        return $param;
    }
}