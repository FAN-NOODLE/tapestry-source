<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use Illuminate\Validation\ValidationException;

class GetStoreIdsController extends BaseController
{
    /**
     * @return array
     * @throws ValidationException
     */
    public function index()
    {
        $input = $this->request->input();
        $rules = [
            'enterprise_id' => 'required',
            'mobile'        => 'required',
        ];
        $message = [
            'mobile.required'        => 'mobile参数必填',
            'enterprise_id.required' => 'enterprise_id参数必填'
        ];
        $this->validate($this->request, $rules, $message);
        $result = StoreWorkerService::getStoreIdsByMobile($input['mobile'], $input['enterprise_id']);
        if (!empty($result)) {
            $result = array_column($result, 'store_id');
        }
        return success($result);
    }
}
