<?php


namespace App\Http\Controllers\V1\Qywx\Tag;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Tag;
use Illuminate\Http\JsonResponse;

class ListController extends BaseController
{
    /**
     * 返回标签组列表
     * @return mixed
     * @throws \Exception
     * @api /app/v1/qywx/tag/list
     */
    public function index()
    {
        $input['qy_domain'] = $this->qyDomain;

        return (new Tag())->tagGroupList($input);
    }

    /**
     * 获取标签列表
     * @return mixed
     * @throws \Exception
     * @api /app/v1/qywx/tag/tag-list
     */
    public function tagList()
    {
        $input = $this->request->only([
            'tag_group_id',
            'keyword', // 标签名称
            'custom_group_id', // 用户分组id
            'page',
            'size'
        ]);

        $input['qy_domain']     = $this->qyDomain;
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Tag())->tagList($input);

        return success([
            'list'  => $list['data'],
            'page'  => $list['meta']['page'],
            'size'  => $list['meta']['size'],
            'total' => $list['meta']['total'],
        ]);
    }

    /**
     * 获取企业标签库所有标签组和标签
     * @return JsonResponse
     * @throws \Exception
     */
    public function allList()
    {
        $list = (new Tag())->allList(['qy_domain' => $this->qyDomain]);

        return success(['list' => $list['data']]);
    }

    /**
     * 获取企业下所有标签列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function selectList()
    {
        $input              = $this->request->input();
        $input['qy_domain'] = $this->qyDomain;

        $list = (new Tag())->tagSelectList($input);

        return success([
            'list'  => $list['data'],
            'page'  => $list['meta']['page'],
            'size'  => $list['meta']['size'],
            'total' => $list['meta']['total'],
        ]);
    }
}