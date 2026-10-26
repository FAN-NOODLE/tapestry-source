<?php

namespace App\Http\Controllers\Service\V1\StoreWorker\Open;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\StoreWorkerService;
use CrsSdk\Shop\Store;

class AddController extends BaseController
{

    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $rules   = [
            'enterprise_id' => 'required',
            'store_no'      => 'required',
            'mem_uid'       => 'required',
        ];
        $message = [
            'enterprise_id.required' => 'enterprise_id参数必填',
            'mem_uid.required'       => 'mem_uid参数必填',
            'store_no.required'      => '门店编号必填',
        ];
        $input   = $this->validate($this->request, $rules, $message);
        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $storeInfo = (new Store())->detail([
            'store_no' => $input['store_no'],
            'enterprise_id' => $input['enterprise_id']
        ]);
        $input['store_id']    = $storeInfo['id'];
        $input['operator_id'] = 0;
        $input['is_boss']     = StoreWorker::BOSS_NO;
        $input['source']      = StoreWorker::SOURCE_WORK;
        $data                 = StoreWorkerService::addStoreWorker($input);
        return success($data);
    }

}
