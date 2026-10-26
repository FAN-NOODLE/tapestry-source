<?php


namespace App\Http\Controllers\V1\StaffRight;


use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use CrsSdk\Shop\StaffRight;

class InheritStaffController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only(['worker_id', 'old_worker_id', 'ids', 'is_all']);

        $rules = [
            'worker_id' => 'required',
            'old_worker_id' => 'required',
        ];

        $message = [
            'worker_id.required' => '继承导购ID不能为空',
            'old_worker_id.required' => '原导购ID不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        (new StaffRight)->InheritStaff($params);

        return success();
    }

}