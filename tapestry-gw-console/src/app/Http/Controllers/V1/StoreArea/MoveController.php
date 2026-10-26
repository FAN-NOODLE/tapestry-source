<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class MoveController extends BaseController
{
    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        return success((new StoreArea())->areaMove($this->getParam()));
    }

    /**
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'pid'           => 'required|integer',
            'sort'          => 'required|integer',
            'name'          => 'required|string',
            'front_id'      => 'required|integer',
        ], [
            'pid.*'     => '父级ID不合法',
            'sort.*'    => '排序值不合法',
            'name.*'    => '区域名称不合法',
            'front_id.*'=> '前端ID不合法',
        ]);
        $saveArray = $this->request->only(['pid', 'sort', 'name', 'front_id']);
        $saveArray['enterprise_id'] = $this->enterpriseId;
        return $saveArray;
    }
}
