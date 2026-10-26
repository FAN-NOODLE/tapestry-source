<?php


namespace App\Http\Controllers\App\V1\StaffCoupon;


use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\Member;
use CrsSdk\Boss\Auth;
use CrsSdk\Shop\StaffCoupon;
use CrsSdk\Shop\StaffExtSetting;

class PreviewShareController extends BaseController
{
    /**
     * @return array|void
     * @throws \Exception
     */
    public function index()
    {
        // 获取用户信息
        $memUser = (new Member())->detail(
            $this->qyDomain,
            $this->appIdentifier,
            ['memUid' => $this->memUid]
        );
        // 验证员工信息
        if (empty($memUser)) {
            error(200016);
        }

        // 获取有效的内购券信息
        $staffCouponList = [];
        $list = (new StaffCoupon())->staffCouponEffectiveList([
            'enterprise_id' => $this->enterpriseId
        ]);
        if (!empty($list)) {
            foreach ($list as $key => $val) {
                $staffCoupon       = [
                    'coupon_name'     => $val['coupon_name'],
                    'coupon_value'    => $val['coupon_value'],
                    'coupon_mininumn' => $val['coupon_mininumn']
                ];
                $staffCouponList[] = $staffCoupon;
            }
        }

        // 获取运营推广设置
        $staffExtSetting = (new StaffExtSetting())->staffExtSettingDetail([
            'enterprise_id' => $this->enterpriseId
        ]);

        if (empty($staffExtSetting['data'])) {
            error(720004);
        }

        // 获取小程序授权信息
        $auth_detail = (new Auth())->authDetail([
            'enterprise_id' => $this->enterpriseId,
            'auth_type'     => 1
        ]);

        return success([
            'mem_username' => $memUser['memUsername'],
            'mem_face'     => str_replace('http:', 'https:', $memUser['memFace']),
            'staff_coupon' => $staffCouponList,
            'app_logo'     => $auth_detail['data']['logo'],
            'app_name'     => $auth_detail['data']['name'],
            'bg_img_path'  => $staffExtSetting['data']['bg_img_path']
        ]);
    }
}
