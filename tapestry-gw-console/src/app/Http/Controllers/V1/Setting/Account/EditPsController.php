<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;
use function GuzzleHttp\Psr7\try_fopen;

/**
 * 编辑管理员账号
 * Class EditController
 * @package App\Http\Controllers\V1\account
 */
class EditPsController extends BaseController
{
    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    public function __construct()
    {
        $isRequireLogin = request()->input('is_require_login', true);

        if($isRequireLogin){
            $this->requireLogin = true;
        }

        parent::__construct();
    }

    /**
     * 账号编辑
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @api {domain}/gw-console/v1/setting/account/edit
     */
    public function index()
    {
        $input = $this->request->only([
            'old_ps', 'new_ps', 'new_ps_two', 'account', 'enterprise_id'
        ]);

        if($this->requireLogin){
            $input['enterprise_id'] = $this->enterpriseId;
        }

        $input['operator_id']   = $this->enterpriseAdminId;
        $input['admin_id']      = $this->enterpriseAdminId;
        $input['is_require_login'] = $this->requireLogin;

        //AccountRequest::edit($input);

        (new Account())->editAccountPs($input);

        return success([]);
    }

}