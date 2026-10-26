<?php
/**
 * Created by sublime.
 * User: zhaoqi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\System;

use CrsSdk\Shop\SystemLabel;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\Label\LabelBaseController;

class GoodController extends LabelBaseController
{
    /**
     * 标签商品列表
     *
     * @param Request $request
     *
     * @return array
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['label_id', 'page', 'size']);

        $rules = [
            'label_id' => 'required|integer',
            'page' => 'required|integer|min:1',
            'size' => 'required|integer|min:1',
        ];

        $message = [
            'label_id.required' => '标签ID错误',
            'label_id.integer' => '标签ID错误',
            'page.required' => '无效的页码',
            'page.integer' => '无效的页码',
            'page.min' => '无效的页码',
            'size.required' => '无效的条数',
            'size.integer' => '无效的条数',
            'size.min' => '无效的条数'
        ];

        $this->validate(request(), $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;

        return (new SystemLabel())->good($params);
    }
}
