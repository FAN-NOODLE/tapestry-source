<?php


namespace App\Http\Controllers\V1\Withdrawal;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;

class ApprovalCheckController extends BaseController
{

    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['admin_id']      = $this->enterpriseAdminId;

        $data = (new Withdrawal())->approvalCheck($input);

        return success($data['data'] ?? $data);
    }

}