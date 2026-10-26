<?php

namespace App\Console\Commands\Message;

use App\Models\MessageScene;
use Illuminate\Console\Command;

class SyncMessageScene extends Command
{

    protected $signature = 'sync_message_scene';

    protected $description = '消息发送场景更新';

    /**
     * @throws \Exception
     */
    public function handle()
    {
        $scenes = MessageScene::query()
            ->with('channels.settings.category:id,name')
            ->orderBy('id', 'asc')
            ->get();

        foreach ($scenes as $scene) {

            $this->info($scene['id'].'-【'.$scene['desc'].'】正在处理：');
            $channel = $scene['channels']->firstWhere('channel', 'mini_program');

            if (is_null($channel)) {
                continue;
            }

            $categoryIds = $channel['settings']->pluck('category_id')->sort()->values()->all();
            $categories  = $channel['settings']->pluck('category')->sortBy('id')->values()->all();

            MessageScene::query()
                ->where('id', $scene['id'])
                ->update([
                    'available_category_ids' => json_encode($categoryIds),
                    'available_categories'   => json_encode($categories, 320)
                ]);
        }
    }
}
