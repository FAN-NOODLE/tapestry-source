<?php

namespace App\Http\Controllers\Service\V1\StoreWorker\Open;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use CrsSdk\Shop\Store;

class EditController extends BaseController
{

    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $rules   = [
            'mem_uid'       => 'required',
            'enterprise_id' => 'required',
            'mobile'        => 'required',
            'name'          => 'required',
            'store_no'      => 'required'
        ];
        $message = [
            'mem_uid.required'       => 'mem_uid手机号参数必填',
            'enterprise_id.required' => 'enterprise_id参数必填',
            'store_no.required'      => 'store_id参数必填',
            'mobile.required'        => '手机号参数必填',
            'name.required'          => '名称参数必填'
        ];
        $input   = $this->validate($this->request, $rules, $message);
        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $storeInfo = (new Store())->detail([
            'store_no' => $input['store_no'],
            'enterprise_id' => $input['enterprise_id']
        ]);
        $input['store_id'] = $storeInfo['id'];
        StoreWorkerService::editStoreWorkerByMemId($input);
        return success([]);
    }

}
