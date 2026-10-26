<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Micros\Common;
use App\Models\Scrm\Members;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class CheckController extends BaseController
{
    /**
     * 员工身份校验
     *
     * @param Request $request
     *
     * @return array
     * @throws Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        $result = (new Member())->check($qyDomain, $appIdentifier, $params);

        return success($result);
    }

    /**
     * 检查用户是否有上传素材的权限
     * @param Request $request
     * @return JsonResponse
     */
    public function CheckSource(Request $request): JsonResponse
    {
        $mediaId      = 393;
        $params       = $request->all();
        $auth         = 0;
        $userid       = $params['userid'];
        $enterpriseId = $params['enterprise_id'];
        $memberInfo   = Members::query()->where(['enterprise_id' => $enterpriseId, "userid" => $userid])->first();
        if (!empty($memberInfo) && !empty($memberInfo['mobile'])) {
            $input = [
                "enterprise_id" => $params['enterprise_id'],
                "mobile"        => $memberInfo['mobile']
            ];
            try {
                $result      = Common::Authority()->setApi("micro-authority/v1/account/account-permissions")->post($input);
                $permissions = $result['permissions'];
                if (in_array($mediaId, $permissions)) {
                    $auth = 1;
                }
            } catch (Exception $e) {
                error(10002, $e->getMessage());
            }
        }
        $res = ['auth' => $auth];
        return successJson($res);
    }
}
