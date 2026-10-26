<?php

use Illuminate\Database\Seeder;
use App\Models\EnterpriseOpenPlatform;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;

class DataOpenPlatformSeeder extends Seeder
{
    // 第三方小程序配置信息
    protected $extBuildData = [
        // 测试环境
        'dev'  => [
            'app_id'          => 'wxa9defa41d9769b5e',
            'host'            => 'shop-api-dev-crs.vchangyi.com',
            'work_host'       => 'work-api-dev-crs.vchangyi.com',
            'enterprise_hash' => '054edce1d4cec91c76001d968841c08f',
        ],
        // 灰度环境
        'test' => [
            'app_id'          => 'wxe11d047c889f9c83',
            'host'            => 'crs-api-test.vchangyi.com',
            'work_host'       => 'work-api-test-crs.vchangyi.com',
            'enterprise_hash' => 'ecd85280aa135bbd0108dd6aa424565a',
        ],
        // 生产环境
        'prod' => [
            'app_id'          => 'wx535abb91ba41bb3b',
            'host'            => 'shop-api-crs.vchangyi.com',
            'work_host'       => 'work-api-crs.vchangyi.com',
            'enterprise_hash' => '79066ecdb52334c1f6eb7795fcffd44e',
        ]
    ];
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        // 新增默认开放平台配置
        try {
            // 老开放平台处理
            $config = config('wechat.open_platform.default');
            $config['ext_build_json'] = $this->extBuildData[env("APP_ENV")];
            // 创建开放平台默认记录
            $openPlatform = EnterpriseOpenPlatform::query()->create($config);
            // 企业表关联默认开放id
            Enterprise::query()->where('deleted_at', '=', 0)->update([
                'open_platform_id' => $openPlatform->id
            ]);

            // 授权表更新build——json
            EnterpriseThirdPartAccount::query()->where([
                "platform" => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,
                "enterprise_type" => EnterpriseThirdPartAccount::RELEASE_STABLE,
                "status" => EnterpriseThirdPartAccount::STATUS_ACTIVE,
                "deleted_at" => 0,

            ])->update([
                'ext_build_json' => $config['ext_build_json']
            ]);
        } catch (\Exception $exception) {

        }

    }
}
