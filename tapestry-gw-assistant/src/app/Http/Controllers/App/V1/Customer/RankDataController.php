<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class RankDataController extends BaseController
{

    const RANK_LENGTH = 50;

    /**
     * 顾客数据排行榜
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->loginMobile();
        $params['qy_domain'] = $this->loginAdminEnterprise['domain'];

        $rankList = (new GoMicroScrm())->rank($params);
        
        return success($rankList);
    }


}