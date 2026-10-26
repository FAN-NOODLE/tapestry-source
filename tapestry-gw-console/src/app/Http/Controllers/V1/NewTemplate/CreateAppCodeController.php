<?php
/**
 * 生成物料二维码
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2020/06/17
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\NewTemplate;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\NewTemplate;

class CreateAppCodeController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->input();

        $input['enterprise_id'] = $this->enterpriseId;
        $input['source_id']     = $this->storeIds;


        $this->writeOperatorLog("新增首页模板");

        $data = (new NewTemplate())->CreateTemplateCode($input);

        return success($data);
    }
}
