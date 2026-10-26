<?php


namespace App\Http\Controllers\V1\StaffBuyManager;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffExtSetting;

class StaffExtDetailController extends BaseController
{
    public function index()
    {

        return (new StaffExtSetting())->staffExtSettingDetail(['enterprise_id' => $this->enterpriseId]);

    }

}