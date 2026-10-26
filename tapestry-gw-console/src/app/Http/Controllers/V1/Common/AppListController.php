<?php

namespace App\Http\Controllers\V1\Common;

use CrsSdk\Work\Common;
use App\Http\Controllers\V1\BaseController;

class AppListController extends BaseController
{
    /**
     * 获取企业商城小程序与绑定小程序列表
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $params['qy_domain']      = $this->qyDomain;
        $params['app_identifier'] = $this->appIdentifier;
        $params['enterprise_id']  = $this->enterpriseId;

        $list = (new Common())->appList($params);

        return success(['list' => $list['data']]);
    }
}
