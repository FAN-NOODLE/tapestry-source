<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/19 0019
 * Time: 11:04
 */

namespace App\Http\Controllers\Service\V1\DistributionSettlement;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\Settlement\ListForAppService;

class ListForAppController extends BaseController
{
    private $listForAppService;

    public function __construct()
    {
        parent::__construct();

        $this->listForAppService = new ListForAppService($this->request->input());
        $enterpriseId = $this->request->input('enterprise_id');
        switchAppHost(['enterprise_id' => $enterpriseId]);
    }

    public function index()
    {
        return success($this->listForAppService->list());
    }
}