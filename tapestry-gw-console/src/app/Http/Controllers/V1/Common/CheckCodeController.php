<?php


namespace App\Http\Controllers\V1\Common;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Common;

class CheckCodeController extends BaseController
{
    /**
     * 验证码验证
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'mobile',
            'code',
        ]);

        $input['client'] = Constant::CLIENT;

        (new Common())->CheckCode($input);

        return success();
    }
}
