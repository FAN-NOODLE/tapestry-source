<?php

namespace App\Http\Controllers\Service\V1\StoreWorker\Open;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use CrsSdk\Shop\Store;

class DetailController extends BaseController
{

    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $rules   = [
            'enterprise_id' => 'required',
            'mem_uid'       => 'required',
            'store_no'      => 'required',
        ];
        $message = [
            'enterprise_id.required' => 'id参数必须为数字',
            'mem_uid.required'       => 'id参数必填',
            'store_no.required'      => '门店ID不能为空',
        ];
        $input   = $this->validate($this->request, $rules, $message);
        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $storeInfo = (new Store())->detail([
            'store_no' => $input['store_no'],
            'enterprise_id' => $input['enterprise_id']
        ]);
        $input['store_id'] = $storeInfo['id'];
        $data              = StoreWorkerService::detail($input);
        return success($data);
    }

}
