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


class SpuController extends BaseController
{

    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 商品(spu)列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function list(Request $request)
    {
        $params = $request->only(['page', 'size', 'keyword', 'is_sale']);

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getSpuList($params);

        return success($result);
    }

    /**
     * 商品(spu)详情
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $params = $request->only(['id']);

        $rules = [
            'id' => 'required',
        ];

        $message = [
            'id.required' => '商品id不能为空'
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->getSpuDetail($params);

        return success($result);
    }

    /**
     * 商品保存
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function save(Request $request)
    {
        $params = $request->only(['id', 'name', 'image_url', 'sale_price', 'video_list', 'image_list']);

        $rules = [
            'id'         => 'required',
            'name'       => 'required',
            'image_url'  => 'required',
            'sale_price' => 'required|integer|min:1',
            'video_list' => 'required|array',
            'image_list' => 'required|array',
        ];

        $message = [
            'id.required'         => '商品id不能为空',
            'name.required'       => '商品名称不能为空',
            'image_url.required'  => '商品缩略图不能为空',
            'sale_price.required' => '商品售价不能为空',
            'sale_price.integer'  => '商品售价类型错误',
            'sale_price.min'      => '商品售价不能为零',
            'image_list.required' => '图片内容不能为空',
            'video_list.required' => '视频内容不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new DoveDream)->saveSpu($params);

        return success($result);
    }
}