<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Services\MemberService;
use CrsSdk\Shop\StaffRight;
use Illuminate\Http\Request;

class AddController extends BaseController
{
    /**
     * 添加员工
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

        if (empty($params['memFace'])) {
            throw new Error(100000, '员工头像不能为空!');
        }

        $params['avatarMediaid'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $params['memFace']);

        unset($params['memFace']);

        $result = (new Member())->add($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}
