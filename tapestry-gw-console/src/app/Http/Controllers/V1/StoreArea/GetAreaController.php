<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class GetAreaController extends BaseController
{
//    protected $requireLogin = false;
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return ['data' => (new StoreArea())->getArea($this->getParam())];
    }

    /**
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'pid'   => 'integer',
        ], [
            'pid.*' => '上级ID不合法',
        ]);
        return [
            'pid' => $this->request->get('pid', 0),
            'enterprise_id' => $this->enterpriseId,
        ];
    }
}
