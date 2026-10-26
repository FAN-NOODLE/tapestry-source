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

class AddController extends LabelBaseController
{
    /**
     * 印象标签添加
     *
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['name', 'describe']);

        // 验证名称是否为空，防止填写空格提交
        if (empty(preg_replace('# #','',$params['name']))) {
            error(800001);
        }

        if (empty(preg_replace('# #','',$params['describe']))) {
            error(800012);
        }

        $rules = [
            'name' => 'between:0,7',
            'describe' => 'between:0,101',
        ];

        $message = [
            'name.between' => '名称介于0-6个字符',
            'describe.between' => '描述介于0-100个字符',
        ];

        $this->validate(request(), $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new ImpressLabel())->add($params);

        if (!$result) {
            notice(100000, [], $params);
        }

        // 日志记录
        $this->writeOperatorLog('添加 ' . $params['name'] . ' 印象标签');

        return $result;
    }
}
