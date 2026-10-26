<?php


namespace App\Http\Controllers\V1\StaffBuyManager;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffExtSetting;

class FillController extends BaseController
{
    /**
     * 运营设置数据添加
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->input();

        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('添加或修改运营设置数据');

        return (new StaffExtSetting())->staffExtSettingFill($input);
    }

}