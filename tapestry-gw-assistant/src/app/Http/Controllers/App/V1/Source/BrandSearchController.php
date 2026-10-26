<?php

namespace App\Http\Controllers\App\V1\Source;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class BrandSearchController extends BaseController
{


    /**
     * 小程序端各素材列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2024/5/7
     */
    public function list(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new GoMicroScrm())->brandSearchList($params));

    }

    public function listBySku(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new GoMicroScrm())->brandSearchListGroupSku($params));

    }

    public function talkieImageSearch(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new GoMicroScrm())->talkieImageSearch($params));

    }

    /**
     * 小程序端各素材列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2024/5/7
     */
    public function getCacheInfo(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        return success((new GoMicroScrm())->getCacheInfo($params));

    }



    /**
     * 小程序端各素材列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2024/5/7
     */
    public function collectSearchList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new GoMicroScrm())->collectSearchList($params));

    }
}
