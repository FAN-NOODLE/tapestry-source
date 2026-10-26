<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:48
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use Illuminate\Http\Request;

class QrCodeController extends BaseController
{
    /**
     * 员工二维码
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        $result = (new Member())->qrCode($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}
