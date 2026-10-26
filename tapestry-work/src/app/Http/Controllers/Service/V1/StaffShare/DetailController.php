<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffShareService;

class DetailController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'integer',
            'store_id'      => 'integer',
            'worker_id'     => 'integer',
            'user_id'       => 'integer',
            'binding_type'  => 'required',
        ];
        $message = [
            'enterprise_id.integer'  => 'enterprise_id必须为int类型',
            'store_id.integer'       => 'store_id必须为int类型',
            'worker_id.integer'      => 'worker_id必须为int类型',
            'user_id.integer'        => 'user_id必须为int类型',
            'binding_type.required'  => 'binding_type参数必填',
        ];
        $input = $this->validate($this->request, $rules, $message);
        $data = (new StaffShareService())->detail($input);

        return success($data);
    }
}
