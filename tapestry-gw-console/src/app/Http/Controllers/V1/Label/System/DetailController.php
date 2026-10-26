<?php
/**
 * Created by sublime.
 * User: zhaoqi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\System;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\SystemLabel;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 系统标签详情
     *
     * @param Request $request
     *
     * @return array
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['label_id']);

        $rules = [
            'label_id' => 'required|integer|min:1',
        ];

        $message = [
            'label_id.required' => '标签ID错误',
            'label_id.integer' => '标签ID错误',
            'label_id.min' => '标签ID错误'
        ];

        $this->validate(request(), $rules, $message);

        $labelId = $request->get('label_id');
        $params['label_id'] = $labelId;
        $params['enterprise_id'] = $this->enterpriseId;
        return (new SystemLabel())->detail($params);

    }
}
