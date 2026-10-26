<?php


namespace App\Http\Controllers\V1\Postage;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\PostageRelation;
use Illuminate\Support\Facades\Log;

class RelationActivityListController extends BaseController
{
    /**
     * 获取关联活动列表
     *
     * @return \Illuminate\Http\JsonResponse
     *@api {domain}/gw-console/v1/postage/relation-activity-list
     */
    public function index()
    {
        $input = $this->request->all();

        $input['enterprise_id'] = $this->enterpriseId;

        try {
            $result = (new PostageRelation())->postageRelationActivityList($input);
        } catch (\Exception $exception) {
            Log::info('获取关联活动列表失败:sdk->PostageRelation:'.$exception->getMessage(), $input);
            notice(921009, ['活动']);
        }

        return success([
            'list' => $result['list'] ?? [],
            'page' => $result['page'] ?? $input['page'],
            'size' => $result['size'] ?? $input['size'],
            'total' => $result['total'] ?? 0,
            'activity_type_list' => $result['activity_type_list'] ?? 0
        ]);

    }

}