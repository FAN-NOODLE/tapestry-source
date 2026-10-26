<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class BatchSaveController extends BaseController
{
    /**
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return success((new StoreArea())->batchSave($this->getParam()));
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'tree'            => 'required|json',
        ], [
            'tree.*'      => '提交数据不合法',
        ]);
        $saveArray['enterprise_id'] = $this->enterpriseId;
        $saveArray['tree'] = $this->request->tree;
        return $saveArray;
    }
}
