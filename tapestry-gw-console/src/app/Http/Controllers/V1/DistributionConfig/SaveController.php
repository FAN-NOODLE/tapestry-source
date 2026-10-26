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

class SaveController extends BaseController
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
        $rules = [
            'distribution_text'   => 'required|max:200',
            'distribution_title'  => 'required|max:20',
            'settlement_standard' => 'required',
            'poster_img'          => 'required|url',
            'id'                  => 'integer'
        ];

        $messages               = [
            'distribution_text.*'   => '分销描述不合法',
            'distribution_title.*'  => '分销标题不合法',
            'settlement_standard.*' => '分销结算标准不合法',
            'poster_img.*'          => '海报背景图不合法',
            'id.*'                  => '分销ID有误'
        ];
        $input                  = $this->validate($request, $rules, $messages);
        $input['enterprise_id'] = $this->enterpriseId;
        try {
            $distribution->save($input);
        } catch (Error $error) {
            Log::error("Work:Distribution->configSave方法出错，错误信息:{$error->getMessage()},错误Code:{$error->getCode()}", $input);
            notice($error->getCode(), $error->getMessage());
        }
        return success();
    }
}
