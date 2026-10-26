<?php
/**
 * 获取任务列表
 */

namespace App\Http\Controllers\App\V1\MarketingPlan;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\GoMicroScrm\Plan;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class TaskListController extends Controller
{
    /**
     * 获取计划列表
     * User: HuiYang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/13 17:44
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;

        try {
            $result = (new Plan())->planList($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    public function taskList(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;
        if (!isset($input['status'])) {
            $input['status'] = -1;
        }
        try {
            $result = (new Plan())->planTaskList($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    public function taskInfo(Request $request)
    {
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;
        $input['member_id'] = $this->memberAuthInfo['member_id'];
        try {
            $result = (new Plan())->planTaskInfo($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    public function tagUserList(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;

        try {
            $result = (new Plan())->tagUserList($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    public function setStatus(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;

        try {
            $result = (new Plan())->setTaskStatus($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    public function updateStatus(Request $request)
    {
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile']        = $this->mobile;

        try {
            $result = (new Plan())->updateTaskStatus($input);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    public function getVideoMedia(Request $request)
    {
        $input = $request->input();

        $result = [];
        try {
            $videoDetail = (new Plan())->getVideoDetail($input);
            $param       = [
                "enterprise_id" => isset($videoDetail['enterprise_id']) ? $videoDetail['enterprise_id'] : '',
                "source_id"     => isset($videoDetail['id']) ? $videoDetail['id'] : '',
                "source_type"   => "video",
                "source_url"    => isset($videoDetail['url']) ? $videoDetail['url'] : '',
                "source_name"   => isset($videoDetail['name']) ? $videoDetail['name'] : ''
            ];

            if ($param['enterprise_id'] && $param['source_id'] && $param['source_url']) {
                $result = (new Plan())->getMakeMedia($param);
            }

        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }
}
