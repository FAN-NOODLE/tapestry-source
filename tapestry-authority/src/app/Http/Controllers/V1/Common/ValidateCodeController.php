<?php

namespace App\Http\Controllers\V1\Common;

use App\Http\Controllers\V1\BaseController;
use App\Library\ValidateCode\ValidateCode;
use Illuminate\Support\Facades\Redis;

class ValidateCodeController extends BaseController
{
    /**
     * 验证码
     *
     * @return array|void
     */
    public function index()
    {
        $obj     = new ValidateCode();
        $code    = $obj->getCode();
        $redis = Redis::connection();
        $redis->del(['validate_code:*']);
        $randKey = 'validate_code:'. uniqid();
        $redis->setex($randKey, 300, $code);
        $base64Img = $obj->doimg();
        return success([
            'code_key'   => $randKey,
            'base64_img' => $base64Img
        ]);
    }

}
