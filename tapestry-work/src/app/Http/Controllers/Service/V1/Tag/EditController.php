<?php


namespace App\Http\Controllers\Service\V1\Tag;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\TagRequest;
use App\Services\TagService;
use Illuminate\Support\Facades\Log;

class EditController extends BaseController
{
    /**
     * 编辑标签组或者标签 1-标签组 2-标签
     */
    const EDIT_TAG_GROUP = 1;
    const EDIT_TAG = 2;
    /**
     * 编辑标签或标签组
     * @return array
     * @throws \Illuminate\Validation\ValidationException*@throws \Exception
     * @throws \Exception
     * @api /service/v1/tag/edit
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'name',
            'tag_group_id',
            'type', // 1- 编辑标签组 2-编辑标签
            'qy_domain'
        ]);

        // 验证
        TagRequest::edit($input);

        // 判断标签是否重复--虽然觉得这样判断so foolish
        if (isset($input['tag_group_id']) && !empty($input['tag_group_id']) && $input['type'] == self::EDIT_TAG) {
            (new TagService())->isRepeatTag($input['qy_domain'], $input['tag_group_id'], [
                    'id' => $input['id'],
                    'name' => $input['name']
            ], 'edit');
        }

        // 编辑
       try {

         (new TagService())->editTag($input);

        } catch (\Exception $exception) {

            Log::info('编辑标签失败'.$exception->getMessage());
           throw new Error($exception->getCode(), $exception->getMessage());

        }

        return success(['data' => []]);
    }
}