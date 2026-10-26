<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/8/31
 * Time: 16:21
 */

namespace App\Http\Controllers\Callback;

use App\Exceptions\Error;
use App\Models\Enterprise;
use Illuminate\Http\Request;

class WorkPluginController
{
    // 安装企业微信小程序跳转地址
    const INSTALL_APP_JUMP_URL = '/myapp/setting/basic-setup/authorize/authorizes';

    public function index(Request $request)
    {
        $enterpriseHash = $request->input('hash', '');

        if (empty(trim($enterpriseHash))) {
           exit('企业参数[hash]缺失');
        }

        $enterprise = Enterprise::query()->where('hash', $enterpriseHash)->first();
        if ($enterprise && $enterprise->is_split_flow == 0) {
            // 不分流企业
            if ($enterprise->version == 1) {
                $enterpriseHash = 'public';
            } else {
                $enterpriseHash = 'store';
            }
        }

        return redirect(config('app.console_host') . '/' . $enterpriseHash . self::INSTALL_APP_JUMP_URL);
    }
}
