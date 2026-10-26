<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class SftpErrorRecordListController extends BaseController
{

    /**
     * 错误日志列表
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'enterprise_id',
            'file_log_id',
        ]);
        if (empty($params['enterprise_id']) || empty($params['file_log_id'])) {
            throw new Error(600000, '文件id不能为空');
        }
        $result = StoreService::getErrorRecordList($params['file_log_id']);
        return success($result);
    }
}
