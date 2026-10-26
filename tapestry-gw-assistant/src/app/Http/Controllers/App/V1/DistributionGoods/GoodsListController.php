<?php


namespace App\Http\Controllers\App\V1\DistributionGoods;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\MicroGoodsV1\MiniProgram;
use CrsSdk\Work\DistributionGoods;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;

class GoodsListController extends BaseController
{

    /**
     * 图片路径
     */
    const test_url = "https://img-crs-test.vchangyi.com/";
    const prod_url = "https://img-crs.vchangyi.com/";
    /**
     * 商品隐藏分销商品展示兼容
     * @var int
     */
    protected  $show_hide = 1;

    /**
     * 推广佣金类型：1=>百分比,2=>金额;
     */
    const PROMOTION_COMMISSION_TYPE_PERCENTAGE = 1;
    const PROMOTION_COMMISSION_TYPE_AMOUNT = 2;

    /**
     * 服务佣金类型：1=>百分比,2=>金额;
     */
    const SERVICE_COMMISSION_TYPE_PERCENTAGE = 1;
    const SERVICE_COMMISSION_TYPE_AMOUNT = 2;

    /**
     * 二级佣金设置
     * 1: 基于商品分销
     */
    const COMMISSION_TYPE_GOODS = 1;

    /**
     * 二级佣金设置
     * 2: 基于推广佣金总额的分账
     */
    const COMMISSION_TYPE_PROMOTION = 2;
    /**
     * 分销商展示 spu维度 修改成sku
     * 为兼容原goods模块商品依spu维度返回商品信息  故查询分销商品依spu维度
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/21
     */
    public function index()
    {

        $params = request()->only(['category_ids', 'new_goods', 'search_value', 'sort_sales', 'sort_commission', 'page', 'size', 'store_id']);
        $params['page'] = $params['page'] ?? 1;
        $params['size'] = $params['size'] ?? 10;

        $return = [
            'list' => [],
            'page' => intval($params['page']),
            'size' => intval($params['size']),
            'total' => 0,
        ];
        return success($return);

        /**
        $params['enterprise_id'] = $this->enterpriseId;
        $params['new_goods'] = $params['new_goods'] ?? 0;
        $distributionResult = (new DistributionGoods())->appGoodsSkuList($params);
        $distributionGoods = $distributionResult['list'] ?? [];
        $setting = $distributionResult['setting'] ?? [];
        $total = $distributionResult['total'] ?? 0;
        if (empty($distributionGoods)) {
            return success([]);
        }

        $spuIds = array_keys($distributionGoods);
        $where = [
            'enterprise_id' => $this->enterpriseId,
            //'page' => $params['page'],
            //'size' => $params['size'],
            'sort_field' => 'sales',
            'sort' => 'DESC',
            'spu_ids' => $spuIds
        ];
        if (!empty($params['category_ids'])) {
            $where['category_id'] = $params['category_ids'];
        }
        if (!empty($params['search_value'])) {
            $where['keywords'] = $params['search_value'];
        }
            $where['page'] = 1;
            $where['size'] = 9999;

        $where['show_hide'] = $this->show_hide;
        $goodsList = (new MiniProgram())->searchGoodsSpuList($where);
        $return = [
            'list' => [],
            'page' => intval($params['page']),
            'size' => intval($params['size']),
            'total' => 0,
        ];
        return success($return);

        $appId = '';
        $traceParams = [];
        if ($this->userType == self::USERTYPE_WORKER) {
            // 获取商城appid
//            try {
//                $shopDetail = (new \CrsSdk\Boss\MiniProgram())->detail([
//                    'enterprise_id' => $this->enterpriseId
//                ]);
//            } catch (Error $error) {
//                if ($error->getCode() == 300100) {
//                    notice(657012);
//                }
//            }
            $appId = $shopDetail['data']['app_id'] ?? '';

            // 获取导购workerid
             $memberInfo = $this->storeWorkerDetail($this->enterpriseId, $params['store_id'] ?? 0);
            $traceParams['mem_uid'] = $memberInfo['mem_uid'] ?? '';
            $traceParams['share_uid'] = $memberInfo['mem_uid'] ?? '';
            $traceParams['source'] = self::SHARE_TYPE;
        }
        foreach ($goodsList['list'] as $goods) {
            $skuList = $goods['sku']['list'] ?? [];
            array_multisort(array_column($skuList, 'sell_price'), SORT_DESC, $skuList);
            $skuIDs = array_column($distributionGoods[$goods['spu_id']]['sku_list'],'sku_id');
            //$sku = current($skuList);
            foreach ($skuList as $k=>$v){
               if(in_array($v['id'],$skuIDs)){//spu下sku 是否为分销商品
                   $commission_value = $this->getCommissionInfo($v,$setting,$distributionGoods,$goods);
                   $traceId = 0;
                   if ($this->userType == self::USERTYPE_WORKER) {

                       $traceParams['spu_id'] = $v['spu_id'];
                       $traceParams['sku_id'] = $v['id'];
                       $traceId = encryptShortLink($traceParams);
                   }
                   $return['list'][] = [
                       'spu_id' => $v['spu_id'],
                       'spu_img' => $v['pictures']['image'],
                       'spu_name' => $v['sku_name'],
                       'sales' => $v['sales'] ?? 0,
                       'commission' => $commission_value,
                       'new_goods' => $distributionGoods[$goods['spu_id']]['sku_list'][$v['id']]['new_goods'] ?? 0,
                       'sku_id' => $v['id'] ?? 0,
                       'sign_price' => $v['sign_price'] ?? 0,
                       'sell_price' => $v['sell_price'] ?? 0,
                       'app_id' => $appId,
                       'path' => 'goods/goods-detail/goods-detail.html?short_param_id=' . $traceId
                   ];
               }
            }
        }
        $return['total'] = count($return['list']);
        if (!empty($params['sort_commission'])) {
            $sort = $params['sort_commission'] == 'DESC' ? SORT_DESC : SORT_ASC;
            array_multisort(array_column($return['list'], 'commission'), $sort, $return['list']);
        }
        if (!empty($params['sort_sales'])) {
            $sort = $params['sort_sales'] == 'DESC' ? SORT_DESC : SORT_ASC;
            array_multisort(array_column($return['list'],'sales'),$sort,$return['list']);
        }
        $start=($params['page'] - 1) * $params['size'];
        $return['list'] = array_slice($return['list'], $start, $params['size']);
        $env = config('app.env');
        if ($env === 'dev' || $env === 'test') {
            foreach ($return['list'] as $key => $val) {
                $val['spu_img'] = $this->imgUrl($val['spu_img']);
                $return['list'][$key] = $val;
            }
        }
        return success($return);
         */
    }

