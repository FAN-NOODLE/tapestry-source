<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearchSw;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class CrowdSwController extends BaseController
{

    public function addCrowd(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params = array_merge($params, $this->getRoleRange());
        return success((new ExternalSearchSw())->addCrowd($params));
    }


    public function crowdList(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdList($params));
    }

    public function crowdInfo(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdInfo($params));
    }

    public function crowdOpenIds(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdOpenIds($params));
    }
    public function crowdEdit(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdEdit($params));
    }
    public function crowdDel(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdDel($params));
    }

    public function crowdUpdateSort(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdUpdateSort($params));
    }

    public function crowdUpdateContact(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdUpdateContact($params));
    }

    public function crowdBatchUpdateContact(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearchSw())->crowdBatchUpdateContact($params));
    }

}