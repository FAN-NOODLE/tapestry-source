<?php


namespace App\Http\Controllers\V1\Goods;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Goods;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class ConsoleGoodsListController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws Exception
     */
    public function index(): JsonResponse
    {
        $input                  = $this->getParams();
        $input['enterpriseId']  = $this->enterpriseId;
        return success((new Goods())->getConsoleList($input));
    }

    /**
     * @return array
     * @throws ValidationException
     */
    public function getParams(): array
    {
        return $this->validate(request(), [
            'searchType'    => 'required|integer',
            'searchParams.searchString.keyword' => 'string',
            'searchParams.searchString.type'    => 'integer',
            'searchParams.searchInt.high'       => 'integer',
            'searchParams.searchInt.low'        => 'integer',
            'searchParams.searchInt.type'       => 'integer',
            'searchParams.searchChannel'        => 'integer',
            'searchParams.searchCreateTime.beginTimestamp'  => 'integer',
            'searchParams.searchCreateTime.endTimestamp'    => 'integer',
            'searchParams.searchGoodsStatus'    => 'integer',
            'searchParams.searchChannelOwnership'   => 'integer',
            'searchParams.searchSort.sortColumn'    => 'integer',
            'searchParams.searchSort.sortType'      => 'integer',
            'page'  => 'integer',
            'size'  => 'integer'
        ], [
            'searchType.*'      => '列表类型不合法',
            'searchParams.searchString.keyword.*'   => '关键字必须是字符串',
            'searchParams.searchString.type.*'      => '关键字搜索类型不合法',
            'searchParams.searchInt.high.*'     => '最大值不合法',
            'searchParams.searchInt.low.*'      => '最小值不合法',
            'searchParams.searchInt.type.*'     => '范围类型不合法',
            'searchParams.searchChannel.*'      => '渠道不合法',
            'searchParams.searchCreateTime.beginTimestamp.*'    => '开始时间不合法',
            'searchParams.searchCreateTime.endTimestamp.*'      => '结束时间不合法',
            'searchParams.searchGoodsStatus.*'  => '商品状态不合法',
            'searchParams.searchChannelOwnership.*' => '渠道归属不合法',
            'searchParams.searchSort.sortColumn.*'  => '排序字段不合法',
            'searchParams.searchSort.sortType.*'    => '排序类型不合法',
            'page.*'    => '分页不合法',
            'size.*'    => '分页大小不合法',
        ]);
    }
}