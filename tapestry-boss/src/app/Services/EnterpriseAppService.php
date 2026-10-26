<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:20
 */

namespace App\Services;

use App\Common\Constant;
use App\Models\EnterpriseApp;
use App\Models\EnterpriseAppInstall;

class EnterpriseAppService extends BaseService
{

    /**
     * 根据企业查询应用列表
     *
     * @param int $enterpriseId 企业ID
     *
     * @return array|void
     */
    public function enterpriseAppList($enterpriseId = 0)
    {
        // 如果查询应用列表
        $enterpriseAppList = EnterpriseApp::query()->get(['id', 'identity']);

        // 如果有默认应用
        if (!empty($enterpriseAppList)) {
            // 已安装应用列表
            $installList = EnterpriseAppInstall::query()
                ->where('enterprise_id', (int)$enterpriseId)
                ->get(['enterprise_id', 'app_id'])
                ->toArray();

            if (!empty($installList)) {
                $installList = array_column($installList, null, 'app_id');
            }

            // 查询应用安装状态
            foreach ($enterpriseAppList as $value) {
                // 默认未安装
                $status = Constant::ENTERPRISE_UNINSTALLED_STATUS;

                // 如果是优惠券则显示已经安装
                if ($value['identity'] == 'coupon') {
                    $value->status = Constant::ENTERPRISE_INSTALL_STATUS;
                } else {
                    // 判断如果已安装
                    if (!empty($installList[$value['id']] ?? '')) {
                        $status = Constant::ENTERPRISE_INSTALL_STATUS;
                    }

                    $value->status = $status;
                }

                $app = (new EnterpriseApp())->app;

                $value->title = isset($app[$value['identity']]['title']) ? $app[$value['identity']]['title'] : '';
                $value->content = isset($app[$value['identity']]['content']) ? $app[$value['identity']]['content'] : '';
            }
        }

        return $enterpriseAppList;
    }
}
