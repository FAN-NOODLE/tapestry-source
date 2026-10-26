<?php

namespace App\Http\Controllers\App\V1\Excitation;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExcitationMedal;
use App\Micros\MicroScrm\ExcitationScrm;
use Illuminate\Http\JsonResponse;

/**
 * 奖章
 */
class MedalController extends BaseController
{
    /**
     * 奖章墙
     * @return JsonResponse
     */
    public function MedalWall():JsonResponse
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationMedal::MedalWall($params);
        return success($result);
    }

    /**
     * 奖章墙详情
     * @return JsonResponse
     */
    public function MedalWallDetail():JsonResponse
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationMedal::MedalWallDetail($params);
        return success($result);
    }

    /**
     * 领取奖章数量和弹框
     * @return JsonResponse
     */
    public function RecordAlreadyGot():JsonResponse
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationMedal::RecordAlreadyGot($params);
        return success($result);
    }

    /**
     * 领取奖章
     * @return JsonResponse
     */
    public function RecordMemberReceived():JsonResponse
    {
        $params = $this->request->all();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        $result = ExcitationMedal::RecordMemberReceived($params);
        return success();
    }

}