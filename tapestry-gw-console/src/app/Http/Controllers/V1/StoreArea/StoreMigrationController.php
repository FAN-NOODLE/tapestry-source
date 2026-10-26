<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class StoreMigrationController extends BaseController
{
    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        (new StoreArea())->storeMigration($this->getParam());
        return success();
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'store_ids'     => 'required|string',
            'from_area_id'  => 'required|integer',
            'to_area_id'    => 'required|integer',
        ], [
            'store_ids.*'   => '门店ID不合法',
            'from_area_id.*'=> '来源区域ID不合法',
            'to_area_id.*'  => '目标区域ID不合法'
        ]);
        $param = $this->request->only(['store_ids', 'from_area_id', 'to_area_id']);
        $param['enterprise_id'] = $this->enterpriseId;
        return $param;
    }
}
