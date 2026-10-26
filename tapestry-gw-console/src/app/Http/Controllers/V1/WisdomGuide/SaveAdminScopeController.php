<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/28 0028
 * Time: 15:29
 */

namespace App\Http\Controllers\V1\WisdomGuide;


use App\Exceptions\Error;
use CrsSdk\MicroAuthority\Account;
use App\Http\Controllers\V1\BaseController;

class SaveAdminScopeController extends BaseController
{
    private $accountSdk;

    public function __construct()
    {
        parent::__construct();
        $this->accountSdk = new Account();
    }

    public function index()
    {
        $this->checkParam();
        try {
            $this->accountSdk->EditAdminEnterpriseDataRange($this->getSdkParam());
        } catch (\Throwable $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }

        return success([]);
    }

    private function checkParam()
    {
        $this->validate($this->request, [
            'guide_mobile' => 'required',
            'range'        => 'required|array',
            'range_type'   => 'required'
        ]);
    }

    private function getSdkParam()
    {
        return [
            'enterprise_id' => $this->enterpriseId,
            'mobile'        => $this->request->input('guide_mobile'),
            'range'         => $this->reformRange()
        ];
    }

    private function reformRange()
    {
        $range = $this->request->input('range');
        if (!$range) {
            throw new Error(10000, '请传入部门或门店id');
        }
        foreach ($range as $objectId) {
            if (strlen($objectId) > 10) {
                throw new Error(10000, 'tapestry二期修改（2022.10.24）：请前端同学将参数中range数组中的取值由接口返回数据中的dp_id改为dp_third_id');
            }
            $reformRange[] = [
                "type"      => $this->request->input('range_type'),
                "object_id" => $objectId
            ];
        }
        return $reformRange;
    }
}