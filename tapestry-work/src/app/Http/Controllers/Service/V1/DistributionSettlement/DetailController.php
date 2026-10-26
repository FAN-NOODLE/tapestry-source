<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/16 0016
 * Time: 15:49
 */

namespace App\Http\Controllers\Service\V1\DistributionSettlement;


use App\Http\Controllers\Controller;
use App\Services\Settlement\DetailService;

class DetailController extends Controller
{
    private $detailService;

    public function __construct()
    {
        parent::__construct();
        $enterpriseId = $this->request->input('enterprise_id');
        switchAppHost(['enterprise_id' => $enterpriseId]);
        $this->detailService = new DetailService(
            $this->request->input('enterprise_id'),
            $this->request->input('id'),
            $this->request->input('version',"")
            );

    }

    public function index()
    {
        return success($this->detailService->detail());
    }

}