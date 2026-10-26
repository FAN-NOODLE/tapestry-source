<?php

namespace App\Console\Commands\Message;

use App\Models\SubscribeMessageSetting;
use Illuminate\Console\Command;
use Illuminate\Support\Arr;
use Illuminate\Support\Str;

class SyncTemplateKeyRules extends Command
{

    protected $signature = 'sync_template_key_rules';

    protected $description = '小程序订阅消息setting表更新key_rules';

    public function handle()
    {
        $settings = SubscribeMessageSetting::query()
            ->where('key_rules', '')
            ->with('belongsTid:tid,rules')
            ->get();
        foreach ($settings as $setting) {
            $settingsConfigRules = explode("|", $setting['belongsTid']['rules']);

            $data = [];
            $kids = json_decode($setting['kids'], 1);
            foreach ($kids as $item) {
                $data[] = Arr::first($settingsConfigRules, function ($value) use ($item) {
                    return Str::endsWith($value, $item);
                });
            }
            SubscribeMessageSetting::query()
                ->where('id', $setting['id'])
                ->update(['key_rules' => json_encode($data, 320)]);
        }
    }
}
