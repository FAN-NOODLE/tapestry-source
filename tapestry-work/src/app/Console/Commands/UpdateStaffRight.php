<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\StaffRight;
use App\Services\MemberService;
use Illuminate\Console\Command;

class UpdateStaffRight extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-staff-right';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新权限人员';

    public function handle()
    {

        // 查询所有权限
        $list = StaffRight::query()->get()->toArray();

        if (empty($list)) {
            echo '无数据不能更新';
            return false;
        }

        $enterpriseIds = array_unique(array_column($list, 'enterprise_id'));

        $qyDomains = [];
        foreach ($enterpriseIds as $enterpriseId) {

            $params = [
                'enterprise_id' => $enterpriseId
            ];

            $detail = (new \CrsSdk\Boss\Enterprise())->detail($params);

            $qyDomain                 = $detail['data']['domain'] ?? 0;
            $qyDomains[$enterpriseId] = $qyDomain;
        }

        // 循环更新数据
        foreach ($list as $value) {

            // 如果是部门
            if ($value['obj_type'] == 1) {
                $ucParams = [
                    'dpId'                   => $value['obj_id'],
                    'departmentChildrenFlag' => 1
                ];

                $memberLIst = MemberService::memberList($qyDomains[$value['enterprise_id']], 'contact', $ucParams);

                $memberUids = array_column($memberLIst, 'memUid');
                if (!empty($memberUids)) {
                    $objectExtraId = implode(',', $memberUids);

                    StaffRight::query()->where('id', $value['id'])->update(['object_extra_id' => $objectExtraId]);
                }

                sleep(1);
            } else {
                StaffRight::query()->where('id', $value['id'])->update(['object_extra_id' => $value['obj_id']]);
            }
        }

        return true;
    }

}
