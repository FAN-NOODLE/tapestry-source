<?php


namespace App\Http\Controllers\App\V1\StaffCoupon;


use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\Member;
use CrsSdk\Boss\Auth;
use CrsSdk\Shop\StaffCoupon;
use CrsSdk\Shop\StaffExtSetting;

class ShareController extends BaseController
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

        // 生成分享页小程序太阳码
        $qr_code = (new StaffCoupon())->createQrCode([
            'mem_uid'       => $this->memUid,
            'path'          => 'pages/employee-friend-coupon/employee-friend-coupon',
            'enterprise_id' => $this->enterpriseId
        ]);

        return success([
            'mem_username'     => $memUser['memUsername'],
            'mem_face'         => str_replace('http:', 'https:', $memUser['memFace']),
            'app_logo'         => $auth_detail['data']['logo'],
            'app_name'         => $auth_detail['data']['name'],
            'staff_share_text' => $staffExtSetting['data']['staff_share_text'],
            'poster_img_path'  => $staffExtSetting['data']['poster_img_path'],
            'share_img_path'   => $staffExtSetting['data']['share_img_path'],
            'qr_code'          => $qr_code['data']
        ]);
    }
}
