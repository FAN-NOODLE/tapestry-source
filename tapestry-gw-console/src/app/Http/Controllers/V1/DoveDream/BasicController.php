<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2020/2/20                   
* +------------------------------------+
* | Time   : 16:39                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\V1\DoveDream;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Validator;
use Illuminate\Http\Request;
use CrsSdk\Shop\DoveDream;


class BasicController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 基本信息详情
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $params = $request->only(['type']);

        $rules = [
            'type' => 'required|in:1,2'
        ];

        $message = [
            'type.required' => '类型不能为空',
            'type.in'       => '类型范围错误'
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getBasic($params);

        return success($result);
    }

    /**
     * 基本信息保存
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function save(Request $request)
    {
        $params = $request->only(['type', 'content']);

        $rules = [
            'type'    => 'required|in:1,2',
            'content' => 'required',
        ];

        $message = [
            'type.required'    => '简介类型不能为空',
            'type.in'          => '简介类型不在范围',
            'content.required' => '文本内容不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->saveBasic($params);

        return success($result);
    }
}