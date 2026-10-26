<?php

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

class ListController extends BaseController
{

    /**
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index()
    {
        //查询条件校验
        // $input = $this->validates();
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $sourceIds = $input['search']['source_ids'] ?? [];

        // 是否有商城权限
        if ($this->hasShopAuth == 1) {
            $this->storeIds = array_merge($this->storeIds, [$this->enterpriseId]);
        }

        if ($sourceIds === []) {
            $sourceIds = $this->storeIds;
        } else {
            // 如果前端传门店ID但是不在该管理员可查看的范围,防止跨权限查看
            $sourceIds = array_intersect($this->storeIds, $sourceIds);
            if ($sourceIds === [] || $this->hasShopAuth === 2) {
                return success([
                    'list'  => [],
                    'total' => 0,
                    'page'  => $input['page'] ?? 1,
                    'size'  => $input['size'] ?? 10
                ]);
            }
        }

        try {
            $input['search']['source_ids'] = array_values($sourceIds);

            // 如果是获取待上架商品，则排除活动商品
            if($input['type'] == 4){
                $input['search']['channel_type'] = [1];
            }

            Log::error('获取商品列表搜素条件', $input);
            $data = (new Goods())->list($input);
        } catch (Error $error) {

            Log::error('获取商品列表失败', $error->data);
            notice(451508, [], $error->data);

        } catch (\Exception $exception) {

            Log::error('获取商品列表失败'.$exception->getMessage(), $input);
            notice(451508);
        }

        //返回结果
        return success($data);
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
            'page'                 => 'integer',
            'size'                 => 'integer',
            'type'                 => 'required|integer',
            'sort_field'           => 'integer',
            'sort'                 => 'string',
            'search'               => 'array',
            'search.sku_id_list'   => 'array',
            'search.source_ids'    => 'array',
            'search.name_type'     => 'string',
            'search.keywords'      => 'string',
            'search.between_type'  => 'string',
            'search.start_between' => 'string',
            'search.end_between'   => 'string',
            'search.time_type'     => 'string',
            'search.start_time'    => 'integer',
            'search.end_time'      => 'integer',
        ];

        $message = [
            'page.integer'                => '页数有误',
            'size.integer'                => '每页记录有误',
            'type.*'                      => '查询类型有误',
            'sort_field.string'           => '排序字段有误',
            'sort.string'                 => '排序方式有误',
            'search.array'                => '搜索条件格式有误',
            'search.sku_id_list.array'    => 'sku集合有误',
            'search.source_ids.array'     => '门店ID有误',
            'search.name_type.string'     => '搜索条件有误',
            'search.keywords.string'      => '搜索关键词有误',
            'search.between_type.string'  => '范围搜索类型有误',
            'search.start_between.string' => '开始值有误',
            'search.end_between.string'   => '结束值有误',
            'search.time_type.string'     => '搜索时间范围类型',
            'search.start_time.integer'   => '开始时间有误',
            'search.end_time.integer'     => '结束时间有误',
        ];

        $this->validate(request(), $rules, $message);

        return request()->input();
    }
}