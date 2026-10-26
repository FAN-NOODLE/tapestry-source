<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/4/9
 * Time: 11:40
 */

namespace App\Http\Controllers\Service\V1\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class GetShopSecretController extends BaseController
{
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id');
        if (empty($enterpriseId)) {
            return error(600001);
        }

        $enterprise = EnterpriseThirdPartAccount::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->first();

        $data = [
            'app_id'     => $enterprise['app_id'],
            'app_secret' => $enterprise['app_secret']
        ];

        return success(['data' => $data]);
    }
}
