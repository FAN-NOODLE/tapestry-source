<?php
namespace App\Http\Controllers\V1\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\PostageRelation;
use Illuminate\Support\Facades\Log;

class PostageTypeRelationController extends BaseController
{
    /**
     * 关联邮费模板
     * @api {domain}/gw-console/v1/postage/create-relation
     * @return array
     */
    public function index()
    {
        $input = $this->request->only([
            'object_type',
            'object_id',
            'relation'
        ]);
        $input['enterprise_id'] = $this->enterpriseId;
        try {
            (new PostageRelation())->postageRelationSkuOrActivity($input);
        } catch (\Exception $exception) {
            Log::info('添加邮费模板失败:sdk->PostageRelation:'.$exception->getMessage(), $input);
            notice(921008);
        }

        return success();
    }

}