<?php
/*
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Services\OpenPlatform;

use App\Common\RedisKey;
use App\Exceptions\Error;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\EnterpriseVirtualOpenPlatform;
use App\Services\BaseService;
use App\Services\EnterpriseNew\EnterpriseNewService;
use Illuminate\Support\Facades\Log;

class EnterpriseService extends BaseService
{
    /**
     * @var EnterpriseVirtualOpenPlatform
     */
    protected $virtualOpenPlatform;

    /**
     * @var EnterpriseThirdPartAccount
     */
    protected $thirdPartAccount;

    public function __construct(
        EnterpriseVirtualOpenPlatform $virtualOpenPlatform,
        EnterpriseThirdPartAccount $thirdPartAccount
    )
    {
        $this->thirdPartAccount = $thirdPartAccount;
        $this->virtualOpenPlatform = $virtualOpenPlatform;
    }

    /**
     * 更新第三方平台信息.
     *
     * @param array $config ['enterprise_id', 'open_app_id']
     *
     * @return mixed
     *
     * @throws Error
     */
    public function updateOrCreateEnterpriseConfig($config)
    {
        if (!key_exists('enterprise_id', $config)) {
            throw new Error(100000, '缺少参数enterprise_id');
        }

        if (!key_exists('open_app_id', $config)) {
            throw new Error(100000, '缺少参数open_app_id');
        }

        $where = ['enterprise_id' => $config['enterprise_id']];

        $query = $this->virtualOpenPlatform->where($where)->first();

        if (!empty($query)) {
            $query->update($config);
        } else {
            $this->virtualOpenPlatform->create($config);
        }

        return true;
    }

    /**
     * 更新第三方配置.
     *
     * @param array $config ['enterprise_id', 'platform', 'app_id']
     *
     * @return mixed
     *
     * @throws Error
     */
    public function updateOrCreateThirdPartConfig($config)
    {
        if (!key_exists('enterprise_id', $config)) {
            throw new Error(100000, '缺少参数enterprise_id');
        }

        $where = [
            'enterprise_id' => $config['enterprise_id'],
            'app_id'        => $config['app_id'],
            'platform'      => $config['platform'],
        ];

        $query = $this->thirdPartAccount->where($where)->first();
        // 允许多个小程序授权
        if (!empty($query)) {
            $query->update($config);
            $result = $query;
            // 清除授权信息缓存
            (new EnterpriseNewService())->clearAuthInfoCache($config['enterprise_id']);
        } else {
            $result = $this->thirdPartAccount->create($config);
        }

        try {
            (new \CrsSdk\GoMicroScrm\MiniProgram())->authorize([
                "enterprise_id" => $config['enterprise_id'],
                "app_name"      => $config['name'],
                "app_id"        => $config['app_id'],
                "logo"          => $config['logo'],
                "origin_id"     => $config['origin_id'],
            ]);
        } catch (\Exception $exception) {
            Log::error('授权信息传递给SCRM异常', [
                'line'    => $exception->getLine(),
                'code'    => $exception->getCode(),
                'file'    => $exception->getFile(),
                'message' => $exception->getMessage(),
            ]);
        }

        return $result;
    }

    /**
     * 取消公众号/小程序的授权.
     * 当 公众号/小程序 从未授权过，则会返回false.
     *
     * @param string $appId
     *
     * @return mixed
     */
    public function negativeThirdPartStatus($appId)
    {
        $thirdPartConfig = $this->thirdPartAccount
            ->query()
            ->where('app_id', '=', $appId)
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->first();
        if ($thirdPartConfig === null) {
            return false;
        }
        $this->thirdPartAccount
            ->query()
            ->where('app_id', $appId)
            ->update(['status' => EnterpriseThirdPartAccount::STATUS_NEGATIVE]);

        // 清除授权信息缓存
        (new EnterpriseNewService())->clearAuthInfoCache($thirdPartConfig->getAttribute('enterprise_id'));

        return true;
    }
}
