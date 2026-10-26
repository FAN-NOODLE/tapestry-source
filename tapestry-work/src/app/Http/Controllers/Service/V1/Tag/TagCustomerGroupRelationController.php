<?php


namespace App\Http\Controllers\Service\V1\Tag;


use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\TagRequest;
use App\Services\CustomerGroupService;
use App\Services\TagService;
use Illuminate\Support\Facades\Log;

class TagCustomerGroupRelationController extends BaseController
{
    /**
     * 关联标签和用户分组-添加和更新
     * @api /service/v1/tag/tag-customer
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id',
            'tag_id',
            'tag_group_id',
            'custom_group_ids'
        ]);

        // 验证
        TagRequest::tagRelation($input);

        // 插入
        try {
            $id = (new TagService())->addTagCustomerGroupRelation($input);
        } catch (\Exception $exception) {
            Log::info('关联用户分组失败'.$exception->getMessage(), $input);
            notice(810006);
        }

        // 发送队列
        $data = [
            'id' => $id,
            'enterprise_id' => $input['enterprise_id'],
            'tag_group_id' => $input['tag_group_id'],
            'tag_id' => $input['tag_id'],
            'custom_group_ids' => $input['custom_group_ids'],
        ];

        (new CustomerGroupService())->sendCMQMsg($data, 'bind_customer_tag_group');

        return success(['data' => []]);
    }
}