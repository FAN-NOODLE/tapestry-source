<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2018/12/28
 * Time: 11:38
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use Illuminate\Http\Request;

class UcCallbackController extends BaseController
{
    /**
     * token
     */
    const TOKEN = 'REDACTED_FOR_SOURCE_SHARING';

    /**
     * key
     */
    const KEY = 'wxqyEncodingAeskey';

    /**
     * @param Request $request
     * @return array
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/30 11:21 AM
     */
    public function index(Request $request)
    {
        $qyDomain = $request->post('qy_domain');

        if (empty($qyDomain)) {
            throw new Error(100000, 'qy_domain 字段不能为空');
        }

        $result['token'] = (new Enterprise())->entSettingGet($qyDomain, self::TOKEN);
        $result['key'] = (new Enterprise())->entSettingGet($qyDomain, self::KEY);

        return success(['data' => $result]);
    }
}