<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\AppManage;

use App\Common\Constant;
use App\Models\EnterpriseAppInstall;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class StatusController extends BaseController
{
    /**
     * 应用列表
     *
     * @author 何岳龙 heyuelong@vchangyi.com
     * @date   2018/12/25 11:16
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        return success([
            'data' => [
                'status' => Constant::ENTERPRISE_INSTALL_STATUS
            ]
        ]);

        $params = $request->only(['enterprise_id', 'identity']);

        // 企业ID不能为空
        if (empty($params['enterprise_id'])) {
            return error(600001);
        }

        // 应用标示不能为空
        if (empty($params['identity'])) {
            return error(203005);
        }

        // 应用标示不能为空
        if (!in_array($params['identity'], Constant::APP_IDENTITY)) {
            return error(203006);
        }

        // 默认未安装
        $status = Constant::ENTERPRISE_UNINSTALLED_STATUS;

        // 如果是优惠券则已安装
        if ('coupon' == $params['identity']) {
            $status = Constant::ENTERPRISE_INSTALL_STATUS;
        } else {
            $total = EnterpriseAppInstall::query()->where($params)->count();
            // 判断如果已安装
            if (!empty($total)) {
                $status = Constant::ENTERPRISE_INSTALL_STATUS;
            }
        }

        return success([
            'data' => [
                'status' => $status
            ]
        ]);
    }
}
