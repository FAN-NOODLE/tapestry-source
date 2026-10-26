<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\MemberExtend;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class UpdateMemberExtendRoleField extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-member-extend-role-field';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新员工扩展表字段role';

    public function handle()
    {
        $env = env('APP_ENV', 'prod');

        $dbLink = '';

        if ('dev' === $env) {
            $dbLink = 'crs_dev_boss.';
        }

        if ('prod' === $env) {
            $dbLink = 'crs_prod_boss.';
        }

        if ('test' === $env) {
            $dbLink = 'crs_prod_boss.';
        }

        if (!$dbLink) {
            exit('未设置数据库连接');
        }

        echo 'dbLink:' . $dbLink;

        $memEpIds       = MemberExtend::query()->pluck('enterprise_id')->toArray();
        $sql            = "SELECT `id`, `domain` from " . $dbLink . "enterprise where `domain` <> '' and deleted_at = 0";
        $enterpriseList = DB::select($sql);
        $domainList     = array_column($enterpriseList, 'domain', 'id');

        foreach ($domainList as $enterpriseId => $domain) {
            if (in_array($enterpriseId, $memEpIds)) {
                $where = [
                    'pageNum'  => 1,
                    'pageSize' => 10000,
                ];

                try {
                    $result = (new Member())->list($domain, 'contact', $where);
                    foreach ($result['list'] ?? [] as $member) {
                        MemberExtend::query()
                            ->where([
                                'mem_uid'       => $member['memUid'],
                                'enterprise_id' => $enterpriseId
                            ])
                            ->update([
                                'role' => $member['memRole']
                            ]);
                    }
                } catch (\Exception $e) {
                    Log::info('update-member-extend-where:', $where);
                    Log::info('update-member-extend-error:', ['message' => $e->getMessage(), 'line' => $e->getLine()]);
                }

            }
        }
    }

}
