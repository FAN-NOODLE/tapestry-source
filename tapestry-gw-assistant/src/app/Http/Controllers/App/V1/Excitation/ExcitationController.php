<?php

namespace App\Http\Controllers\App\V1\Excitation;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExcitationScrm;
use Illuminate\Http\JsonResponse;

class ExcitationController extends BaseController
{

    /**
     * We Journey 首页
     *
     * @return JsonResponse
     */
    public function index():JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationScrm::index($params);
        return success($result);
    }

    /**
     * 等级信息
     * @return JsonResponse
     */
    public function gradePoints():JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationScrm::gradePoints($params);
        return success($result);
    }

    /**
     * 积分列表
     * @return JsonResponse
     */
    public function pointsList():JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id']     = $this->memberAuthInfo['member_id'];

        $result = ExcitationScrm::pointsList($params);
        return success($result);
    }

    /**
     * 积分规则说明
     * @return JsonResponse
     */
    public function integralRuleDescription():JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationScrm::integralRuleDescription($params);
        return success($result);
    }

    /**
     * 积分回去信息
     * @return JsonResponse
     */
    public function getPointsRules():JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id']     = $this->memberAuthInfo['member_id'];

        $result = ExcitationScrm::getPointsRules($params);
        return success($result);
    }
    /**
     * 等级变动记录
     * @return JsonResponse
     */
    public function levelRecord():JsonResponse
    {
        $params = request()->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id']     = $this->memberAuthInfo['member_id'];

        $result = ExcitationScrm::levelRecord($params);
        return success($result);
    }

    /**
     * 金币列表
     * @return JsonResponse
     */
    public function coinsList():JsonResponse
    {
        $params = request()->all();


        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id']     = $this->memberAuthInfo['member_id'];

        $result = ExcitationScrm::coinsList($params);
        return success($result);
    }

    /**
     * 金币配置获取
     * @return JsonResponse
     */
    public function coinsSettingsInfo():JsonResponse
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        //获取商城配置
        $result = ExcitationScrm::coinsSettingsInfo($params);
        $loginToken = '';
        if (isset($params['login_str'])) {
            //获取welearning商城跳转参数token
            $loginToken = hash('sha256', env("WE_LEARNING_TOKEN_SALT", 'sinostrong') . $params['login_str']);
        }
        return success(['coins_settings_info' => $result,'login_token' => strtoupper($loginToken)]);
    }


}