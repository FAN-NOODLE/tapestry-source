<?php

namespace App\Http\Controllers\App\V1\Common;


use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\StoreArea;
use Illuminate\Support\Facades\Auth;

class AreaStoreListController extends Controller
{
    protected $isRequireLogin = false;

    public function index()
    {
        $enterpriseId = Auth::user()['enterprise_id'];

        $list = (new StoreArea())->areaStoreList(['enterprise_id' => $enterpriseId]);

        return success($list);
    }
}
