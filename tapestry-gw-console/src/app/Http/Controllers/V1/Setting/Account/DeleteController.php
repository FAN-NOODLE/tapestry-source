<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;

/**
 * 删除管理员账号
 * Class DeleteController
 * @package App\Http\Controllers\V1\account
 */
class DeleteController extends BaseController
{

    public function index()
    {
        $input = $this->request->only([
            'id',
        ]);

        $input['enterprise_id'] = $this->enterpriseId;

        // 验证
        AccountRequest::del($input);

        (new Account())->deleteAccount($input);

        $this->writeOperatorLog('删除账号:id-'.$input['id']);
        return success([]);
    }

}