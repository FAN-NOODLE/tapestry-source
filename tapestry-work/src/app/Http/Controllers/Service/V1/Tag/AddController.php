<?php


namespace App\Http\Controllers\Service\V1\Tag;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\TagRequest;
use App\Services\CustomerGroupService;
use App\Services\TagService;
use Illuminate\Support\Facades\Log;

class AddController extends BaseController
{
    /**
     * 标签和标签组添加-添加标签组时必须增加标签
     * @return array
     * @throws \Illuminate\Validation\ValidationException*@throws \Exception
     * @throws \Exception
     * @api /service/v1/tag/add
     */
    public function index()
    {
        $input = $this->request->only([
            'qy_domain',
            'enterprise_id',
            'tag_group_id',
            'tag_group_name',
            'order',
            'tag_list',
        ]);

        // group_id和group_name必须传一个
        if (
            (!isset($input['tag_group_id']) || empty($input['tag_group_id'])) &&
            (!isset($input['tag_group_name']) || trim($input['tag_group_name']) == '')
        ) {
            notice(810004);
        }

        // 验证
        TagRequest::add($input);

        $input['order'] = $input['order'] ?? 0;
        $input['tag_group_id'] = $input['tag_group_id'] ?? 0;
        $input['tag_group_name'] = $input['tag_group_name'] ?? '';

        // 判断标签是否重复--虽然觉得这样判断so foolish
        if (isset($input['tag_group_id']) && !empty($input['tag_group_id'])) {
            (new TagService())->isRepeatTag($input['qy_domain'], $input['tag_group_id'], $input['tag_list']);
        }

        try {

            // 返回true
            (new TagService())->addTag($input);

        } catch (\Exception $exception) {

            Log::info('添加标签失败'.$exception->getMessage(), $input);
            notice(810005);

        }

        return success(['data' => []]);
    }

}