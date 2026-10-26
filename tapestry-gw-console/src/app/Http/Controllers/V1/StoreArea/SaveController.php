<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class SaveController extends BaseController
{
    /**
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return success((new StoreArea())->areaSave($this->getParam()));
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'id'            => 'integer',
            'pid'           => 'required|integer',
            'sort'          => 'required|integer',
            'name'          => 'required|string',
            'front_id'      => 'required|integer',
        ], [
            'id.*'      => '区域ID不合法',
            'pid.*'     => '父级ID不合法',
            'sort.*'    => '排序值不合法',
            'name.*'    => '区域名称不合法',
            'front_id.*'=> '前端ID不合法',
        ]);
        $saveArray = $this->request->only(['pid', 'sort', 'name', 'front_id']);
        $saveArray['enterprise_id'] = $this->enterpriseId;
        $saveArray['id'] = $this->request->get('id', 0);
        return $saveArray;
    }
}
