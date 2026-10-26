<?php

namespace App\Http\Controllers\V1\GuideManage;

use CrsSdk\Work\GuideManage;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class DeleteController extends BaseController
{
    /**
     * 删除导购
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params                  = $this->request->only(['mem_uid']);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain']        = $this->qyDomain;

        (new GuideManage())->delete($params);

        $this->writeOperatorLog('删除导购');

        return success();
    }
}
