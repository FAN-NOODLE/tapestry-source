<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/19 0019
 * Time: 10:49
 */

namespace App\Http\Controllers\App\V1\DistributionSettlement;


use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Http\Controllers\App\V1\VoidLoginProxy;
use CrsSdk\Work\DistributionSettlement;

class ListController extends BaseController
{
    private $settlementSdk;

    private  $workId;

    public function __construct()
    {

        parent::__construct();

        $this->settlementSdk = new DistributionSettlement();
        $this->workId = $this->storeWorkerDetail($this->enterpriseId,$this->request->input('store_id'))['id'];


    }

    public function index()
    {
        return success($this->settlementSdk->listForApp($this->getParam()));
    }

    private function getParam()
    {
        $page = $this->request->input('page')??Constant::PAGE;
        $size = $this->request->input('size')??Constant::LIMIT;

        return [
            'page'=>$page,
            'size'=>$size,
            'enterprise_id'=>$this->enterpriseId,
            'work_id'=>$this->workId
        ];
    }

}