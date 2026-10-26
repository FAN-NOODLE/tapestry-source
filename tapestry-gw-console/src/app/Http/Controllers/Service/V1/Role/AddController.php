<?php


namespace App\Http\Controllers\Service\V1\Role;


use App\Http\Controllers\V1\BaseController;
use App\Models\Role;
use Illuminate\Http\Request;

class AddController extends BaseController
{
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);

        $data = [
            'enterprise_id' => $enterpriseId,
            'name'          => '店长',
            'describe'      => '店长',
            'type'          => 4
        ];

        Role::query()->create($data);

        return success();
    }
}
