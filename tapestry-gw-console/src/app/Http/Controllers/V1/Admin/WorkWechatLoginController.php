<?php

namespace App\Http\Controllers\V1\Admin;

use App\Modules\LoginFactory;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Common;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class WorkWechatLoginController extends BaseController
{
    protected $requireLogin = false;

    /**
     * 1.跳转统一登录授权页面 如果是企业微信过来的可以直接走2
     *
     * @return \Illuminate\Http\RedirectResponse|\Laravel\Lumen\Http\Redirector
     * @throws \Illuminate\Validation\ValidationException
     */
    public function jump()
    {
        $this->validate($this->request, [
            'ref' => 'required|string'
        ]);

        $input = $this->request->all();

        $url = 'http://www.vchangyi.com/sso/index.php?corp_id=wxa7044ee8255576b0&ref=' . $input['ref'] . '&user_type=member';

        return redirect($url);
    }

    /**
     * 2.登录
     *
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function login()
    {
        $this->validate($this->request, [
            'type'                    => 'required|in:' . LoginFactory::SCAN_CODE . ',' . LoginFactory::WORK_WECHAT,
            'auth_code'               => 'string',
            'token'                   => 'string',
            'corpid'                  => 'string',
            'work_wechat_login_token' => 'string',
            /*'register_enterprise'     => 'array',*/
            'bind_account'            => 'array'
        ]);

        $input = $this->request->all();

        // 加载Login模块
        $login = LoginFactory::getInstance($input['type'], $input);

        try {
            $login->init($input['work_wechat_login_token'] ?? null);
        } catch (\Exception $exception) {
            Log::error('企业微信登录初始化异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'input'   => $input
            ]);

            error(100000, [$exception->getMessage()]);
        }

        $data = [
            'code'                    => 'success',
            'desc'                    => '可以登录',
            'auth_enterprise'      => null,
            'admin'                   => null,
            'tmp_token'               => '',
            'work_wechat_login_token' => (string)$login
        ];

        // 判断授权的管理员是否有效
        if ($login->authEffective === false) {
            $data['code'] = 'auth_invalid';
            $data['desc'] = '授权无效';

            return success($data);
        }

        // 是否需要注册企业
        if ($login->enterprise === null) {

            // 存在注册企业信息 进行注册
            if (isset($input['register_enterprise']) && !empty($input['register_enterprise'])) {
                $login->registerEnterprise($input['register_enterprise']);
            } else {
                $data['code'] = 'enterprise_not_found';
                $data['desc'] = '需要注册企业';

                // 需要注册企业给用户默认的企业信息
                $data['auth_enterprise'] = $login->ucEnterprise;

                return success($data);
            }
        }

        // 检查企业是否过试用期
        if ($login->isExpire === true) {
            $data['code'] = 'enterprise_expired';
            $data['desc'] = '企业已过试用期';

            return success($data);
        }

        // 检查是否绑定管理员 是:查出来管理员id，然后前端调get-token接口完成登录
        if ($login->admin === null) {
            if (isset($input['bind_account']) && !empty($input['bind_account'])) {
                // 验证验证码
                try {
                    (new Common())->checkSmsCode([
                        'event'  => 'word_wechat_bind_account',
                        'mobile' => $input['bind_account']['mobile'],
                        'code'   => $input['bind_account']['code']
                    ]);
                } catch (\Exception $exception) {
                    $data['code'] = 'verify_code_error';
                    $data['desc'] = '验证码错误';

                    return success($data);
                }

                // 绑定
                $login->bindAccount($input['bind_account']['mobile']);
            } else {
                $data['code'] = 'admin_not_found';
                $data['desc'] = '管理员不存在';
                $data['epId'] = $login->enterprise['id'];

                return success($data);
            }
        }

        // 调用get-token接口需要一个临时token
        $tmpToken = Crypt::encrypt(json_encode([
            'expire_time' => time() + 3600
        ], JSON_UNESCAPED_UNICODE));

        $data['admin'] = $login->admin;

        // 这个是前端需要admin_id字段
        $data['admin']['admin_id'] = $data['admin']['id'];
        $data['tmp_token']         = $tmpToken;

        return success($data);
    }
}
