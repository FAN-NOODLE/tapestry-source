<?php

namespace App\Services;

use App\Common\Constant;
use App\Jobs\StoreWorkerGenerateQrCodeJob;
use App\Models\Scrm\Members;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\MicroAuthority\Role;
use Exception;
use Illuminate\Database\Concerns\BuildsQueries;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\DB;


class ScrmMemberService
{


    /**
     * @param $params
     * @return BuildsQueries|Builder|Model|mixed|object|null
     */
    public static function detail($params)
    {
        $res = Members::query()->when(!empty($params['enterprise_id']), function ($query) use ($params) {
            return $query->where('enterprise_id', '=', $params['enterprise_id']);
        })
            ->when(!empty($params['mobile']), function ($query) use ($params) {
                return $query->where('mobile', '=', $params['mobile']);
            })
            ->when(!empty($params['userid']), function ($query) use ($params) {
                return $query->where('userid', $params['userid']);
            })
            ->first();
        if (isset($res->qr_poster_code_id) && $res->qr_poster_code_id == 0) {

            // 获取domain
            $domain = (new EnterpriseService)->detail($params['enterprise_id'])['data']['domain'];

            // 异步处理生成导购海报
            $data = [
                'worker_ids'    => [$res->id],
                'enterprise_id' => $params['enterprise_id'],
                'store_id'      => 0,
                'domain'        => $domain,
                'operator_id'   => $input['operator_id'] ?? 0,
                'operator_name' => $input['operator_name'] ?? '系统'
            ];
            dispatch(new StoreWorkerGenerateQrCodeJob($data));

        }
        return $res ? $res->toArray() : [];
    }

    /**
     * 根据条件获取人员列表
     *
     * @param array $params
     * @return boolean
     * @throws Exception
     */
    public function upgradeMemberRange(array $params = []): bool
    {
        $mobile       = $params['mobile'];
        $enterpriseId = $params['enterprise_id'];
        $type         = $params['type'];

        try {
            DB::beginTransaction();

            $memberInfo = Members::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('mobile', $mobile)
                ->first();
            if ($memberInfo->type == $params['type'] && $type == Constant::STORE_WORKER_TYPE_ENTERPRISE) {
                //如果是普通导购，身份未变则跳过
                return true;
            }
            $memberInfo->type = $params['type'];
            $memberInfo->save();

            // 通过手机号查询账号信息
            $account       = (new Account());
            $accountDetail = $account->detailAccount([
                'enterprise_id' => $enterpriseId,
                'mobile'        => $mobile,
            ]);

            // 修改账号范围
            $updateRequest = [
                'id'            => $accountDetail['id'],
                'account'       => $accountDetail['account'],
                'username'      => $accountDetail['username'],
                'mobile'        => $accountDetail['mobile'],
                'enterprise_id' => $params['enterprise_id'],
                'range_type'    => 2,
                'range'         => []
            ];

            if ($type == Constant::STORE_WORKER_TYPE_MANAGER) {
                $updateRequest['range'] = $params['range'];
            }

            // 获取当前账号已存在的角色
            $accountRoleList    = $account->detailAccount([
                'id'            => $accountDetail['id'],
                'enterprise_id' => $params['enterprise_id']
            ]);
            $accountExitRoleIds = array_column($accountRoleList['new_role_list'], 'id');

            // 查询所有角色
            $roleList = (new Role())->lists([
                'page'          => 1,
                'size'          => 100,
                'enterprise_id' => $enterpriseId,
            ]);

            $markIndexArr = array_column($roleList['list'], "id", 'mark');

            // 企业管理
            if ($type == Constant::STORE_WORKER_TYPE_MANAGER) {
                // 同时拥有后台管理系统和智慧导购小程序权限
                $updateRequest['admin_available_system'] = [1, 2];
                $roleIds                                 = [
                    $markIndexArr[Constant::GUIDE_ENTERPRISE_GUIDE],
                    $markIndexArr[Constant::GUIDE_ENTERPRISE_ADMIN],
                    $markIndexArr[Constant::CONSOLE_ENTERPRISE_ADMIN]
                ];

                $updateRequest['role_ids'] = array_unique(array_merge($accountExitRoleIds, $roleIds));
            }

            // 企业导购
            if ($type == Constant::STORE_WORKER_TYPE_ENTERPRISE) {
                // 仅有智慧导购小程序权限
                $updateRequest['admin_available_system'] = [2];
                $roleIds                                 = [
                    $markIndexArr[Constant::GUIDE_ENTERPRISE_GUIDE],
                ];

                // 过滤掉数组中指定的角色
                $exitRoleIds = array_filter($accountExitRoleIds, function ($value) use ($markIndexArr) {
                    $enterpriseAdminRole = [
                        $markIndexArr[Constant::GUIDE_ENTERPRISE_ADMIN],
                        $markIndexArr[Constant::CONSOLE_ENTERPRISE_ADMIN]
                    ];
                    return !in_array($value, $enterpriseAdminRole);
                });

                $updateRequest['role_ids'] = array_unique(array_merge($exitRoleIds, $roleIds));
            }

            $account->editAccount($updateRequest);
            DB::commit();

        } catch (Exception $e) {
            DB::rollBack();
            throw $e;
        }

        return true;
    }
}
