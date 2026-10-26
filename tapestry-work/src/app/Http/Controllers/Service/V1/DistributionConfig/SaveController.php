<?php
/**
 * 全员分销-基础配置-设置
 * Created by PhpStorm.
 * User: weihua
 * Date: 2020-09-24
 * Time: 18:09
 */

namespace App\Http\Controllers\Service\V1\DistributionConfig;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Distribution;
use Illuminate\Http\Request;

class SaveController extends BaseController
{
    public function index(Request $request)
    {
        $input            = $this->getParams($request);
        $distributionData = [
            'enterprise_id'       => $input['enterprise_id'],
            'distribution_title'  => $input['distribution_title'],
            'distribution_text'   => $input['distribution_text'],
            'settlement_standard' => $input['settlement_standard'],
            'poster_img'          => $input['poster_img'],
        ];
        try {
            $id = $input['id'];
            // 更新分销配置表
            $distribution = Distribution::query()
                ->where('id', $id)
                ->firstOrCreate(['id' => $id]);
            $distribution->update($distributionData);
        } catch (Error $error) {
            \error($error->getCode(), $error->getMessage());
        } catch (\Exception $exception) {
            notice('worker_distribution_save_error');
        }
        return success();
    }

    /**
     * 获取请求参数
     *
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams(Request $request): array
    {
        $rules   = [
            'enterprise_id'       => 'required|integer',
            'distribution_title'  => 'required|max:20',
            'distribution_text'   => 'required|max:200',
            'settlement_standard' => 'required',
            'poster_img'          => 'required|url',
            'id'                  => 'integer'
        ];
        $message = [
            'enterprise_id.*'       => 'enterprise_id有误',
            'distribution_title.*'  => 'distribution_title有误',
            'distribution_text.*'   => 'distribution_text有误',
            'settlement_standard.*' => 'settlement_standard有误',
            'poster_img.*'          => 'poster_img有误',
        ];
        $input       = $this->validate($request, $rules,$message);
        $input['id'] = $input['id'] ?? 0;
        return $input;
    }

}