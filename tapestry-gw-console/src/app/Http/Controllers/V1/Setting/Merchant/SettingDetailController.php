<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/21 17:07
 */

namespace App\Http\Controllers\V1\Setting\Merchant;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;
use App\Exceptions\Error;

class SettingDetailController extends BaseController
{
    /**
     * 获取商户信息
     *
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new Enterprise())->enterpriseMerchantSettingInfo($params);

        return success($data['data']);
    }
}
