<?php
/**
 * 批量添加员工成为导购
 */

namespace App\Console\Commands;


use App\Common\Constant;
use App\Jobs\SyncHistoryExternalJob;
use App\Library\Uc\External;
use App\Models\EnterpriseMembersToGuide;
use App\Services\StoreWorkerService;
use App\Library\Uc\Member;
use Illuminate\Console\Command;

class BatchAddMemberToGuide extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'batch_add_member_to_guide';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '批量添加员工成为导购';

    /**
     * 指定某些员工成为导购
     *
     * @throws \Exception
     */
    public function handle()
    {
        // 1、查询所有待导入的员工列表
        $waitAddMembers = EnterpriseMembersToGuide::query()->select(['domain', 'userid', 'enterprise_id'])->where(
            [
                ['is_exec', 0],
            ]
        )->get()->toArray();
        // 2、向uc查询员工信息
        foreach ($waitAddMembers as $waitMember) {
            $start_at = time();

            $err = "添加导购成功";
            try {
                $result = (new Member())->detailByUserId($waitMember['domain'], 'contact', ['memUserid' => $waitMember['userid']]);
                if ($result['memUid'] == "") {
                    continue;
                }
                // 3、把存在的员工插入到导购中
                $input['type'] = 3;
                $input['enterprise_id'] = $waitMember['enterprise_id'];
                $input['member'] = [
                    [
                        'mem_uid'   => $result['memUid'],
                        'user_name' => $result['memUsername'],
                        'mobile'    => $result['memMobile'],
                        'avatar'    => $result['memFace'],
                        'user_id'   => $result['memUserid']
                    ],
                ];
                $msg = StoreWorkerService::addGuide($input);
                $upWhere = [
                    'userid'        => $waitMember['userid'],
                    'enterprise_id' => $waitMember['enterprise_id']
                ];
                $is_exec = 1;
                $memUid = $result['memUid'];
                if (!is_bool($msg)) {
                    $err = $msg;
                    $is_exec = 2;
                }
                $end_at = time();
                $upData = [
                    'is_exec'  => $is_exec,
                    'result'   => $err,
                    'start_at' => $start_at,
                    'end_at'   => $end_at,
                    'mem_uid'  => $memUid
                ];
                EnterpriseMembersToGuide::query()
                    ->where($upWhere)
                    ->update($upData);
            }catch (\Exception $exception){
                $end_at =time();
                $err=$exception->getMessage();
                $upWhere = [
                    'userid'        => $waitMember['userid'],
                    'enterprise_id' => $waitMember['enterprise_id']
                ];
                $upData = [
                    'is_exec'  => 2,
                    'result'   => $err,
                    'start_at' => $start_at,
                    'end_at'   => $end_at,
                ];

                EnterpriseMembersToGuide::query()
                    ->where($upWhere)
                    ->update($upData);
            }
        }


        die("批量添加完成~~~");
    }
}