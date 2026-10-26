<?php


namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class BindStoreMemController extends BaseController
{

    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'store_id'      => 'integer',
            'user_id'       => 'integer',
        ];
        $message = [
            'enterprise_id.required' => 'enterprise_id参数必填',
            'enterprise_id.integer'  => 'enterprise_id必须为int类型',
            'store_id.integer'       => 'store_id必须为int类型',
            'user_id.integer'        => 'user_id必须为int类型',
        ];
        $input = $this->validate($this->request, $rules, $message);
        $data = StoreWorkerService::storeMemBind($input);

        return success($data);
    }
}
