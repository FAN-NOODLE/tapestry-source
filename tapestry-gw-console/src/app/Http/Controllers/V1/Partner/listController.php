<?php


namespace App\Http\Controllers\V1\Partner;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Partner;
use Illuminate\Http\Request;

class listController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only([
            'page',
            'size',
            'search_type',
            'search_value',
            'partner_status',
            'commission',
            'commission_min',
            'commission_max',
            'orderby',
            'order_by_type',
            'worker_ids'
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new Partner())->partnerList($params);

        return success($result);
    }

}