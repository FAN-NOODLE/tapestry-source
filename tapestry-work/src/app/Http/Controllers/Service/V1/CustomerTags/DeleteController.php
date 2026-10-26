<?php
/**
 * Description: 删除用户自定义标签
 * User: lintong
 * Email: 18511245965@163.com
 * Date: 2020/8/2
 * Time: 17:33
 */

namespace App\Http\Controllers\Service\V1\CustomerTags;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\UserCustomizeTagsRelation;
use App\Services\UserCustomizeTagsService;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    /**
     * 解除自定义标签绑定
     * @param Request $request
     * @return array
     */
    public function DeleteRelation(Request $request)
    {
        $post = $request->only('id');
        try {
            UserCustomizeTagsRelation::query()
                ->where('id', $post['id'])
                ->delete();
            return success(['code' => 200, 'message' => '操作成功！', 'data' => ['id' => $post['id']]]);
        } catch (\Exception $e) {
            error(820006);
        }
    }

    /**
     * 删除店员自定义标签
     * @param Request $request
     * @return array
     */
    public function DeleteTag(Request $request)
    {
        $post = $request->only('id');
        (new UserCustomizeTagsService())->deleteTag($post);
        return success(['code' => 200, 'message' => '操作成功！', 'data' => ['id' => $post['id']]]);
    }
}