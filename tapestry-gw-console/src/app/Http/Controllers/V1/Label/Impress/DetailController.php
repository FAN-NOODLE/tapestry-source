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

class DetailController extends LabelBaseController
{
    /**
     * 印象标签详情
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
        return (new ImpressLabel())->detail($params);
    }
}
