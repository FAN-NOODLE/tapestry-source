<?php


namespace App\Http\Controllers\V1\PartnerSettlement;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\PpartnerDistribution;
use Illuminate\Http\Request;

class listController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only([
            'page',
            'size',
            'enterprise_id',
            'search_type',
            'search_value',
            'status_type',
            'status_value',
            'worker_ids',
            'order_by_type',
            'orderby'
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new PpartnerDistribution())->Distributionlist($params);

        return success($result);
    }

}