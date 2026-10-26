<?php
/**
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

class AddCategoryTemplateController extends BaseController
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
            'name'  => 'required',
            'level' => 'required|integer|in:1,2,3',
            'type'  => 'required|in:1,2',
        ];

        $messages = [
            'name.*'  => '模版名称有误',
            'level.*' => '模版等级有误',
            'type.*'  => '模版类型有误',
        ];

        $this->validate($request, $rule, $messages);

        // 组装数据
        $input                  = $request->only([
            'name', 'level', 'type'
        ]);
        $input['enterprise_id'] = $this->enterpriseId;

        $CategoryTemplate = new CategoryTemplate();
        try {
            $result = $CategoryTemplate->categoryTemplateAdd($input);

            if (isset($result['id']) && !empty($result['id'])) {
                // init 模版下级分类
                $result['cate_init'] = $CategoryTemplate->categoryAdd([
                    'enterprise_id' => $input['enterprise_id'],
                    'cate_tpl_id'   => $result['id']
                ]);
            }

            // 日志记录
            $this->writeOperatorLog('添加【' . $input['name'] . '】模版分类');
        } catch (Error $error) {
            Log::info('修改分类关联商品权重失败', $error->data);
            notice(100000, [$error->getMessage()]);
        } catch (\Exception $e) {
            Log::info('修改分类关联商品权重失败' . $e->getMessage(), $input);
            notice(100000, [$e->getMessage()]);
        }

        return success($result);
    }
}