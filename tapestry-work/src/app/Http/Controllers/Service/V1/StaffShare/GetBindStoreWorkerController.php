<?php


namespace App\Http\Controllers\Service\V1\StaffShare;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffShareService;

class GetBindStoreWorkerController extends BaseController
{

    public function index()
    {
        $rules = [
            'user_id' => 'required|integer',
        ];
        $message = [
            'user_id.required' => 'user_id参数必填',
            'user_id.integer'  => 'user_id参数必须为int类型',
        ];
        $this->validate($this->request, $rules, $message);
        $userId = $this->request->input('user_id');
        $data = (new StaffShareService())->getBindWorkerInfo($userId);

        return success($data);
    }
}