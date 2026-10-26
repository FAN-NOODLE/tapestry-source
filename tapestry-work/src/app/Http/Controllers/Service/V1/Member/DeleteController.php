<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:48
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\CashCouponMember;
use App\Models\ContactQrCodeParams;
use App\Models\GuideUserRelation;
use App\Models\MemberExtend;
use App\Models\StaffRight;
use App\Services\StaffSettingService;
use CrsSdk\Shop\StaffRight as StaffRightSDK;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class DeleteController extends BaseController
{
    /**
     * 删除员工
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        // 如果删除的是导购并且导购分配了客户就不允许删除
        $exists = GuideUserRelation::query()
            ->where([
                'mem_uid' => $params['memUid']
            ])
            ->exists();

        if ($exists) {
            throw new Error(10000, '该导购下存在已分配用户,禁止删除');
        }

        (new Member())->delete($qyDomain, $appIdentifier, $params);

        $memberExtendInfo = MemberExtend::query()
            ->where([
                'mem_uid' => $params['memUid']
            ])->first();
        if ($memberExtendInfo && $memberExtendInfo->config_id) {
            (new Member())->delContactQrCode($qyDomain, ['id' => $memberExtendInfo->config_id]);
        }

        // 更新联系我二维码关联员工/部门
        $contactParams = ContactQrCodeParams::query()
            ->where('object_extra', 'like', '%' . $params['memUid'] . '%')
            ->get()
            ->toArray();
        foreach ($contactParams as $v) {
            if ($v['type'] == ContactQrCodeParams::MEMBER) {
                ContactQrCodeParams::query()->where('id', $v['id'])->delete();
            } else {
                $memUids = explode(',', $v['object_extra']);
                $diffMemUids = array_diff($memUids, [$params['memUid']]);
                $diffMemUids = implode(',', $diffMemUids);
                ContactQrCodeParams::query()->where('id', $v['id'])->update(['object_extra' => $diffMemUids]);
            }
        }

        // 如果是离职,则删除本地member_extend
        MemberExtend::query()
            ->where([
                'mem_uid' => $params['memUid']
            ])
            ->delete();

        // 更新企业微信红包发送人
        CashCouponMember::query()
            ->where([
                'object' => $params['memUid']
            ])
            ->update([
                'is_leave' => CashCouponMember::MEMBER_IS_LEAVE
            ]);

        StaffRight::query()->where('obj_type', StaffRight::TYPE_STAFF)->where('obj_id', $params['memUid'])->delete();

        try {
            // 内购员工用户关系解绑
            (new StaffRightSDK())->staffUnbindUser(['mem_uid' => $params['memUid']]);
        } catch (\Exception $exception) {
            Log::error('内购员工用户关系解绑异常：'. $exception->getMessage());
        }

        try {
            // 内购员工分享关系解绑
            (new StaffSettingService())->unBindingStaffShareBymemUid([$params['memUid']]);
        } catch (\Exception $exception) {
            Log::error('内购员工分享关系解绑异常：'. $exception->getMessage());
        }

        return success(['data' => []]);
    }
}