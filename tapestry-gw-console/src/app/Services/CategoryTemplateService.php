<?php
/**
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/9/12
 * Time: 16:52
 */

namespace App\Services;

use App\Common\Constant;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\Shop\CategoryTemplate;

class CategoryTemplateService extends BaseService
{

    /**
     * 商品列的名称
     */
    CONST IMPORT_GOODS_ROWS_NAME = '关联商品SPUid';

    /**
     * 表头在excel中的行
     */
    CONST IMPORT_HEAD_LINE = 1;

    /**
     * 分类导入时数据最小列
     */
    CONST MIN_IMPORT_ROW = 2;

    /**
     * 模板错误
     */
    CONST ROOT_CATEGORY_NOT_EXISTS = '顶级分类不存在';
    CONST SPU_ID_ERROR = '请检查商品spuId';

    /**
     * 最大分页条数
     */
    const  LIST_MAX_PAGE_SIZE = 500;

    /**
     * 处理导入数据并提交导入服务
     *
     * @param $excelData
     * @param $input
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2020/3/24 17:17
     */
    public function importCategory($excelData, $input)
    {
        //表头信息
        $header = $excelData[self::IMPORT_HEAD_LINE];
        //初始化商品列
        $goodsRow = 0;
        //获取商品列
        foreach ($header as $key => $val) {
            if ($val == self::IMPORT_GOODS_ROWS_NAME) {
                $goodsRow = $key;
            }
            //删除多余的表头列
            if ($goodsRow && $key > $goodsRow) {
                unset($header[$key]);
            }
        }

        //商品列不存在报错
        if (!$goodsRow) {
            notice(100000, ['excel表头异常，请不要修改表头']);
        }
        //获取商品id集合并删除多余的列
        $goodsSpuIdArr = [];
        foreach ($excelData as $key => $val) {
            //获取所有商品spuId
            if (isset($val[$goodsRow]) && $val[$goodsRow]) {
                if ($val[$goodsRow] != self::IMPORT_GOODS_ROWS_NAME) {
                    $goodsSpuIdArr[] = $val[$goodsRow];
                }
            }
            //删除多余的列
            foreach ($val as $k => $v) {
                if ($k > $goodsRow) {
                    unset($excelData[$key][$k]);
                }
            }
            //补全不存在的列
            for($i = 0; $i <= $goodsRow; $i ++) {
                $excelData[$key][$i] = $excelData[$key][$i] ?? '';
            }
        }
        $success   = 0; //成功的数据条数
        $error     = 0; //失败的数据条数
        $errorData = []; //失败的数据集合

        $cateTplSdk = new CategoryTemplate();
        try {
            $tplInfo = $cateTplSdk->categoryTemplateDetail([
                'enterprise_id' => $input['enterprise_id'],
                'id'            => $input['cate_tpl_id'],
            ]);
            // 判断数据与模板层级是否匹配
            // 数据至少2列（1级分类模板：1列分类 1列数据，2级分类模板：2列分类 1列商品，以此类推）
            $listLevel = count($excelData[1]);
            if ($listLevel < self::MIN_IMPORT_ROW || $listLevel != ($tplInfo['level'] + 1)) {
                notice(100000, ['模板数据与模板层级不符']);
            }
        } catch (\Exception $e) {
            notice(100000, [$e->getMessage()]);
        }
        unset($excelData[0]);
        unset($excelData[1]);
        //获取所有可用的商品spuId集合
        try {
            $result      = (new Goods())->spuList([
                'enterprise_id' => $input['enterprise_id'],
                'source_ids'    => [$input['enterprise_id']],
                'spu_ids'       => $goodsSpuIdArr,
                'page'          => 1,
                'size'          => 1000,
            ]);
            $existsIdArr = [];
            if (!empty($result['list'])) {
                $existsIdArr = array_column($result['list'], 'id');
            }
            foreach ($excelData as $key => $val) {
                if (empty($val[0])) {
                    //一级分类不存在
                    $errorData[] = array_merge($val, [self::ROOT_CATEGORY_NOT_EXISTS]);
                    $error++;
                    unset($excelData[$key]);
                    continue;
                }
                if ($val[$goodsRow] && !in_array($val[$goodsRow], $existsIdArr)) {
                    //商品spuId错误
                    $errorData[] = array_merge($val, [self::SPU_ID_ERROR]);
                    unset($excelData[$key]);
                    $error++;
                    continue;
                }
                $success++;
            }
        } catch (\Exception $e) {
            notice(100000, [$e->getMessage()]);
        }

        if (!empty($excelData)) {
            $data = [
                'cate_tpl_id' => $input['cate_tpl_id'],
                'list'        => array_values($excelData),
            ];
            //导入数据
            try {
                $cateTplSdk->importCategory($data);
            } catch (\Exception $e) {
                notice(100000, [$e->getMessage()]);
            }
        }
        return [
            'success_num' => $success, //成功的数据条数
            'error_num'   => $error, //失败的数据条数
            'error_title' => array_merge($header, ['错误原因']),
            'error_data'  => $errorData //失败的数据集合
        ];
    }

    /**
     * 分类绑定的商品
     * @param $condition
     * @return array
     * @throws \Exception
     */
    public function getCategorySpuListAll($condition)
    {
        $page     = $condition['page'] ?? Constant::PAGE;
        $page_max = 0;
        $list     = [];

        do {
            // 绑定的商品ID
            $result = (new CategoryTemplate())->categoryGoodsList($condition);

            // 获取商品信息
            if (isset($result['list']) && !empty($result['list'])) {
                $list = array_merge($list, $result['list']);
            }

            // 计算总页数
            if ($page_max === 0) {
                $page_max = ceil($result['total'] / $condition['size']);
            }

            $page++;
        } while ($page <= $page_max);

        return $list;

    }
}
