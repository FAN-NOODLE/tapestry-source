<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/19 22:24
 */

namespace App\Http\Controllers\V1\DistributionConfig;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionConfig;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class PartnerDetailController extends BaseController
{
    /**
     * 分销配置存储
     *
     * @param Request $request
     * @param WorkerDistribution $workerDistribution
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request,DistributionConfig $distribution)
    {
        $input = $request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        try {
            $detail = $distribution->partnerDetail($input);
            return success($detail);
        } catch (Error $error) {
            Log::error("Work:Distribution->PartnerDetail，错误信息：{$error->getMessage()},错误码：{$error->getCode()}", $input);
            notice(1721006);
        }
    }
}