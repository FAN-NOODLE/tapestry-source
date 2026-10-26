<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\MarketingPlan;
use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskResource;
use App\Models\MarketingPlanVisibleRange;
use App\Models\MemberBannerRelation;
use App\Models\StaffRight;
use App\Services\MarketingPlanService;
use App\Services\MemberService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Log;

class AutoCreateMarketingPlanPreviewPoster extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:create_marketing_plan_preview_poster';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '营销任务生成后台预览海报';

    public function handle()
    {
        Log::info($this->description . '--开始', ['time' => Carbon::now()]);

        //查询需要生成预览海报的素材
        $allResource = MarketingPlanTaskResource::query()
            ->where('type',MarketingPlanTaskResource::RESOURCE_TYPE_APPLET_BANNER)
            ->where('poster_url','')
            ->where('resource_url',"!=",'')
            ->where('applet_path',"!=",'')
            ->get(['id as resource_id','applet_path','enterprise_id','resource_url','sun_code'])
            ->toArray();
        $service = new MarketingPlanService();
        foreach ($allResource as $resource){
            $service->createPreviewPoster($resource);
        }

        Log::info($this->description . '--end', ['time' => Carbon::now()]);
    }
}
