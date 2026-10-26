<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-02
 * Time: 16:15
 */

namespace App\Http\Controllers\V1\DistributionConfig;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionConfig;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class DetailController extends BaseController
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
            $detail = $distribution->detail($input);
            return success($detail);
        } catch (Error $error) {
            Log::error("Work:Distribution->detail，错误信息：{$error->getMessage()},错误码：{$error->getCode()}", $input);
            notice(1721006);
        }
    }
}
