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

class SettingSaveController extends BaseController
{
    /**
     * 保存商户信息
     *
     * @param Request $request
     * @return mixed|void
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['logo', 'short_name']);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new Enterprise())->enterpriseMerchantSettingSave($params);

        return success($result['data'] ?? []);
    }
}
