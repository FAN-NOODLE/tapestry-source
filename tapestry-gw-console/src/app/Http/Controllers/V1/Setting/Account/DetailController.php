<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;

/**
 * 管理员账号详情
 * Class DetailController
 * @package App\Http\Controllers\V1\account
 */
class DetailController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'id'
        ]);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['qy_domain'] = $this->qyDomain;

        AccountRequest::detail($input);

        $detail = (new Account())->detailAccount($input);

        return success($detail);
    }
}