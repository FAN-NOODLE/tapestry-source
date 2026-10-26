<?php

namespace App\Http\Controllers\App\V1\Source;

use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\MicroScrm\GoMicroScrm;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class ListController extends Controller
{


    /**
     * 素材产品类别列表
     * @param Request $request
     * @return JsonResponse
     */
    public function categoryList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $list                    = (new GoMicroScrm())->getSourceCategoryLists($params);
        return success($list);

    }

    /**
     *  素材分组归属门店列表
     * @param Request $request
     * @return JsonResponse
     */
    public function groupList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $list                    = (new GoMicroScrm())->getSourceRelationList($params);
        return success($list);
    }

    /**
     * 根据门店获取素材分组列表
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function getListForStore(Request $request): JsonResponse
    {
        $params                   = $request->all();
        $params['enterprise_id']  = $this->enterpriseId;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $params['member_user_id'] = $this->memberAuthInfo['userid'];
        $list                     = (new GoMicroScrm())->getListForStore($params);
        return success($list);
    }

    /**
     * 小程序端各素材列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/1
     */
    public function index(Request $request): JsonResponse
    {
        $param=$request->all();
        $rules   = [
            'page' => 'required|integer',
            'size' => 'required|integer',
            'type' => 'required|integer',
        ];
        $message = [
            'type.required' => 'type不能为空',
            'type.integer'  => 'type类型有误',
        ];
        $params  = $this->validate($request, $rules, $message);

        $params['enterprise_id']  = $this->enterpriseId;
        $params['mobile']         = $this->mobile;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $params['member_user_id'] = $this->memberAuthInfo['userid'];
        $param=array_merge($param,$params);

        $list = (new GoMicroScrm())->selectList($param);

        return success($list);

    }

    /**
     * 查询各素材上新数量
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/11/1
     */
    public function getNewNum(Request $request): JsonResponse
    {
        $params=$request->all();
        $params['enterprise_id']  = $this->enterpriseId;
        $params['mobile']         = $this->loginMobile();
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $params['member_user_id'] = $this->memberAuthInfo['userid'];
        $list                     = (new GoMicroScrm())->getNewSourceNum($params);

        return success($list);

    }


    /**
     * 上报计次
     * @return void
     */
    public function reportTimes(): JsonResponse
    {
        $request = $this->request->all();

        $request['enterprise_id'] = $this->enterpriseId;
        $result = (new GoMicroScrm())->reportTimes($request);
        return success($result);
    }


    /**
     * 验证是否用户首次打开素材交互提示
     *
     * @return JsonResponse
     */
    public function checkFirst(): JsonResponse
    {
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $params['userid'] = $this->memberAuthInfo['userid'];
        try{
            $result=(new GoMicroScrm())->checkFirst($params);
        }catch (Exception $e){
            error(10002, $e->getMessage());
        }

        return success(["first"=>$result]);
    }


    /**
     * 小程序端收藏各素材列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/1
     */
    public function collectSelectList(Request $request): JsonResponse
    {
        $param=$request->all();
        $rules   = [
            'page' => 'required|integer',
            'size' => 'required|integer',
            'type' => 'required|integer',
        ];
        $message = [
            'type.required' => 'type不能为空',
            'type.integer'  => 'type类型有误',
        ];
        $params  = $this->validate($request, $rules, $message);

        $params['enterprise_id']  = $this->enterpriseId;
        $params['mobile']         = $this->mobile;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $params['member_user_id'] = $this->memberAuthInfo['userid'];
        $param=array_merge($param,$params);

        $list = (new GoMicroScrm())->collectSelectList($param);

        return success($list);

    }

    /**
     * 收藏素材
     * @param Request $request
     * @return JsonResponse
     */
    public function collect(Request $request): JsonResponse
    {
        $param=$request->all();
        $params['enterprise_id']  = $this->enterpriseId;
        $params['mobile']         = $this->mobile;
        $params['member_id']      = $this->memberAuthInfo['member_id'];
        $param=array_merge($param,$params);

        $res = (new GoMicroScrm())->collect($param);
        return success($res);
    }
}
