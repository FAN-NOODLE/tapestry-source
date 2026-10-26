<?php
/**
 * Created by PhpStorm.
 * User: zy
 * Date: 2020/3/4
 * Time: 15:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Services\CategoryTemplateService;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SpuListController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $input                  = $this->getParams($request);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['page']          = $request->input('page', Constant::PAGE);
        $input['size']          = $request->input('size', Constant::LIMIT);
        $input['source_ids']    = [$this->enterpriseId];
        $input['is_selection']  = Constant::IS_SPU_SELECTION;
        $input['search_type']   = 5;


        try {

            $cateCond = [
                'cate_id'       => $input['cate_id'],
                'enterprise_id' => $this->enterpriseId,
                'page'          => Constant::PAGE,
                'size'          => CategoryTemplateService::LIST_MAX_PAGE_SIZE
            ];
            // 获取已绑定的商品
            $cateSpuList = (new CategoryTemplateService())->getCategorySpuListAll($cateCond);

            // 获取商品信息
            if (!empty($cateSpuList)) {
                $input['spu_ids'] = array_unique(array_column($cateSpuList, 'spu_id'));
            }

            $list = (new Goods())->spuList($input);

        } catch (Error $error) {
            Log::info('分类关联商品列表失败', $error->data);
            notice(100000, [$error->getMessage()], $error->data);
        } catch (\Exception $e) {
            Log::info('分类关联商品列表失败' . $e->getMessage(), $input);
            notice(100000, [$e->getMessage()]);
        }

        return success($list);
    }

    /**
     * 验证数据
     * @param $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    protected function getParams($request)
    {
        // 验证参数
        $rule = [
            'cate_id' => 'required|integer',
            'search'  => 'string',
        ];

        $messages = [
            'cate_id.*' => '分类不存在',
        ];

        return $this->validate($request, $rule, $messages);

    }

    /**
     * SPU列表新
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function lists()
    {
        $input                  = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        return success((new Goods())->spuList($input));
    }
}
