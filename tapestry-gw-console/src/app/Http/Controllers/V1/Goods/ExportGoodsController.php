<?php
/**
 * 批量导出商品
 */

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\GoodsExportJob;

class ExportGoodsController extends BaseController
{
    /**
     * 文件名称（列表显示和导出名称）
     * @var
     */
    const FILE_NAME = [
        -1 => '商品管理-全部商品列表（SPU）',
        1 => '商品管理-全部商品列表（SKU）',
        2 => '商品管理-普通商品列表',
        3 => '商品管理-活动商品列表',
        4 => '商品管理-待上架商品列表',
        5 => '商品管理-已失效商品列表',
    ];

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'goods/export-goods';

    /**
     * 商品列表出tag
     */
    const TAG = 'export-goods';

    /**
     * 商品导出
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        //查询条件校验
        $input = $this->validates();
        $input['enterprise_hash'] = $this->request->header('enterprise-hash');
        $input['enterpriseId']          = $this->enterpriseId;
        $input['enterprise_id']          = $this->enterpriseId;
        $input['route_url']     = $input['route_url'] ?? self::ROUTE_URL;
        $input['tag'] = $input['tag'] ?? self::TAG;
        $input['operate_name']  = $this->enterpriseAdmin['username'] ?? '未知';
        $input['operate_id']    = $this->enterpriseAdminId;
        // 导出
        $this->dispatch(new GoodsExportJob($input, $input, self::FILE_NAME[$input['searchType']]));
        return success();
    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        return $this->validate(request(), [
            'exportChannel' => 'required|integer|in:1,2',
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
            'searchParams.searchReleaseStatus'  => 'integer',
            'searchParams.searchNeedDelete'     => 'integer',
            'searchParams.searchChannelOwnership'   => 'integer',
            'searchParams.searchSort.sortColumn'    => 'integer',
            'searchParams.searchSort.sortType'      => 'integer',
            'searchParams.channelIds'               => 'array',
            'page'  => 'integer',
            'size'  => 'integer'
        ], [
            'exportChannel.*'   => '商品导出渠道不合法',
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
            'searchParams.searchReleaseStatus.*'  => '商品发布状态不合法',
            'searchParams.searchNeedDelete.*'     => '商品删除状态不合法',
            'searchParams.searchChannelOwnership.*' => '渠道归属不合法',
            'searchParams.searchSort.sortColumn.*'  => '排序字段不合法',
            'searchParams.searchSort.sortType.*'    => '排序类型不合法',
            'searchParams.channelIds.*'             => '多选渠道ID不合法',
            'page.*'    => '分页不合法',
            'size.*'    => '分页大小不合法',
        ]);
    }
}
