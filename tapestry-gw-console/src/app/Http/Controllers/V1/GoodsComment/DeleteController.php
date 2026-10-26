<?php


namespace App\Http\Controllers\V1\GoodsComment;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Comment;

class DeleteController extends BaseController
{
    /**
     * 删除评论
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        return success((new Comment())->deleteComment($input));
    }
}