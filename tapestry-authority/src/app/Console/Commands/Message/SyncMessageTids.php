<?php

namespace App\Console\Commands\Message;

use App\Models\EnterpriseThirdPartAccount;
use App\Models\SubscribeMessageTid;
use EasyWeChat\Factory;
use Illuminate\Console\Command;

class SyncMessageTids extends Command
{

    protected $signature = 'sync_message_tids';

    protected $description = '同步订阅消息关键字';

    public function handle()
    {
        $enterprises = EnterpriseThirdPartAccount::query()
            ->where('wx_category_ids', '<>', '')
            ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('platform', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('type', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->orderBy('enterprise_id', 'asc')
            ->get();

        $config = config('wechat.defaults');

        foreach ($enterprises as $enterprise) {
            $this->info('企业：' . $enterprise['enterprise_id'] . '-' . $enterprise['name'] . '正在处理');
            $config = array_merge($config, ['app_id' => $enterprise['app_id'], 'secret' => $enterprise['app_secret']]);

            $app = Factory::miniProgram($config);

            $categoryIds = json_decode($enterprise['wx_category_ids'], 1);

            foreach ($categoryIds as $categoryId) {
                $this->insertTid($app, $categoryId);
                $this->updateTid($app, $categoryId);
            }
        }
    }

    /**
     * @param $app
     * @param $categoryId
     */
    public function insertTid($app, $categoryId)
    {
        $categoryTidCount = SubscribeMessageTid::query()->where('categoryId', $categoryId)->count();
        if ($categoryTidCount) {
            $this->info('分类：' . $categoryId . '已同步，无需重复同步，tid数：' . $categoryTidCount);
            return;
        }
        $flag = 1;
        $i    = 0;
        $tids = [];
        while ($flag) {
            $limit = 30;
            $start = 0 + $limit * $i;

            try{
                $res = $app->subscribe_message->getTemplateTitles([$categoryId], $start, $limit);
            }catch(\Exception $exception){
                $res = [
                    'errcode'=>40125,
                    'errmsg'=>'invalid appsecret'
                ];
            }

            if (0 <> $res['errcode']) {
                $this->error('分类：' . $categoryId . '同步tid Error:' . $res['errmsg'] . '-' . $res['errmsg']);
                $flag = 0;
                continue;
            }
            if ($res['count'] < $limit * $i) {
                $flag = 0;
            } else {
                $tids[] = $res['data'];
                $i++;
            }
        }
        $tidRes = collect($tids)->collapse()->unique()->sort()->values()->all();

        SubscribeMessageTid::query()->insert($tidRes);

    }

    /**
     * @param $app
     * @param $categoryId
     */
    public function updateTid($app, $categoryId)
    {
        $categoryTids = SubscribeMessageTid::query()
            ->where('categoryId', $categoryId)
            ->where('keywords', '')
            ->get();
        if ($categoryTids->isEmpty()) {
            $this->info('分类：' . $categoryId . 'keywords已同步，无需重复同步');
            return;
        }
        foreach ($categoryTids as $tid) {
            $key = $app->subscribe_message->getTemplateKeywords($tid['tid']);
            if (0 <> $key['errcode']) {
                $this->error('分类：' . $categoryId . '-Tid：' . $tid['tid'] . '-同步keywords Error:' . $key['errmsg'] . '-' . $key['errmsg']);
                return;
            }
            $keyData = $key['data'] ?? [];
            array_walk($keyData, function (&$value) use ($tid) {
                $value['tid']       = $tid['tid'];
                $value['key_rules'] = $value['rule'] . $value['kid'];
            });
            $keywords = implode('|', array_column($keyData, 'name'));
            $rules    = implode('|', array_column($keyData, 'key_rules'));
            $this->info('分类：' . $categoryId .'-Tid：' . $tid['tid'] . '-正在同步keywords');
            SubscribeMessageTid::query()->where('tid', $tid['tid'])
                ->update([
                    'keywords' => $keywords,
                    'rules'    => $rules,
                ]);
        }
    }
}
