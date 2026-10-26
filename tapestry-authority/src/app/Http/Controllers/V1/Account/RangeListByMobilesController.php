<?php


namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;

/**
 * 根据多个手机号获取管理范围列表
 */
class RangeListByMobilesController extends BaseController
{

    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id',
            'mobiles'
        ]);
        $adminList = EnterpriseAdmin::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('status', '=', EnterpriseAdmin::STATUS_OPEN)
            ->whereIn('mobile',$input['mobiles'])
            ->get(['id','mobile'])
            ->keyBy('id')
            ->toArray();
        $rangeList = EnterpriseDataRange::query()->whereIn('enterprise_admin_id', array_keys($adminList))
            ->where('type', EnterpriseDataRange::TYPE_DEPARTMENT)
            ->get(['enterprise_admin_id','object_id'])
            ->toArray();
        $haveDataArr = [];
        foreach ($rangeList as $val) {
            if (!empty($haveDataArr[$adminList[$val['enterprise_admin_id']]['mobile']])) {
                array_push($haveDataArr[$adminList[$val['enterprise_admin_id']]['mobile']], $val['object_id']);
            } else {
                $haveDataArr[$adminList[$val['enterprise_admin_id']]['mobile']] = [$val['object_id']];
            }
        }
        foreach ($input['mobiles'] as $mobile) {
            if (!in_array($mobile, array_keys($haveDataArr))) {
                $haveDataArr[$mobile] = [];
            }
        }

        return success($haveDataArr);
    }
}