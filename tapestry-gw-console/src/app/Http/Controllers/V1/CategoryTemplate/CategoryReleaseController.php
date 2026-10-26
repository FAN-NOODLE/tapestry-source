<?php
/**
 * Created by PhpStorm.
 * User: zy
 * Date: 2020/3/19
 * Time: 15:31
 */

namespace App\Http\Controllers\V1\CategoryTemplate;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Category;
use CrsSdk\Shop\CategoryTemplate;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class CategoryReleaseController extends BaseController
{

    /**
     * 分类模板类型 1商城 2门店
     */
    CONST CATEGORY_TPL_CENTER = 1;

    /**
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index(Request $request)
    {

        // 接受参数
        $input = $request->only([
            'cate_tpl_id',
            'store_ids'
        ]);

        // 参数验证规则
        $rules = [
            'cate_tpl_id' => 'required|integer',
            'store_ids'   => 'array'
        ];

        $this->validate($request, $rules);


        $sourceIdArr = $input['store_ids'];
        $tplInfo     = (new CategoryTemplate())->categoryTemplateDetail([
            'enterprise_id' => $this->enterpriseId,
            'id'            => $input['cate_tpl_id']
        ]);
        if ($tplInfo['type'] == self::CATEGORY_TPL_CENTER) {
            $sourceIdArr = [$this->enterpriseId];
        }
        $res = $this->pushCategory($input['cate_tpl_id'], $sourceIdArr);
        return success($res);

    }

    /**
     * 推送分类数据到商品微服务
     *
     * @param $cateTplId
     * @param $storeIdArr
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2020/3/20 14:31
     */
    private function pushCategory($cateTplId, $storeIdArr)
    {
        $categorySdkObj = (new Category());
        $returnArr      = [];
        foreach ($storeIdArr as $storeId) {
            $data['source_id']   = $storeId;
            $data['cate_tpl_id'] = $cateTplId;
            try {
                $categorySdkObj->releaseCategory($data);
                $returnArr['success_store_ids'][] = $storeId;
            } catch (\Error $e) {
                Log::error('门店分类推送失败', [
                    'status' => $e->getCode(),
                    'msg'    => $e->getMessage(),
                    'data'   => $data,
                ]);
                $returnArr['fail_store_ids'][] = $storeId;
            } catch (\Exception $e) {
                Log::error('门店分类推送失败', [
                    'msg'  => $e->getMessage(),
                    'data' => $data,
                ]);
                $returnArr['fail_store_ids'][] = $storeId;
            }
        }

        return $returnArr;
    }

}