<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\Member;
use App\Services\DepartmentService;
use CrsSdk\Boss\Auth;
use Illuminate\Support\Facades\Log;
use SimpleSoftwareIO\QrCode\Facades\QrCode;

class MeController extends BaseController
{
    /**
     * 获取当前登录用户详情
     *
     * @param DepartmentService $departmentService
     *
     * @return array|mixed
     * @throws \Exception
     */
    public function index(DepartmentService $departmentService)
    {
        // 调用uc方法获取数据
        $ucResult = (new Member())->detail($this->qyDomain, $this->appIdentifier, ['memUid' => $this->memUid]);

        if (!isset($ucResult)) {
            error(600103);
        }

        $dpId = '';
        // 部门类型id
        $dptId = 0;
        if (isset($ucResult['dpName'])) {
            // 取出所属部门id
            $dpId  = $ucResult['dpName'][0]['dpId'];
            $dptId = $ucResult['dpName'][0]['dptId'];
        }

        // 默认不是门店类型
        $isStore = 2;

        if ($dptId) {
            // 查询门店类型的部门类型id
            $storeDpTid = $departmentService->getStoreDptId($this->qyDomain, $this->appIdentifier);
            if ($storeDpTid == $dptId) {
                // 是门店类型
                $isStore = 1;
            }
        }

        if (empty($dpId)) {
            // 部门信息不存在
           error(600102);
        }

        $data = (new Member())->qrCode($this->qyDomain, $this->appIdentifier, ['memUid' => $this->memUid]);


        $shopData = (new Auth())->getShopAppid(['enterprise_id' => $this->enterpriseId])['data'];

        // 员工角色
        $role = '';
        if (isset($ucResult['role']) && !empty($ucResult['role'])) {
            $role = $ucResult['role']['roleName'];
        }

        // 赋值微信头像
        $memFace = $ucResult['memFace'];

        return success([
            'data' => [
                'mem_uid'         => $this->memUid,
                'mem_username'    => $ucResult['memUsername'],
                'mem_face'        => stripslashes(strpos($memFace, 'https') !== false ? $memFace : str_replace("http", "https", $memFace)),
                'mem_role'        => $role,
                'mem_mobile'      => $ucResult['memMobile'],
                'department_name' => $ucResult['dpName'][0]['dpName'],
                'is_manager'      => $role == '店长' ? 1 : 2,// 是否是管理者（1：是，2：不是）
                'is_store'        => $isStore,//  是否是门店类型的组织员工（1：是，2：否）
                'shop_app_id'     => $shopData['shop_app_id'],//  商城小程序APPID
                'avatar_new'      => stripslashes(strpos($memFace, 'https') !== false ? $memFace : str_replace("http", "https", $memFace)),
                'qr_code_new'     => stripslashes(strpos($data['qrCode'], 'https') !== false ? $data['qrCode'] : str_replace("http", "https", $data['qrCode'])),
            ]
        ]);
    }

//    /**
//     * 图片转base64
//     *
//     * @param string $imageUrl
//     *
//     * @return string
//     */
//    public function imageToBase64($imageUrl = '')
//    {
//        Log::info('url', ['usl' => $imageUrl]);
//
//        if (!$imageUrl) {
//            return '';
//        }
//
//        try {
//            $type         = getimagesize($imageUrl);
//            $content      = file_get_contents($imageUrl);
//            $file_content = base64_encode($content);
//
//            switch ($type[2]) {//判读图片类型
//                case 1:
//                    $img_type = "gif";
//                    break;
//                case 2:
//                    $img_type = "jpg";
//                    break;
//                case 3:
//                    $img_type = "png";
//                    break;
//                default:
//                    return '';
//            }
//
//            //合成图片的base64编码
//            $img = 'data:image/' . $img_type . ';base64,' . $file_content;
//
//            return $img;
//        } catch (\Exception $e) {
//            Log::info('to-base64-error', ['msg' => $e->getMessage(), 'line' => $e->getLine(), 'file' => $e->getFile()]);
//            return '';
//        }
//    }
}