    public function getCommissionInfo($value,$setting,$distributionGoods,$goods)
    {
        $commission = [];
        if(!empty($setting) && $this->userType == self::USERTYPE_PARTNER){//设置了二级分销
                //基于商品分销
                if($setting['commission_type'] == self::COMMISSION_TYPE_GOODS){
                    //合伙人推广佣金
                    $commission['partner_commission'] = $this->priceCalculation($value['sell_price'],$distributionGoods[$goods['spu_id']]['sku_list'][$value['id']]['commission']);
                    //导购推广佣金
                    //$commission['promotion_commission'] = (int)floor(($value['sell_price'] / 100) * ($setting['commission_rate'] / 10000) * 100);
                }
                //基于推广佣金总额的分账
                if($setting['commission_type'] == self::COMMISSION_TYPE_PROMOTION){
                    //总的推广佣金
                    $promotion_commission =  $this->priceCalculation($value['sell_price'],$distributionGoods[$goods['spu_id']]['sku_list'][$value['id']]['commission']);
                    //合伙人推广佣金
                    $commission['partner_commission'] = $this->priceCalculation($promotion_commission,$setting['commission_rate']);
                    //导购推广佣金 = 总推广佣金-合伙人佣金
                    //$commission['promotion_commission'] = (int)floor($promotion_commission - $commission['partner_commission']);
                }
            }
        //未设置或二级分销关闭
        if($this->userType == self::USERTYPE_WORKER || empty($setting)){
            $commission['promotion_commission'] = $distributionGoods[$goods['spu_id']]['sku_list'][$value['id']]['commission'] / 10000 * $value['sell_price'] ?? 0;
        }

        return $commission_value = isset($commission['partner_commission']) ? $commission['partner_commission'] : $commission['promotion_commission'];
    }

    /**
     * 价格精度计算
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/25 15:01
     */
    public function priceCalculation($order_goods_price,$rate)
    {
        return (int)floor(bcmul($order_goods_price,$rate,2)/10000);
    }

    /**
     * 图片域名替换
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/5 18:19
     */
    public function imgUrl($img)
    {
        if(empty($img)){
            return '';
        }
        // 如果是开发和测试环境
            $imgArray = explode('/',$img);
            if($imgArray[2] == "crs-dev-1251064102.cos.ap-shanghai.myqcloud.com"){
                $img = self::test_url.$imgArray[3];
            }
            return $img;
        }
}