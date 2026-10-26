<?php
/**
 * Created by PhpStorm.
 * User: liupengwei
 * Date: 2020/5/14
 * Time: 15:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SpuOptionListController extends BaseController
{
    /**
     * SPU选项框列表接口标识
     */
    const IS_SIMPLE_SPU_LIST = 1;

    /**
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index(Request $request)
    {
        $input = $this->request->all();
        $input['page']          = $request->input('page', Constant::PAGE);
        $input['size']          = $request->input('size', Constant::LIMIT);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['is_simple_list'] = self::IS_SIMPLE_SPU_LIST;

        try {

            $list = (new Goods())->spuList($input);

        } catch (Error $error) {
            Log::info('商品列表获取失败', $error->data);
            notice(100000, [$error->getMessage()], $error->data);
        } catch (\Exception $e) {
            Log::info('商品列表获取失败' . $e->getMessage(), $input);
            notice(100000, [$e->getMessage()]);
        }

        return success($list);
    }
}
