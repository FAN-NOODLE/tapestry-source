<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Common\RedisKey;
use App\Models\Enterprise;
use Illuminate\Http\Request;
use CrsSdk\Shop\SubscribeMessage;
use Illuminate\Support\Facades\Log;
use App\Models\EnterpriseWorkPlugin;
use App\Library\SubscribeMessageSdk;
use App\Models\EnterpriseThirdPartAccount;
use App\Http\Controllers\Service\V1\BaseController;

class ShopSettingController extends BaseController
{
    /**
     * 保存商城小程序的secret
     * @param Request $request
     * @return array|void
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function index(Request $request)
    {
        $enterpriseId      = $request->post('enterprise_id');
        $miniProgramSecret = $request->post('mini_program_secret', '');
        $corpId            = $request->post('corp_id', '');
        $plPluginId        = $request->post('pl_plugin_id', '');
        $appSecret         = $request->post('app_secret', '');
        $agentId           = $request->post('agent_id', '');

        // 切换SHOP域名
        switchAppHost($enterpriseId);
        // 查询企业信息
        $enterprise = Enterprise::query()
            ->where('id', $enterpriseId)
            ->first();

        if (!empty($agentId) && !empty($miniProgramSecret)) {
            $hashLockKey = redisKey('ENTERPRISE_HASH_INFO', ['hash' => $enterprise->hash]);
            $domainLockKey = redisKey('ENTERPRISE_DOMAIN_INFO', ['domain' => $enterprise->domain]);

            // 企业信息编辑 .删除缓存
            delRedis([
                RedisKey::REDIS_APP_NOVICE . $enterprise->id,
                RedisKey::REDIS_APP_GROUP . $enterprise->id,
                RedisKey::REDIS_APP_SECKILL . $enterprise->id,
                RedisKey::REDIS_ENTERPRISE_INFO . $enterprise->hash,
                $hashLockKey,
                $domainLockKey
            ]);

            // 更新企业表contact_status的状态,2为禁用状态,3为自建应用权限
            Enterprise::query()
                ->where('id', $enterpriseId)
                ->update(['contact_status' => 3]);

            $where = [
                'enterprise_id' => $enterpriseId,
                'identifier'    => 'crs-shop'
            ];

            $update = [
                'enterprise_id' => $enterpriseId,
                'agent_id'      => $agentId,
                'secret'        => $miniProgramSecret,
                'corp_id'       => $corpId,
                'identifier'    => 'crs-shop',
                'pl_plugin_id'  => $plPluginId
            ];

            $query = EnterpriseWorkPlugin::query()->where($where)->first();

            if (!empty($query)) {
                $query->update($update);
                delRedis([RedisKey::ENTERPRISE_WORK_PLUGIN_BY_ENTERPRISE_ID . $enterpriseId]);
            } else {
                EnterpriseWorkPlugin::query()->create($update);
            }
        }

        // 查询企业app_id
        $enterpriseAccount = EnterpriseThirdPartAccount::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'type'          => EnterpriseThirdPartAccount::TYPE_SHOP,
                'status'        => EnterpriseThirdPartAccount::STATUS_ACTIVE,
            ])
            ->first();

        // 如果app_secret和app_id不为空
        if (!empty($appSecret) && !empty($enterpriseAccount->app_id)) {
            try {
                // 获取当前小程序分类
                $getCategoryData = (new SubscribeMessageSdk($enterpriseId, $enterpriseAccount->app_id, $appSecret))->getCategory();
            } catch (\Exception $exception) {
                return error(610001);
            }

            if ($getCategoryData['errcode'] == 0) {
                $wxCategoryIds = json_encode(array_column($getCategoryData['data'], 'id'));

                EnterpriseThirdPartAccount::query()->where('id', $enterpriseAccount->id)->update(
                    [
                        'wx_category_ids' => $wxCategoryIds,
                        'wx_category'     => json_encode($getCategoryData['data']),
                        'app_secret'      => $appSecret
                    ]
                );
                delRedis([RedisKey::MINI_PROGRAM_AUTH_INFO_BY_ENTERPRISE_ID . $enterpriseId]);
                // 组装同步模板参数
                $params = [
                    'app_id'        => $enterpriseAccount->app_id,
                    'app_secret'    => $appSecret,
                    'enterprise_id' => $enterpriseId
                ];
                Log::info('同步模板参数:' . json_encode($params));

                // 传递hash到模板同步
                $_SERVER['enterprise_hash'] = $enterprise->is_split_flow == 1 ? $enterprise->hash : 'public';
                try {
                    (new SubscribeMessage())->syncTemplates($params);
                } catch (\Exception $exception) {
                    Log::info('同步模板失败:' . json_encode([$exception->getCode(), $exception->getMessage()]));
                }
                Log::info('同步模板成功:');

            }
        }

        return success(['data' => []]);
    }
}
