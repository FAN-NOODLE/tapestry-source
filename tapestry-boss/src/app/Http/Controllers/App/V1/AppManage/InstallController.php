<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\AppManage;

use App\Models\EnterpriseApp;
use App\Models\EnterpriseAppInstall;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;

class InstallController extends BaseController
{

    /**
     * 安装应用
     * @author 何岳龙 heyuelong@vchangyi.com
     * @date 2018/12/25 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id', 0);
        $appId        = $request->post('app_id', 0);

        // 企业ID不能为空
        if (empty($enterpriseId)) {
            return error(600001);
        }

        // 应用ID不能为空
        if (empty($appId)) {
            return error(203001);
        }

        // 查询应用ID是否存在
        $appData = EnterpriseApp::query()->where('id', $appId)->first(['id', 'identity']);

        if (empty($appData->id ?? '')) {
            return error(203003);
        }

        $where = [
            'enterprise_id' => $enterpriseId,
            'app_id' => $appId
        ];

        // 查询企业应用是否安装过
        $total = EnterpriseAppInstall::query()->where($where)->count();

        if ($total) {
            return error(203002);
        }

        try {
            // 应用标示
            $where['identity'] = $appData->identity;

            // 添加企业应用
            EnterpriseAppInstall::create($where);
        } catch (\Exception $e) {
            return error(203004);
        }

        return success([
            'data' => [],
        ]);
    }
}
