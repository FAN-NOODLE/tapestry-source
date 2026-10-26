<?php


namespace App\Http\Controllers\V1\PartnerSettlement;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\PpartnerDistribution;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only([
            'page',
            'size',
            'enterprise_id',
            'id'
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new PpartnerDistribution())->detail($params);

        return success($result);
    }

}