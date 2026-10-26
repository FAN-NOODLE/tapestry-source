<?php

namespace App\Http\Controllers\V1\GuideManage;

use CrsSdk\Work\GuideManage;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 添加导购
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params                  = $this->request->only(['mem_uids']);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain']        = $this->qyDomain;

        (new GuideManage())->add($params);

        $this->writeOperatorLog('添加导购');

        return success();
    }
}
