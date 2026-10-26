<?php


namespace App\Http\Controllers\V1\Qywx\Tag;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Tag;

class TagCustomerGroupController extends BaseController
{
    /**
     * 标签关联用户分组
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input  = $this->request->only([
            'tag_id',
            'tag_group_id',
            'custom_group_ids'
        ]);
        $input['enterprise_id'] = $this->enterpriseId;

        return (new Tag())->tagRelationCustomerGroup($input);
    }

}