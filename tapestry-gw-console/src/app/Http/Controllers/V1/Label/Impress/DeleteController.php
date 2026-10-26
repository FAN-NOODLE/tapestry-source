<?php
/**
 * Created by sublime.
 * User: zhaoqi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\Impress;

use CrsSdk\Shop\ImpressLabel;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\Label\LabelBaseController;

class DeleteController extends LabelBaseController
{
    /**
     * 印象标签删除
     *
     * @param Request $request
     *
     * @return array
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['label_ids']);

        $rules = [
            'label_ids' => 'required|array'
        ];

        $message = [
            'label_ids.required' => '标签id集合不能为空',
            'label_ids.array' => '标签id集合为数组',
        ];

        $this->validate(request(), $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;

        $result =  (new ImpressLabel())->delete($params);

        if (!$result) {
            notice(100000, [], $params);
        }

        // 日志记录
        $this->writeOperatorLog('删除了 ' . count($params['label_ids']) . ' 个印象标签');

        return $result;
    }
}
