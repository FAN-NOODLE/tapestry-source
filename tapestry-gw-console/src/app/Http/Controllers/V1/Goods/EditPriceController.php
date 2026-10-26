<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/25
 * Time: 20:15
 */
namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Services\PostageService;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use CrsSdk\MicroGoodsV1\PostageRelation;
use Illuminate\Support\Facades\Log;

class EditPriceController extends BaseController
{

    public function index()
    {
        $input = $this->validates();
        $this->writeOperatorLog('修改价格');
        try {
            (new MaterialChannel())->edit($input);

        } catch (Error $error) {

            Log::error('获取渠道列表失败', $error->data);
            notice(100000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('获取渠道列表失败' . $e->getMessage(), $input);
            notice(100000, $e->getMessage());

        }

        return success(['data' => '']);
    }

    protected function validates()
    {
        $input = $this->request->all();

        // 参数验证规则
        $rules = [
            'source_id'      => 'required|integer',
            'sku_id'         => 'required|integer',
            'on_line_price'  => 'integer',
            'off_line_price' => 'integer',
        ];

        // 提示语
        $messages = [
            'source_id.*'  => '门店ID不能为空',
            'sku_id.*'     => 'skuId不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        if (!isset($input['on_line_price']) && !isset($input['off_line_price'])) {
            notice(100000, ['请输入价格']);
        }

        $update = [
            'sku_id' => $input['sku_id'],
            'type'   => 'update',
        ];

        if (isset($input['on_line_price'])) $update['on_line_price'] = $input['on_line_price'];
        if (isset($input['off_line_price'])) $update['off_line_price'] = $input['off_line_price'];

        $postageRelation = new PostageRelation();

        // 查询sku的邮费模板
        $thisRelation = $postageRelation->relationList([
            'object_id'     => $input['sku_id'],
            'object_type'   => 1,
            'enterprise_id' => $this->enterpriseId
        ])['list'];

        if (!empty($thisRelation)) {
            $deliveryArr = array_column($thisRelation, 'delivery_type');
            // 计算sku的配送方式
            $update['sale_type'] = (new PostageService())->getSaleTypeByDeliver($deliveryArr);
        } else {
            $update['sale_type'] = 1;
        }

        $input['channel_type'] = 1;
        $input['object_name'] = '普通商品';
        $input['enterprise_id'] = $this->enterpriseId;
        $input['object_id'] = $input['sku_id'];
        $input['data'] = [$update];
        $input['operate_id'] = $this->enterpriseAdminId;

        return $input;
    }
}