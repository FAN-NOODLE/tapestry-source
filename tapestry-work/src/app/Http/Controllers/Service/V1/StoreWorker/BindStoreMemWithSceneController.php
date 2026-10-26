<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class BindStoreMemWithSceneController extends BaseController
{
    /**
     * 绑定门店员工关系(带参数访问门店)
     * @return array|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $rules = [
            'enterprise_id'  => 'required',
            'scene_store_id' => 'required',
            'scene_user_id'  => 'required',
            'store_id'       => 'required'
        ];
        $message = [
            'enterprise_id.required'  => 'enterprise_id参数必填',
            'scene_store_id.required' => 'scene_store_id参数必填',
            'scene_user_id.required'  => 'scene_user_id参数必填',
            'store_id.required'       => 'store_id参数必填'
        ];
        $this->validate($this->request, $rules, $message);
        $params = $this->request->only([
            'enterprise_id',
            'scene_store_id',
            'scene_user_id',
            'store_id',
            'scene_worker_id'
        ]);
        $data = StoreWorkerService::storeMemBindWithScene($params);

        return success($data);
    }
}
