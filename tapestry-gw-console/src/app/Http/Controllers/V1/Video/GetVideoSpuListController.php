<?php


namespace App\Http\Controllers\V1\Video;


use App\Http\Controllers\Controller;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\Video;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class GetVideoSpuListController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $params = $this->getParams();
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new Video())->getVideoSpuList($params));
    }

    /**
     * 获取参数
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate($this->request, [
            'spu_ids' => 'array',
            'spu_id' => 'integer',
            'status' => 'integer',
            'video_status' => 'integer',
            'video_cate_id' => 'integer',
            'spu_name' => 'string',
            'page' => 'integer',
            'size' => 'integer',
        ], [
            'spu_ids.*' => '参数[spu_ids]不合法',
            'spu_id.*' => '参数[spu_id]不合法',
            'status.*' => '参数[status]不合法',
            'video_status.*' => '参数[video_status]不合法',
            'video_cate_id.*' => '参数[video_cate_id]不合法',
            'spu_name.*' => '参数[spu_name]不合法',
            'page.*' => '参数[page]不合法',
            'size.*' => '参数[size]不合法'
        ]);
    }
}