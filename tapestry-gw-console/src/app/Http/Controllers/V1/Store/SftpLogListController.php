<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;

class SftpLogListController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     */
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $res = Sdk::work()->setApi('service/v1/store/sftp-log-list')
            ->post($params);
        return success($res);
    }
}
