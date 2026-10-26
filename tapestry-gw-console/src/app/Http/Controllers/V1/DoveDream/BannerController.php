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

class BannerController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 获取轮播图
     * @return array
     * @throws \Exception
     */
    public function detail()
    {
        $conditions = [
            'enterprise_id' => $this->enterpriseId,
        ];

        $result = (new DoveDream)->getBanners($conditions);

        return success($result);
    }

    /**
     * 保存轮播图
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function save(Request $request)
    {
        $params = $request->only(['content']);

        $rules = [
            'content' => 'required'
        ];

        $message = [
            'content.required' => '轮播图不能为空'
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->saveBanner($params);

        return success($result);
    }
}