<?php


namespace App\Console\Commands;


use App\Library\Util\Logger;
use App\Micros\Common;
use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use Illuminate\Support\Facades\DB;

class ForbidLeaveAdmin
{
    protected $signature = 'forbid-leave-admin';

    protected $description = '禁用离职管理员';

    public function handle()
    {
        self::listEnterpriseAdmin();
    }


    private static function listEnterpriseAdmin()
    {
        $admins = EnterpriseAdmin::query()
            ->where('is_super', EnterpriseAdmin::NO_SUPER)
            ->where('status', 1)
            ->get()
            ->toArray();

        $epMap = [];
        foreach ($admins as $admin) {
            $epMap[$admin['enterprise_id']][] = $admin['mobile'];
        }
        foreach ($epMap as $epId => $mobiles) {
            self::memberInfo($epId, $mobiles);
        }
    }

    private static function memberInfo($epId, $mobiles)
    {
        $qyDomain = Enterprise::query()->find($epId)->value('domain');

        $info = Common::work()
            ->setApi("/service/v1/admin/list")
            ->post(["enterprise_id" => $epId, "qy_domain" => $qyDomain, 'mobiles' => $mobiles]);
        if($info['list']){
            foreach ($info['list'] as $member) {
                EnterpriseAdmin::query()
                    ->where('enterprise_id', $epId)
                    ->where('mobile', $member['memMobile'])
                    ->update(['status' => 2]);
            }
        }

    }
}