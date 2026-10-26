<?php

namespace App\Http\Controllers\App\V1\Source;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class UploadController extends BaseController
{

    /**
     * 验证用户身份是否有上传素材权限
     * @param Request $request
     * @return JsonResponse
     */
    public function checkAuth(Request $request): JsonResponse
    {
        $mediaId                 = 393;
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['userid']=$this->memberAuthInfo['userid'];
        $auth                    = 0;
        try {
            $result      = Common::Authority()->setApi("micro-authority/v1/account/account-permissions")->post($params);
            $permissions = $result['permissions'];
            if (in_array($mediaId, $permissions)) {
                $auth = 1;
            }
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        $res = ['auth' => $auth];
        return success($res);
    }

    /**
     * 小程序端上传素材
     * @param Request $request
     * @return JsonResponse
     */
    public function uploadMedia(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $params['userid']      =  $this->memberAuthInfo['userid'];
        try {
            $result = Common::MicroScrm()->setApi('micro-scrm/v1/source/upload-media')->post($params);
            return success($result);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
    }

    /**
     * 获取小程序端的上传签名
     * @param Request $request
     * @return JsonResponse
     */
    public function credential(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = Common::MicroScrm()->setApi('micro-scrm/v1/source/credential')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 获取压缩结果
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function compress(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try {
            $result = Common::MicroScrm()->setApi('micro-scrm/v1/source/compress')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * KOL素材，图片/视频上传文件自动压缩为朋友圈要求尺寸
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function compressFile(Request $request): JsonResponse
    {
        $params                  = $request->all();
        try {
            $result = Common::MicroScrm()->setApi('micro-scrm/v1/source/compress-file')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success(["file_url"=>$result]);
    }

    /**
     * KOL素材，图片/视频上传文件自动压缩为朋友圈要求尺寸
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function editVideoCreateTime(Request $request): JsonResponse
    {
        $params                  = $request->all();
        try {
            $result = Common::MicroScrm()->setApi('micro-scrm/v1/source/edit-video-create-time')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success(["file_url"=>$result]);
    }

    /**
     * 验证图片资源地址是否满足企业朋友圈的素材要求
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function getMediaInfo(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = Common::MicroScrm()->setApi('micro-scrm/v1/source/get-media-info')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success(["media_id"=>$result]);
    }

    /**
     * ios端视频生成封面图
     * @param Request $request
     * @return JsonResponse
     */
    public function makeCover(Request $request): JsonResponse
    {
        $params                  = $request->all();
        try {
            $result= Common::MicroScrm()->setApi('micro-scrm/v1/source/make-cover')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }
}