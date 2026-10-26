<?php
/**
 * @Description
 * @Author:         tuzhengping
 * @Date:           2020-03-19 17:19
 * @Version:        1.0
 */

namespace App\Http\Controllers\V1\CategoryTemplate;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CategoryTemplate;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class DetailCategoryTemplateController extends BaseController
{
    /**
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function index(Request $request)
    {
        // 验证参数
        $rule = [
            'id' => [
                'required',
                'integer'
            ]
        ];

        $msg = [
            'id.*' => '模版不存在',
        ];

        $input                  = $this->validate($request, $rule, $msg);
        $input['enterprise_id'] = $this->enterpriseId;

        try {
            $info = (new CategoryTemplate())->categoryTemplateDetail($input);
        } catch (Error $error) {
            Log::info('获取分类模版失败', $error->data);
            notice(100000, [$error->getMessage()]);
        } catch (Exception $e) {
            Log::info('获取分类模版失败' . $e->getMessage(), $input);
            notice(100000, [$e->getMessage()]);
        }

        return success($info);

    }


}