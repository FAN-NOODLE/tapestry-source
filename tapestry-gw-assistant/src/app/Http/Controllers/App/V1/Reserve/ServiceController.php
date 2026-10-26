<?php

namespace App\Http\Controllers\App\V1\Reserve;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\{Request, JsonResponse};

class ServiceController extends BaseController
{
    /**
     * @var Sdk
     */
    protected $microScrm;

    public function __construct()
    {
//        parent::__construct();
        $this->microScrm = Common::microScrm();
    }

    /**
     * 服务列表
     * @param Request $request
     * @return JsonResponse
     */
    public function serviceList(Request $request): JsonResponse
    {
        $params = $request->all();
//        $params['enterprise_id'] = $params['enterprise_id'];
        $result = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/service/list')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }


}