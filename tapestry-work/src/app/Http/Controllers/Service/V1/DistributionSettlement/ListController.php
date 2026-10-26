<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/16 0016
 * Time: 11:02
 */

namespace App\Http\Controllers\Service\V1\DistributionSettlement;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\Settlement\ListService;

class ListController extends BaseController
{
    private $listService;

    public function __construct()
    {
        parent::__construct();

        $this->listService = new ListService($this->request->input());
        $enterpriseId = $this->request->input('enterprise_id');
        switchAppHost(['enterprise_id' => $enterpriseId]);
    }

    public function index()
    {
        return success($this->listService->list());
    }

}