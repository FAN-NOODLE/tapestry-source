<?php


namespace App\Http\Controllers\V1\Postage;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\PostageRelation;
use CrsSdk\Shop\Store;
use Exception;

class RelationStoreListController extends BaseController
{

    /**
     * 获取关联商品列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params                  = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        // 获取商品信息
        $data = (new PostageRelation())->postageRelationStoreList($params);
        $params['store_ids']     = $data['list'];

        try {
            $list = (new Store())->lists($params);
        } catch (Exception $e) {
            notice(453003);
        }

        return success($list);
    }

}