<?php

namespace App\Console\Commands\Message;

use App\Models\SubscribeMessageSetting;
use Illuminate\Console\Command;
use Illuminate\Support\Str;

class GetJsonForTemplateKeyValues extends Command
{

    protected $signature = 'get_json_for_template_key_values';

    protected $description = '获取小程序订阅消息setting表key_values字段json文件';

    public function handle()
    {
        $data = $this->getJson();
        $path = base_path('app/Modules/Message/Content/mini_message_' . date('mdHis') . '.json');

        file_put_contents($path, json_encode($data, 320));
        // 执行完成后 使用快捷键 (option + command + L) 可以格式化json （windows CTR + ALT + L）
    }

    public function getJson()
    {
        $settings = SubscribeMessageSetting::query()
            ->whereNotIn('scene_type', [6, 11, 17, 18])
            ->where('kids', '<>', '')
            ->whereIn('scene_type', [21, 22, 23, 24, 25, 26])
            ->orderBy('id', 'asc')
//            ->orderBy('scene_type', 'asc')
//            ->orderBy('category_id', 'asc')
            ->get();

        foreach ($settings as &$setting) {

            $setting['category_id'] = json_encode([
                'category_id' => $setting['category_id'],
                'tid'         => $setting['tid'],
                'scene_type'  => $setting['scene_type'],
                'scene_desc'  => $setting['scene_desc'],
            ], 320);
            $keyRules               = json_decode($setting['key_rules'], 1);
            $setting['key_values']  = json_decode($setting['key_values'], 1);
            if (empty($setting['key_values'])) {
                $values = [];
                foreach ($keyRules as $keyRule) {
                    $values[$keyRule] = [
                        'from'      => '',
                        'translate' => preg_replace('/[^a-z]/i', '', Str::camel($keyRule)),
                        'template'  => ''
                    ];
                }
                $setting['key_values'] = $values;
            }
            unset($setting['tid'], $setting['scene_type'], $setting['scene_desc']);
            unset($setting['type'], $setting['touch_ids'], $setting['page']);
            unset($setting['created_at'], $setting['updated_at']);
        }

        return $settings;
    }
}
