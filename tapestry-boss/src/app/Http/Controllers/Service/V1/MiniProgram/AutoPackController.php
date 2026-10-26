<?php
/**
 * 自动化打包
 * 提审
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use App\Exceptions\Error;
use App\Models\EnterpriseMiniProgramAutoPack;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\EnterpriseLogService;
use CrsSdk\Build\AutoPack;
use Illuminate\Support\Facades\Log;

class AutoPackController extends MiniBaseController
{

    // 好物圈开启
    private $enterpriseMallButtonOn = 1;

    public function index()
    {
        // 参数组装
        $params = $this->getParams();
        $enterpriseId = $params['enterprise_id'];
        $versionId = $params['version_id'];
        unset($params['enterprise_id'], $params['version_id']);
        $params['tag'] = '';

        $enterprise = switchAppHost($enterpriseId);

        // 查询发布tag
        $releaseTag = EnterpriseThirdPartAccount::query()->where([
                'enterprise_id' => $enterpriseId,
                'type'          => EnterpriseThirdPartAccount::TYPE_SHOP,
                'status'        => EnterpriseThirdPartAccount::STATUS_ACTIVE,
                'verify_status' => EnterpriseThirdPartAccount::VERIFIED
            ])->first(['release_tag', 'release_branch']);
        $branch = '';
        if (!empty($releaseTag)) {
            $params['tag'] = $releaseTag->release_tag ?? '';
            $branch = $releaseTag->release_branch ?? '';
        }

        // 添加默认值小程序未开启好物圈配置显示空数组
        $params['plugins'] = !empty($params['plugins']) ? $params['plugins'] : [];
        // 小程序不打包模块
        $params['no_modules'] = !empty($params['no_modules']) ? $params['no_modules'] : [];
        // 跳转小程序列表
        $params['jump_mini_app'] = !empty($params['jump_mini_app']) ? $params['jump_mini_app'] : [];

        // 回调地址  取env配置，没有默认dev
        $config = [
            'branch'   => $branch,
            'config'   => $params,
            'callback' => env('APP_HOST', 'http://boss-api-dev.crs.dev-test.vchangyi.com') . '/callback/auto-pack-audit'
        ];

        if ($enterprise['version'] != 1) {
            $config['repo'] = 'ssh://git@gitlab-ce.k8s.tools.vchangyi.com:32201/crs/crs-store-wxapp.git';
        }

        Log::info('boos自动打包返回数据：' . json_encode($config));
        EnterpriseLogService::debug('自动化打包并且提审：' . $enterpriseId, $config);
        try {
            $result = (new AutoPack())->packTask($config);
            $taskConfig = [
                'version'       => $params['version'],
                'version_id'    => $versionId,
                'enterprise_id' => $enterpriseId,
                'task_id'       => $result['_id'],
                'status'        => $result['status'],
                'task_config'   => json_encode($config),
                'statusDesc'    => $result['statusDesc'] ?? ''
            ];
            EnterpriseMiniProgramAutoPack::query()->create($taskConfig);
        } catch (\Exception $e) {
            Log::info('自动打包失败：' . $e->getMessage());
            Log::info('自动打包失败，参数：' . json_encode($config));

            EnterpriseLogService::error('自动打包失败：' . $enterpriseId, [
                'params'  => $config,
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine()
            ]);
            return new Error(100000, '打包失败：' . $e->getMessage());
        }

        return success($result);
    }

    public function getParams()
    {
        $params = $this->request->only([
            'style',
            'label_data',
            'pages_config',
            'no_modules',
            'plugins',
            'jump_mini_app',
            'ext',
            'ext_build',
            'version',
            'enterprise_id',
            'version_id',
            'goods_detail_mall_open',
            'sdk',
            'no_ex_code',
        ]);

        $rules = [
            'style'                     => 'required',
            'style.colors'              => 'required',
            'style.rgb_colors'          => 'required',
            //            'label_data'                    => 'required',
            //            'label_data.*.label_name'       => 'required',
            //            'label_data.*.icon_selected'    => 'required',
            //            'label_data.*.icon_no_selected' => 'required',
            //            'label_data.*.label_type'       => 'required',
            //            'label_data.*.sort'             => 'required',
            'pages_config'              => 'required',
            'ext'                       => 'required',
            'ext.app_name'              => 'required',
            'ext.app_id'                => 'required',
            'ext.host'                  => 'required',
            'ext.enterprise_hash'       => 'required',
            'ext_build'                 => 'required',
            'ext_build.app_id'          => 'required',
            'ext_build.host'            => 'required',
            'ext_build.enterprise_hash' => 'required',
            'version'                   => 'required',
            'enterprise_id'             => 'required',
            'version_id'                => 'required',
            'goods_detail_mall_open'    => 'required'
        ];

        $message = [
            'style.required'                     => '店铺风格不能为空',
            'style.colors.required'              => '店铺风格色值不能为空',
            'style.colors.rgb_colors.required'   => '店铺风格rgb色值不能为空',
            //            'label_data.required'                    => '店铺导航不能为空',
            //            'label_data.*.label_name.required'       => '店铺导航名称不能为空',
            //            'label_data.*.icon_selected.required'    => '店铺导航选中图不能为空',
            //            'label_data.*.icon_no_selected.required' => '店铺导航未选中图不能为空',
            //            'label_data.*.label_type.required'       => '店铺导航类型不能为空',
            //            'label_data.*.sort.required'             => '店铺导航排序不能为空',
            'pages_config.required'              => '模版数据不能为空',
            'ext.required'                       => '小程序配置不能为空',
            'ext.app_name.required'              => '小程序名称不能为空',
            'ext.app_id.required'                => '小程序app_id不能为空',
            'ext.host.required'                  => '小程序数据源不能为空',
            'ext.enterprise_hash.required'       => '小程序企业hash不能为空',
            'ext_build.required'                 => '第三方小程序配置不能为空',
            'ext_build.app_id.required'          => '第三方小程序app_id不能为空',
            'ext_build.host.required'            => '第三方小程序数据源不能为空',
            'ext_build.enterprise_hash.required' => '第三方小程序企业hash不能为空',
            'version.required'                   => '小程序版本号不能为空',
            'enterprise_id.required'             => '企业id不能为空',
            'version_id.required'                => '版本id不能为空',
            'goods_detail_mall_open.required'    => '好物推荐开关不能为空'
        ];

        $this->validate($this->request, $rules, $message);
        return $params;
    }
}
