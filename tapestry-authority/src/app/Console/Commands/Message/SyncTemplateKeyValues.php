<?php

namespace App\Console\Commands\Message;

use App\Models\SubscribeMessageSetting;
use Illuminate\Console\Command;

class SyncTemplateKeyValues extends Command
{

    protected $signature = 'sync_template_key_values';

    protected $description = '小程序订阅消息setting表更新key_values';

    public function handle()
    {
        $path    = base_path('app/Modules/Message/Content/mini_message.json');
        $json    = file_get_contents($path);
        $metaArr = json_decode($json, 1);

        foreach ($metaArr as $item) {
            SubscribeMessageSetting::query()
                ->where('id', $item['id'])
                ->update(['key_values' => json_encode($item['key_values'], 320)]);
        }
    }
}
