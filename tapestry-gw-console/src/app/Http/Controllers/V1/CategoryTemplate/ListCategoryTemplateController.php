<?php
/**
 * 模版分类列表
 * Created by PhpStorm.
 * User: zy
 * Date: 2020/3/4
 * Time: 15:31
 */

namespace App\Http\Controllers\V1\CategoryTemplate;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CategoryTemplate;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class ListCategoryTemplateController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {

        // 验证参数
        $rule = [
            'type' => 'required|in:1,2',
        ];

        $messages = [
            'type.*' => '模版类型有误',
        ];

        $input                  = $this->validate($request, $rule, $messages);
        $input['size']          = $request->input('size', 10);
        $input['page']          = $request->input('page', 1);
        $input['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new CategoryTemplate())->categoryTemplateList($input);
        } catch (Error $error) {
            Log::info('获取分类模版列表失败', $error->data);
            notice(100000, [$error->getMessage()], $error->data);
        } catch (\Exception $e) {
            Log::info('获取分类模版列表失败' . $e->getMessage(), $input);
            notice(100000, [$e->getMessage()]);
        }

        return success($list);

    }
}