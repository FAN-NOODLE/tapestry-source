<?php

namespace App\Http\Controllers\App\V1\Invite;

use App\Http\Controllers\Controller;
use App\Micros\Common;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Micros\Work\Work;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

/**
 * 用户端 邀请函
 */
class IndexController extends Controller
{

    /**
     * @var Sdk
     */
    protected $microScrm;

    public function __construct()
    {
        parent::__construct();
        $this->microScrm = Common::microScrm();
    }

    /**
     * 查看邀请函信息
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/9
     */
    public function getRecordInfo(Request $request): JsonResponse
    {
        // 参数校验
        $rules   = [
            'invite_record_id' => 'required|integer',
        ];
        $message = [
            'invite_record_id.required' => 'invite_record_id参数不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/invite/get-record-info')->post($params);
        } catch (Exception $e) {
            if ($e->getMessage() == "无权访问") {
                error(10001, "权限不足");
            }
            error(10002, $e->getMessage());
        }

        return success($result);

    }

    /**
     * 发送邀约单记录
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/11/9
     */
    public function createOrder(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/invite/create-order')->post($params);
        } catch (Exception $e) {
            if ($e->getMessage() == "无权访问") {
                error(10001, "权限不足");
            }

            error(10002, $e->getMessage());
        }

        return success($result);

    }


    /**
     * 我的页面
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/10/21
     */
    public function myCenter(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $userDetail = (new GoMicroScrm())->myCenter($params);

        if ($userDetail['union_id']) {

            $userDetail['customer_type'] = '';
            $userDetail['name']          = '';
//            $userDetail['mobile']                   = '';
            $userDetail['client_type'] = '';

            // 查询crm
            $userReq['enterprise_id'] = $params['enterprise_id'];
            $userReq['union_id']      = $userDetail['union_id'];

            $userCrm = (new Work())->crmMemberInfo($userReq);
            if ($userCrm['status']) {
                $userDetail['customer_type'] = $userCrm['customer_type'];
                $userDetail['name']          = $userCrm['name'];
//                $userDetail['mobile']                   = $userCrm['mobile'];
                $userDetail['client_type'] = $userCrm['client_type']??"";
            }

        }

        return success($userDetail);

    }

    /**
     * 我的页面 - 添加专属导购
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/11/16
     */
    public function mySa(Request $request): JsonResponse
    {
        // 参数校验
//        $rules   = [
//            'work_external_user_id'      => 'required|integer',
//        ];
//        $message = [
//            'work_external_user_id.integer' => 'work_external_user_id参数类型有误',
//            'work_external_user_id.required'     => 'work_external_user_id参数不能为空',
//        ];
//        $this->validate($request, $rules, $message);
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $userSA = (new GoMicroScrm())->mySa($params);

        return success($userSA);

    }

    /**
     * 邀请函模板信息
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @author: wangle@vchangyi.com
     * @Time: 2022/12/7
     */
    public function inviteInfo(Request $request): JsonResponse
    {
        // 参数校验
        $rules                   = [
            'id' => 'required|integer',
        ];
        $message                 = [
            'id.required' => '邀请函id参数不能为空',
        ];
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $this->validate($request, $rules, $message);
        $info = (new GoMicroScrm())->inviteInfo($params);

        $res                         = [];
        $res['id']                   = $info['id'];
        $res['title_name']           = $info['title_name'];
        $res['content']              = $info['content'];
        $res['reserve_service_id']   = $info['reserve_service_id'];
        $res['reserve_service_name'] = $info['reserve_service_name'];
        $res['reserve_service_time'] = $info['reserve_service_time'];
        $res['style_type']           = 1;
        $res['background_url']       = json_decode($info['background_url'])[0];
        $res['send_content']         = $info['content'];
        $res['user_is_show']         = $info['user_is_show'];
        $res['member_is_show']       = $info['member_is_show'];

        return success($res);

    }

    /**
     * 用户确认邀请函的时候SW&&KS需要验证到店时间的限制
     * @param Request $request
     * @return JsonResponse
     */
    public function checkReserveCreate(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $info = (new GoMicroScrm())->checkReserveCreate($params);
        return success($info);
    }
}