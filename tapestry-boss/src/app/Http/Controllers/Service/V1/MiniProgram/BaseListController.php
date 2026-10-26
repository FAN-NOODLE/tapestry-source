<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:10
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use App\Models\EnterpriseThirdPartAccount;
use App\Services\OpenPlatform\MiniProgramService;

class BaseListController extends MiniBaseController
{

    /**
     * 小程序基础配置信息列表
     * @param MiniProgramService $miniProgramService
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();

        $list = EnterpriseThirdPartAccount::query()
            ->where('enterprise_third_part_account.app_id', '<>', '')
            ->where('enterprise_third_part_account.app_secret', '<>', '')
            ->where('enterprise_third_part_account.status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('enterprise_third_part_account.platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('wx_category_ids', '<>', '')
            ->join('enterprise', 'enterprise.id', '=', 'enterprise_third_part_account.enterprise_id')
            ->when('shop' == ($params['version'] ?? ''), function ($query) {
                // 老版
                return $query->where('enterprise.version', 1);
            })
            ->when('gw-shop' == ($params['version'] ?? ''), function ($query) {
                // 老版
                return $query->where('enterprise.version', '<>', 1);
            })
            ->orderBy('enterprise_third_part_account.enterprise_id', 'asc')
            ->get([
                'enterprise_third_part_account.enterprise_id',
                'enterprise_third_part_account.app_id',
                'enterprise_third_part_account.app_secret',
                'enterprise_third_part_account.wx_category',
                'enterprise_third_part_account.wx_category_ids',
                'enterprise.version'
            ])
            ->toArray();

        return [
            'count' => count($list),
            'data'  => $list
        ];
    }
}
