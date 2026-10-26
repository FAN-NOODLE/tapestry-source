<?php
/**
 * Description:新增用户自定义标签及绑定关系
 * User: lintong
 * Email: 18511245965@163.com
 * Date: 2020/8/2
 * Time: 15:45
 */

namespace App\Http\Controllers\Service\V1\CustomerTags;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\UserCustomizeTagsService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class AddController extends BaseController
{
    /**
     * 检验参数
     * @param $post
     * @throws \Illuminate\Validation\ValidationException
     */
    protected function checkParams($post)
    {
        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'worker_id'     => 'required|integer',
        ];

        $message   = [
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'enterprise_id.min'      => '企业ID数据类型有误',
            'worker_id.required'     => '员工ID不能为空',
            'worker_id.integer'      => '员工ID数据类型有误',
        ];
        $validator = Validator::make($post, $rules, $message);
        if ($validator->fails()) {
            throw new \Illuminate\Validation\ValidationException($validator);
        }
    }

    /**
     * 添加自定义标签
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function addNewTag(Request $request)
    {
        $post = $request->only([
            'enterprise_id',
            'name',
            'worker_id',
        ]);
        $this->checkParams($post);
        $post['name'] = $post['name'] ?? '';
        $result       = (new UserCustomizeTagsService())->addNewTags($post);
        $result       = is_array($result) ? $result : ['id' => $result];
        return success(['code' => 200, 'message' => '添加成功', 'data' => $result]);
    }

    /**
     * 绑定用户自定义标签
     * @param Request $request
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    public function addNewRelation(Request $request)
    {
        $post = $request->only([
            'enterprise_id',
            'relation_id',
            'worker_id',
            'tag_names'
        ]);
        $this->checkParams($post);
        // 需要绑定的标签id
        $tagsName = $post['tag_names'] ?? [];
        if (!is_array($tagsName) || empty($tagsName)) {
            notice(820000, ['标签参数有误']);
        }
        $post['tag_names']   = $tagsName;
        $post['relation_id'] = (int)$post['relation_id'];
        if (empty($post['relation_id'])) {
            notice(820000, ['relation_id参数有误']);
        }
        $result = (new UserCustomizeTagsService())->addNewRelation($post);
        $result = is_array($result) ? $result : [];
        return success(['code' => 200, 'message' => '绑定成功', 'data' => $result]);
    }
}