<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearch;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class CrowdController extends BaseController
{

    public function addCrowd(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params = array_merge($params, $this->getRoleRange());
        try{
            $res=(new ExternalSearch())->addCrowd($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }

    }


    public function crowdList(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdList($params);
            return successJson($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }

    public function crowdInfo(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdInfo($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }

    public function crowdOpenIds(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdOpenIds($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }
    public function crowdEdit(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdEdit($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }
    public function crowdDel(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdDel($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }

    public function crowdUpdateSort(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdUpdateSort($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }

    public function crowdUpdateContact(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        try{
            $res=(new ExternalSearch())->crowdUpdateContact($params);
            return success($res);
        }catch (\Exception $e){
            error(10002,$e->getMessage());
        }
    }


}