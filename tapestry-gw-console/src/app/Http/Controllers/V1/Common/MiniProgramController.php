<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/12 18:54
 */

namespace App\Http\Controllers\V1\Common;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\MiniProgram;

class MiniProgramController extends BaseController
{
    public function index()
    {
        $shopDetail = (new MiniProgram())->detail([
            'enterprise_id' => $this->enterpriseId
        ]);

        $app['app_id'] = $shopDetail['data']['app_id'] ?? '';
        $app['name'] = $shopDetail['data']['name'] ?? '';

        return success([$app]);
    }
}