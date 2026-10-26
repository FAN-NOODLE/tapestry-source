<?php

namespace App\Http\Controllers\V1\Label\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\Label\LabelBaseController;
use CrsSdk\MicroGoodsV1\Label;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class DeleteController extends LabelBaseController
{
    /**
     * 删除标签
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['label_id']);

        $rules = [
            'label_id' => 'required|integer',
        ];

        $message = [
            'label_id.*'  => '标签不存在',
        ];

        $this->validate(request(), $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;

        try {
            (new Label())->delete($params);
            // 日志记录
            $this->writeOperatorLog('删除 ' . $params['label_id'] . ' 标签');
        } catch (Error $error) {
            notice(450000, [$error->getMessage()], $error->data);
        } catch (\Exception $exception) {
            Log::info('删除标签失败' . $exception->getMessage(), $params);
            notice(450000, [$exception->getMessage()]);
        }

        return success();
    }

}
