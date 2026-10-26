<?php


namespace App\Http\Controllers\V1\Withdrawal;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;

class ApprovalListController extends BaseController
{

    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $data = (new Withdrawal())->ApprovalList($input);

        return success($data);
    }

}