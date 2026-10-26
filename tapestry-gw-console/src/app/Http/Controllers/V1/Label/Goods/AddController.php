<?php

namespace App\Http\Controllers\V1\Label\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\Label\LabelBaseController;
use CrsSdk\MicroGoodsV1\Label;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AddController extends LabelBaseController
{
    /**
     * 添加标签
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['name']);

        $rules = [
            'name' => 'required|between:0,50',
        ];

        $message = [
            'name.between'  => '名称介于0-50个字符',
            'name.required' => '名称不能为空',
        ];

        $this->validate(request(), $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            $result = (new Label())->add($params);
            // 日志记录
            $this->writeOperatorLog('添加 ' . $params['name'] . ' 标签');
        } catch (Error $error) {
            notice(450000, [$error->getMessage()], $error->data);
        } catch (\Exception $exception) {
            Log::info('添加标签失败' . $exception->getMessage(), $params);
            notice(450000, [$exception->getMessage()]);
        }

        return success($result);
    }

}
