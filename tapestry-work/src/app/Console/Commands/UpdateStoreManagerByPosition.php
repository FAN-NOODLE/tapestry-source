<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Common\Constant;
use App\Models\Scrm\Departments;
use App\Models\Scrm\Members;
use App\Services\ScrmMemberService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class UpdateStoreManagerByPosition extends Command
{


    const UPDATE_INTERVAL_TIME = 60;

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-store-manager-by-position';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '根据用户企微职位字段更新用户管理权限';


    public function handle()
    {
        $enterpriseId = env('COACH_ENTERPRISE_ID');
        if (empty($enterpriseId)) {
            return;
        }
        $guideSettingKey = redisKey('UPDATE_STORE_MANAGER_BY_POSITION_TIME', ['enterprise_id' => $enterpriseId]);
        $lastUpdateTime = getOrSetData($guideSettingKey);
        if ((time() - $lastUpdateTime) < self::UPDATE_INTERVAL_TIME) {
            return;
        }
        //获取已经存在的部门
        $memberList = Members::query()->where('enterprise_id', $enterpriseId)
            ->where('status', Members::MEMBERS_STATUS_ACTIVATED)
            ->get(['mobile', 'position', 'id']);
        $scrmMemberSvc = new ScrmMemberService();
        foreach ($memberList as $val) {
            $position = $val->position;
            if (empty($position) || in_array($position, Members::COACH_STORE_BOSS)) {
                //如果是品牌BOSS或者不存在职位则不处理
                continue;
            }
            if (!in_array($position, array_merge(Members::COACH_STORE_BOSS, Members::COACH_STORE_SA, Members::COACH_STORE_SM))) {
                //如果不在系统设置的职位里也不处理
                continue;
            }
            if (empty($val->mobile)) {
                continue;
            }

            try {
                $updateParams = [
                    'enterprise_id' => $enterpriseId,
                    'mobile' => $val->mobile
                ];
                if (in_array($position, Members::COACH_STORE_SA)) {
                    //普通导购更新权限
                    $updateParams['type'] = Constant::STORE_WORKER_TYPE_ENTERPRISE;
                } else if (in_array($position, Members::COACH_STORE_SM)) {
                    //获取导购所在部门
                    $val->load('departmentsRel.departments');
                    $memberArr = $val->toArray();
                    $departmentIds = array_column(array_column($memberArr['departments_rel'], 'departments'), 'third_id');
                    $updateParams['type'] = Constant::STORE_WORKER_TYPE_MANAGER;
                    foreach ($departmentIds as $dpId) {
                        $updateParams['range'][] = [
                            'object_id' => $dpId,
                            'type'=> Constant::STORE_WORKER_TYPE_MANAGER
                        ];
                    }

                }
                if (!empty($updateParams['type'])) {
                    //更新用户身份及管理范围
                    $scrmMemberSvc->upgradeMemberRange($updateParams);
                }
            } catch (\Exception $e) {
                Log::error('根据职位更新用户(' . $val->mobile . ')管理权限失败:' . $e->getMessage() . '，file:' . $e->getFile() . $e->getLine());
            }

        }
        getOrSetData($guideSettingKey, time(), 24 * 3600 * 3);

    }


}
