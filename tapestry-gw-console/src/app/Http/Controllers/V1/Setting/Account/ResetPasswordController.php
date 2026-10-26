<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;


/**
 * 编辑管理员账号
 * Class EditController
 * @package App\Http\Controllers\V1\account
 */
class ResetPasswordController extends BaseController
{



    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        if(!$this->isSuperAdmin) {
            notice(2300001);
        }

         Sdk::authority()
            ->setApi('micro-authority/v1/account/reset-password')
            ->post($input);
        return success([]);

    }

}