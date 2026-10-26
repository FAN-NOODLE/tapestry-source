<?php


namespace App\Http\Controllers\V1\GoodsComment;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Comment;

class StatusController extends BaseController
{
    /**
     * 评论上下架
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        return success((new Comment())->statusComment($input));
    }
}