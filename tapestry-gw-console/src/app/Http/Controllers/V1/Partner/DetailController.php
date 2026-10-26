<?php


namespace App\Http\Controllers\V1\Partner;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Partner;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only([
            'page',
            'size',
            'id',
            'enterprise_id',
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new Partner())->detail($params);

        return success($result);
    }

}