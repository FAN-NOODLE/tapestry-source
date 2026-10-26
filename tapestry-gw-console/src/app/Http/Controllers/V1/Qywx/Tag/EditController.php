<?php


namespace App\Http\Controllers\V1\Qywx\Tag;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Tag;

class EditController extends BaseController
{
    /**
     * 编辑标签组|标签
     * @api /app/v1/qywx/tag/edit
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'name',
            'tag_group_id',
            'type', // 1- 编辑标签组 2-编辑标签
        ]);

        $input['qy_domain'] = $this->qyDomain;
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('编辑标签和标签组');

        return (new Tag())->update($input);
    }

}