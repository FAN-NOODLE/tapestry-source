<?php
/**
 * Created by PhpStorm.
 * User: 刘鹏伟
 * Date: 2020/3/20
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

class CateSkuListController extends BaseController
{
    /**
     * sku售卖方式 1：仅线上售卖
     */
    const ON_LINE_SALE = 1;

    /**
     * sku售卖方式 2：仅线下自提
     */
    const OFF_LINE_SALE = 2;

    /**
     * sku售卖方式 3：线上线下均有售
     */
    const ALL_SALE = 3;

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        //查询条件校验
        $input = $this->validates();

        // 根据搜索类型设置查询条件
        if (!empty($input['keywords']) && !empty($input['search_key'])) {
            switch ($input['search_key']) {
                case 'name':
                    $input['name'] = $input['keywords'];
                    break;
                case 'sku_id':
                    $input['sku_id'] = (int)$input['keywords'];
                    break;
                case 'parent_sku_id':
                    $input['parent_sku_id'] = (int)$input['keywords'];
                    break;
            }
            unset($input['search_key']);
        }

        $input['enterprise_id'] = $this->enterpriseId;

        try {

            $list = (new Goods())->sourceSkuList($input);

            $this->format($list['list']);

        } catch (Error $error) {

            Log::info('获取商品列表失败', $error->data);
            notice(451508, [], $error->data);

        } catch (\Exception $exception) {

            Log::info('获取商品列表失败' . $exception->getMessage(), $input);
            notice(451508);
        }

        //返回结果
        return success($list);
    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'source_id'      => 'required|int',
            'page'           => 'integer',
            'size'           => 'integer',
            'search_key'     => 'string',
            'keywords'       => 'string',
            'channel_status' => 'integer',
            'category_id'    => 'integer'
        ];

        $message = [
            'source_id.required'     => '请指定对应门店',
            'source_id.integer'      => '门店ID数据类型不正确',
            'page.integer'           => '页数有误',
            'size.integer'           => '每页记录有误',
            'search_key.string'      => '查询类型有误',
            'keywords.string'        => '排序字段有误',
            'channel_status.integer' => '上下架类型有误',
            'category_id.integer'    => '分类ID有误'
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }

    /**
     * @param $res
     * @return array
     */
    private function format(&$res)
    {
        if (!$res) return [];

        foreach ($res as &$v) {
            if (in_array($v['sale_type'], [self::ON_LINE_SALE, self::ALL_SALE])) {
                // 线上或线上线下都有售，则展示线上价格
                $price = $v['on_line_price'];
            } else {
                // 仅线下售卖，则展示线下价格
                $price = $v['off_line_price'];
            }
            $v = [
                'parent_sku_id'          => $v['parent_sku_id'],
                'material_id'            => $v['material_id'],
                'material_channel_id'    => $v['material_channel_id'],
                'sku_id'                 => $v['id'],
                'sku_img'                => $v['sku_img'] ?? '',
                'sku_name'               => $v['sku_name'] ?? '',
                'category_name'          => is_array($v['category_name']) ? implode('/', $v['category_name']) : '',
                'category_list'          => $v['category_list'] ?? [],
                'channel_status'         => $v['channel_status'] ?? 1,
                'guide_price'            => $v['guide_price'],
                'on_line_price'          => $price,
                'on_line_current_stock'  => $v['on_line_current_stock'],
                'off_line_current_stock' => $v['off_line_current_stock'],
            ];
        }

        return $res;
    }
}