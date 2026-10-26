<?php
/**
 * Created by sublime.
 * User: zhao qi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\Impress;

use CrsSdk\Shop\ImpressLabel;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\Label\LabelBaseController;

class StatusController extends LabelBaseController
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
        $params = $request->only(['label_ids', 'status']);

        $rules = [
            'label_ids' => 'required|array',
            'status' => 'required|integer'
        ];

        $message = [
            'label_ids.required' => '标签id集合不能为空',
            'label_ids.array' => '标签id集合为数组',
            'status.required' => '状态不能为空',
            'status.integer' => '状态为整型',
        ];

        $this->validate(request(), $rules, $message);

        if (!in_array($params['status'], [self::LABEL_STATUS_OPEN, self::LABEL_STATUS_CLOSE])) {
            error(800011);
        }

        $params = request()->input();
        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new ImpressLabel())->status($params);

        if (!$result) {
            notice(100000, [], $params);
        }

        // 日志记录
        $this->writeOperatorLog(($params['status'] == self::LABEL_STATUS_OPEN) ? '启用' : '禁用' . '了 ' . count($params['label_ids']) . ' 个印象标签');

        return $result;
    }
}
