<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class AreaStoreListController extends BaseController
{
    public function index()
    {
        return success([]);
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $result = (new StoreArea())->areaStoreList($input);

        return success($result);
    }
}
