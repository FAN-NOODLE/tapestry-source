<?php


namespace App\Http\Controllers\V1\Qywx\Tag;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Tag;

class DeleteController extends BaseController
{
    /**
     * 删除标签标签组
     * @api /app/v1/qywx/tag/del
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'tag_id',
            'tag_group_id'
        ]);

        $input['qy_domain'] = $this->qyDomain;
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('删除标签和标签组');

        return (new Tag())->delete($input);
    }

}