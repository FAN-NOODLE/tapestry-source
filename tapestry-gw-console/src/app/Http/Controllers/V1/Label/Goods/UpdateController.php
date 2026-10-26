<?php
/**
 * Created by PhpStorm.
 * Function: xxx
 * User: 龚国新
 * Date: 2020/3/5
 * Time: 15:37
 */

namespace App\Http\Controllers\V1\Label\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\Label\LabelBaseController;
use CrsSdk\MicroGoodsV1\Label;
use Illuminate\Support\Facades\Log;

class UpdateController extends LabelBaseController
{
    /**
     * 修改标签
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        //请求数据校验
        $params                  = $this->validates();
        $params['enterprise_id'] = $this->enterpriseId;

        try {
            (new Label())->update($params);
            // 日志记录
            $this->writeOperatorLog('修改【' . $params['label_id'] . '】标签');
        } catch (Error $error) {
            notice(450000, [$error->getMessage()], $error->data);
        } catch (\Exception $exception) {
            Log::info('修改标签失败' . $exception->getMessage(), $params);
            notice(450000, [$exception->getMessage()]);
        }

        //返回结果
        return success();
    }

    /**
     * 请求数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'label_id' => 'required|integer',
            'name'     => 'required|between:0,50'
        ];

        $message = [
            'label_id.required' => '标签ID不能为空',
            'label_id.integer'  => '标签ID数据类型有误',
            'name.between'      => '名称介于0-50个字符',
            'name.required'     => '名称不能为空',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }
}