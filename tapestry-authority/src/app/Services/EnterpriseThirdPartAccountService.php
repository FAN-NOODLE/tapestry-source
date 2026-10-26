<?php
/**
 * Created by PhpStorm.
 * User: wpf
 * Date: 2020/4/12
 * Time: 16:52
 */

namespace App\Services;

use App\Models\EnterpriseThirdPartAccount;

class EnterpriseThirdPartAccountService extends BaseService
{
    public function getMiniProgram($enterpriseId)
    {
        return (new EnterpriseThirdPartAccount())->getMiniProgram($enterpriseId);
    }
}
