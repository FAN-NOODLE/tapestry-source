<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class ListController extends BaseController
{
    /**
     * @throws \Exception
     */
    public function index()
    {
        return success((new StoreArea())->areaList($this->getParam()));
    }

    public function getParam()
    {
        if (!$this->enterpriseId) {
            throw new Error(1000000, '请登录后再进行操作');
        }
        return [
            'enterprise_id'  => $this->enterpriseId
        ];
    }
}
