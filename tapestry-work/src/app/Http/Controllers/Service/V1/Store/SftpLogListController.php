<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class SftpLogListController extends BaseController
{

    /**
     * 门店列表
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'enterprise_id',
            'file_name',
            'start_time',
            'end_time',
            'page',
            'size',
            'is_page'
        ]);
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        }
        $fields = ['id','file_name','created_at'];
        $result = StoreService::getSftpLogList($params, $params['is_page'] ?? true, $fields);
        return success($result);
    }
}
