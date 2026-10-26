<?php

namespace App\Console\Commands;

use App\Models\MemberGuide;
use App\Models\User;
use App\Models\EnterpriseExtend;
use App\Models\GuideUserRelation;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class AutoAllocateUserToGuide extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:allocate_user_to_guide';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '每一分钟给客户分配导购';

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
        // 查询分配客户开关为开的企业
        $enabledAllocateEnterpriseIds = EnterpriseExtend::query()
            ->where('allocation_status', EnterpriseExtend::OPEN)
            ->pluck('enterprise_id')
            ->toArray();

        if ($enabledAllocateEnterpriseIds) {
            foreach ($enabledAllocateEnterpriseIds as $enterpriseId) {
                foreach (User::query()
                             ->where([
                                 'enterprise_id' => $enterpriseId,
                                 'is_allocated'  => User::IS_ALLOCATED_NO, // 未分配用户
                             ])
                             ->where('type', '!=', 0)
                             ->where('union_id', '!=', '')
                             ->where('mobile', '!=', '')
                             ->cursor() as $user) {

                    $guideUserRelation = GuideUserRelation::query()
                        ->where('union_id', $user['union_id'])
                        ->first(GuideUserRelation::$selectField);

                    if ($guideUserRelation && !empty($guideUserRelation['user_id'])) {
                        User::query()->where(['enterprise_id' => $enterpriseId, 'union_id' => $user['union_id']])->update(['is_allocated' => User::IS_ALLOCATED_YES]);
                        continue;
                    }

                    if ($guideUserRelation && empty($guideUserRelation['user_id'])) {
                        GuideUserRelation::query()->where('union_id', $user['union_id'])->update(['user_id' => $user['id']]);
                        User::query()->where(['enterprise_id' => $enterpriseId, 'union_id' => $user['union_id']])->update(['is_allocated' => User::IS_ALLOCATED_YES]);

                    } else {
                        // 查询分配状态为开启导购
                        $notReachCloseGuide = MemberGuide::query()
                            ->select(DB::raw('member_guide.mem_uid, member_guide.enterprise_id, member_guide.guide_id, member_guide.limit_user_count as limit_user_count, count(member_guide.mem_uid) count'))
                            ->leftJoin('guide_user_relation', function ($join) {
                                $join->on([
                                    ['member_guide.mem_uid', '=', 'guide_user_relation.mem_uid']
                                ])->where('guide_user_relation.deleted_at', 0);
                            })
                            ->where('member_guide.enterprise_id', $enterpriseId)
                            ->where('member_guide.allocation_status', MemberGuide::OPEN)
                            ->where('member_guide.deleted_at', 0)
                            ->groupBy('member_guide.mem_uid', 'member_guide.guide_id')
                            ->havingRaw('count(guide_user_relation.mem_uid) < limit_user_count')
                            ->get()
                            ->toArray();

                        $guide = $notReachCloseGuide;

                        if ($guide) {
                            $memIds = array_column($guide, 'mem_uid');

                            $memberCount = count($memIds);

                            // 查询出分配的最后一个导购
                            $latestGuideUserRelation = GuideUserRelation::query()
                                ->where('enterprise_id', $enterpriseId)
                                ->orderBy('id', 'DESC')
                                ->first(GuideUserRelation::$selectField);

                            $memberKey = 0;

                            if ($latestGuideUserRelation) {
                                $memIdsKey = array_search($latestGuideUserRelation->mem_uid, $memIds);

                                if ($memIdsKey !== false) {
                                    if ($memIdsKey + 1 != $memberCount) {
                                        $memberKey = $memIdsKey + 1;
                                    }
                                }
                            }

                            $guideUserRelationData = [
                                'enterprise_id' => $enterpriseId,
                                'guide_id'      => $guide[$memberKey]['guide_id'],
                                'mem_uid'       => $guide[$memberKey]['mem_uid'],
                                'user_id'       => $user['id'],
                                'union_id'      => $user['union_id'],
                                'relation_at'   => time(),
                                'created_at'    => time(),
                                'updated_at'    => time(),
                            ];

                            if ($guideUserRelationData) {
                                GuideUserRelation::query()->create($guideUserRelationData);
                                User::query()->where(['enterprise_id' => $enterpriseId, 'union_id' => $user['union_id']])->update(['is_allocated' => User::IS_ALLOCATED_YES]);
                            }
                        }
                    }
                }
            }
        }
    }
}
