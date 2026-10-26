<?php


namespace App\Http\Controllers\V1\Admin;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\MicroAuthority\Common;

class SendSmsCodeController extends BaseController
{
    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    public function index()
    {
        $account = $this->request->input('mobile', '');
        $type    = $this->request->input('type', 0);

        if (empty($account)) {
            // 手机号不能为空
            error(451013);
        }

        // type == 1 验证该手机号是否在当前系统注册
        if ($type == 1) {
            $params = ['account' => '+86-' . ltrim($account, '+86-')];
            try {
                (new Account())->checkAccount($params);
            }catch (\Exception $e){
                notice(451026);
            }
        }

        (new Common())->sendCode(['mobile' => $account, 'client' => self::CLIENT]);

        return success();
    }
}
