<?php


namespace App\Http\Controllers\Service\V1\Tag;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\TagService;
use Illuminate\Support\Facades\Log;

class DeleteController extends BaseController
{
    /**
     * 删除标签组或标签
     * @api /service/v1/tag/del
     * tag_id
     * tag_group_id
     * @return array
     */
    public function index()
    {
        $input = $this->request->only([
            'qy_domain',
            'tag_id',
            'tag_group_id'
        ]);

        // 企业标识必传
        if (!isset($input['qy_domain']) || empty($input['qy_domain'])) {
            notice(200018);
        }

        // 调用删除
        try {

            (new TagService())->delTag($input);

        } catch (\Exception $exception) {

            Log::info('删除标签或标签组失败:'.$exception->getMessage(), $input);
            notice(810013);

        }

        return success(['data' => []]);
    }

}