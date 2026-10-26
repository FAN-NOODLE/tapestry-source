<?php


namespace App\Http\Controllers\V1\Qywx\Tag;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\TagRequest;
use CrsSdk\Work\Tag;

class AddController extends BaseController
{
    /**
     * 添加标签组和标签
     * @api /app/v1/qywx/tag/add
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'tag_group_id',
            'tag_group_name',
            'order',
            'tag_list',
        ]);

        TagRequest::add($input);

        $input['qy_domain'] = $this->qyDomain;
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('添加标签和标签组');

        return (new Tag())->add($input);
    }

}