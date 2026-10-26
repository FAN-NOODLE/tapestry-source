<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:24
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\MemberExtend;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use App\Exceptions\Error;

class ContactWayQrCodeController extends BaseController
{
    /**
     * 员工联系我二维码
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 员工id
        $memUserId = $request->post('mem_user_id', '');
        $qyDomain      = $request->post('qy_domain', '');
        $ecwState      = $request->post('ecw_state', '');

        $params = [
            'memUserid' => $memUserId,
            'ecwScene' => 2, // 场景，1-在小程序中联系，2-通过二维码联系
            'ecwType' => 1, // 联系方式类型,1-单人, 2-多人
            'ecwSkipVerify' => 1, // 外部客户添加时是否无需验证，1为无需验证
            'ecwState' => $ecwState
        ];

        $result = (new Member())->getContactQrCode($qyDomain, $params);

        if (isset($result['code']) && $result['code'] == '41054') {
            throw new Error(41054, '该员工尚未加入企业');
        }

        return success(['data' => [
            'qr_code' => $result['qrCode'],
            'config_id' => $result['id']
        ]]);
    }
}
