<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class DeleteController extends BaseController
{
    /**
     * @throws \Exception
     */
    public function index()
    {
        return success((new StoreArea())->areaDelete($this->getParam()));
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'id'    => 'required|integer'
        ], [
            'id.*'  => '区域ID不合法,无法删除!'
        ]);
        return [
            'id'            => $this->request->id,
            'enterprise_id' => $this->enterpriseId,
        ];
    }
}
