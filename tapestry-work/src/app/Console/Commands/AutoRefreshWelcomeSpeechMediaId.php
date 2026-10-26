<?php

namespace App\Console\Commands;

use App\Models\WelcomeSpeech;
use App\Services\MemberService;
use CrsSdk\Boss\Enterprise;
use App\Library\Uc\Common;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoRefreshWelcomeSpeechMediaId extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:refresh_welcome_speech_media_id';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时任务更新欢迎语mediaId';

    /**
     * MemberService
     */
    protected $memberService;

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        $this->memberService = new MemberService();
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     *
     * @throws \Prettus\Repository\Exceptions\RepositoryException
     */
    public function handle()
    {
        $appIdentifier = 'contact';

        // 查询当前未结束的任务
        $welcomeSpeech = WelcomeSpeech::query()
            ->where('updated_at', '<', (time() - 2 * 23 * 60 * 60))
            ->whereIn('type', [WelcomeSpeech::IMAGE, WelcomeSpeech::MINI_PROGRAM])
            ->get()
            ->toArray();

        // 为空
        if (empty($welcomeSpeech)) {
            return;
        }

        foreach ($welcomeSpeech as $val) {
            if (empty($val['original_data'])) {
                continue;
            }
            $originalData = json_decode($val['original_data'], true);
            if (empty($originalData)) {
                continue;
            }
            $pic          = $originalData['pic'];
            $object       = '';
            $enterpriseId = $val['enterprise_id'];
            $enterprise   = (new Enterprise)->detail(['enterprise_id' => $enterpriseId]);
            $qyDomain     = $enterprise['data']['domain'] ?? '';
            if ($qyDomain) {
                if (WelcomeSpeech::IMAGE == $val['type']) {
                    $object = json_decode($val['object'], true);

                    $object['mediaId'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $pic);

                    $object = json_encode($object, JSON_UNESCAPED_UNICODE);
                }
                if (WelcomeSpeech::MINI_PROGRAM == $val['type']) {
                    $object = json_decode($val['object'], true);

                    $object['picMediaId'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $pic);

                    $object = json_encode($object, JSON_UNESCAPED_UNICODE);
                }
                if ($object) {
                    WelcomeSpeech::query()
                        ->where('id', $val['id'])
                        ->update([
                            'object' => $object,
                        ]);
                }
            }
        }
        return 'ok';
    }
}
