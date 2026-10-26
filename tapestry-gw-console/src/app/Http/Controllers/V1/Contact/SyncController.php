<?php

namespace App\Http\Controllers\V1\Contact;

use CrsSdk\Work\Contact;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class SyncController extends BaseController
{
    /**
     * 获取上次同步时间
     * @return JsonResponse
     * @throws \Exception
     */
    public function syncLastTime()
    {
        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new Contact())->getSyncOrgLastTime($params);

        return success($data['data'] ?? []);
    }

    /**
     * 获取同步状态
     * @return JsonResponse
     * @throws \Exception
     */
    public function status()
    {
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain'] = $this->qyDomain;

        $data = (new Contact())->getSyncOrgStatus($params);

        return success(['status' => $data['data']]);
    }

    /**
     * 同步组织结构
     * @return JsonResponse
     * @throws \Exception
     */
    public function syncOrg()
    {
        $params['qy_domain']      = $this->qyDomain;
        $params['enterprise_id']  = $this->enterpriseId;
        $params['app_identifier'] = $this->appIdentifier;
        $params['admin_id']       = $this->enterpriseAdminId;
        $params['admin_name']     = $this->enterpriseAdmin['username'];

        (new Contact())->syncOrg($params);

        return success();
    }
}