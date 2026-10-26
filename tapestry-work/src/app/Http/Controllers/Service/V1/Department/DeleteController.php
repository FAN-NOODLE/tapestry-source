<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:31
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Department;
use App\Models\ContactQrCodeParams;
use Illuminate\Http\Request;
use App\Exceptions\Error;

class DeleteController extends BaseController
{
    /**
     * 删除部门
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

        // 查询联系我二维码关联部门
        $res = ContactQrCodeParams::query()->where('object', $params['dpId'])->where('type', ContactQrCodeParams::DEPARTMENT)->first();

        if ($res) {
            throw new Error(20020501, '该部门已被关联到联系我二维码中，请先编辑二维码取消设置后在删除');
        }

        (new Department())->delete($qyDomain, $appIdentifier, $params);

        return success(['data' => []]);
    }
}
