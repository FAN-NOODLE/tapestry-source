<?php


namespace App\Http\Controllers\V1\Partner;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Partner;
use Illuminate\Http\Request;

class SwitchController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only([
            'forbidden_reason',
            'switch',
            'id'
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new Partner())->partnerSwitch($params);

        return success($result);
    }

}