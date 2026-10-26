<?php
/*
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Providers;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Facade;
use Illuminate\Support\Facades\Log;

class ThirdPartFacade extends Facade
{
    /**
     * 默认为 Server.
     *
     * @return string
     */
    public static function getFacadeAccessor()
    {
        return 'third_part';
    }

    /**
     * @param $enterpriseId
     * @param $miniProgramType
     *
     * @return \EasyWeChat\OpenPlatform\Authorizer\MiniProgram\Application
     * @throws Error
     */
    public static function miniProgram($enterpriseId, $miniProgramType = 1)
    {
        $app = app('third_part')->init($enterpriseId, $miniProgramType)->mini_program;
        if (!$app) {
            throw new Error(100000, '小程序初始化参数不存在！');
        }

        return $app;
    }
}
