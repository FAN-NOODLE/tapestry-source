<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/28 0028
 * Time: 16:29
 */

namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseAdmin;
use App\Services\EnterpriseDataRangeService;

class EditAdminEnterpriseDataRangeController extends BaseController
{
    public function index()
    {

        $this->validate($this->request, [
            'mobile'        => 'required',
            'enterprise_id' => 'required',
            'range'         => 'required'
        ], [
            'mobile.required'        => '手机号必传',
            'enterprise_id.required' => '企业ID必传',
            'range.required'         => '管理企业必传',
        ]);

        $params = $this->request->all();
        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $params['enterprise_id'])
            ->where('mobile', '=', $params['mobile'])
            ->first();

        $dataRange = [
            'range'         => $params['range'],
            'id'            => $admin->id,
            'enterprise_id' => $admin->enterprise_id
        ];

        (new EnterpriseDataRangeService())->editDataRange($dataRange);

        return success();
    }

}