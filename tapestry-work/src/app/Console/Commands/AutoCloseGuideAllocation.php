<?php

namespace App\Console\Commands;

use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class AutoCloseGuideAllocation extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:close_guide_allocation';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '每一分钟关闭已经分配够足额的导购的分配权限';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     *
     * @throws \Exception
     */
    public function handle()
    {
        // 关闭分配状态为1、开启，分配客户超过分配上限的导购的分配权限

        // 查询分配状态为1的导购
        $isOpenMember = MemberGuide::query()
            ->where('member_guide.allocation_status', MemberGuide::OPEN)
            ->get()
            ->toArray();

        if ($isOpenMember) {
            $isOpenMemberUids = array_column($isOpenMember, 'mem_uid');

            // 查询分配状态为1的导购下已经分配的数量
            $isOpenMemberGuideUserRelation = GuideUserRelation::query()
                ->select(DB::raw('mem_uid, COUNT(mem_uid) as count'))
                ->whereIn('mem_uid', $isOpenMemberUids)
                ->groupBy('mem_uid')
                ->get()
                ->toArray();

            $isOpenMemberGuideUserRelation = array_column($isOpenMemberGuideUserRelation, null, 'mem_uid');

            // 初始化分配客户超过分配上限的导购ID集合
            $closeMemberIds = [];

            foreach ($isOpenMember as $value) {
                if (!isset($isOpenMemberGuideUserRelation[$value['mem_uid']])) {
                    continue;
                }
                if ($isOpenMemberGuideUserRelation[$value['mem_uid']]['count'] >= $value['limit_user_count']) {
                    $closeMemberIds[] = $value['mem_uid'];
                }
            }

            if ($closeMemberIds) {
                MemberGuide::query()->whereIn('mem_uid', $closeMemberIds)->update(['allocation_status' => MemberGuide::CLOSE]);
            }
        }
    }
}
