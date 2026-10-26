<?php

namespace App\Console\Commands;

use App\Services\MemberService;
use App\Services\SourceServices;
use CrsSdk\Console\Media;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoSourceMediaIdNew extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:source_media_id_new';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时更新素材图片mediaId';

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
     * @return string
     * @throws \Exception
     */
    public function handle()
    {
        $gwConsole = config('client.api_url.console') . '/gw-console';
        config(['client.api_url.console' => $gwConsole]);

        $data = (new Media())->list(['type' => '1','size' => 100000]);

        (new SourceServices())->sysSourceMediaId($data['data']);

        return 'ok';
    }
}
