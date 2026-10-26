<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/20 0020
 * Time: 11:07
 */

namespace App\Http\Controllers\V1\ActivityCommon;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;
use CrsSdk\MicroGoodsV1\Activity as GoodsActivity;
use App\Exceptions\Error;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use Illuminate\Support\Facades\Log;

class ActivityBaseController extends BaseController
{
    /**
     * 活动sdk
     *
     * @var Activity|null
     */
    protected $activityApi = null;

    /**
     * 活动商品sdk
     *
     * @var GoodsActivity|null
     */
    protected $goodsActivityApi = null;

    /**
     * 输入参数
     *
     * @var array|null|string
     */
    protected $params = null;

    /**
     * 类型
     *
     * @var array|null|string
     */
    protected $type = null;

    /**
     * 类型代号
     *
     * @var array|null|string
     */
    protected $typeCode = null;
    /**
     * 方法
     *
     * @var array|null|string
     */
    protected $method = null;

    /**
     * 门店id
     *
     * @var null
     */
    protected $sourceId = null;

    /*******活动状态开始*********/
    /**
     * 草稿
     */
    const STATUS_DRAFT = 0;

    /**
     * 未上架
     */
    const STATUS_INIT = 1;

    /**
     * 已上架
     */
    const STATUS_UP = 2;

    /**
     * 已下架
     */
    const STATUS_DOWN = 3;

    /**
     * 已结束
     */
    const STATUS_OVER = 4;

    /**
     * 暂停
     */
    const STATUS_SUSPEND = 5;

    /**
     * 面向前端的假暂停状态
     */
    const STAUS_SHAM_SUSPEND = 6;

    const STATUS_RECOVER = 100;

    /**
     * 定时上架
     */
    const STATUS_CRON_UP = 6;

    /**
     * 定时下架
     */
    const STATUS_CRON_DOWN = 7;
    /*******活动状态结束*********/

    /*******活动类型开始*********/
    /**
     * 助力活动商品
     */
    const TYPE_HELP = 6;
    const TYPE_HELP_GIFT = 7;

    /**
     * 新人购活动商品
     */
    const TYPE_NEWCOMER = 11;

    /**
     * 限时购活动
     */
    const TYPE_FLASH_SALE = 12;

    /**
     * 助力抽奖活动
     */
    const TYPE_LUCK = 13;

    /**
     * 拼团
     */
    const TYPE_GROUP = 14;

    /**
     * 砍价
     */
    const TYPE_HAGGLE = 15;

    /**
     * 先享卡
     */
    const TYPE_DISCOUNT_CARD = 17;

    /**
     * 满赠
     */
    const TYPE_FULL_GIFT = 20;

    /**
     * 直购
     */
    const TYPE_DIRECT_BUY = 22;

    /*******活动类型结束*********/

    /**
     * 一口价
     */
    const TYPE_BUYOUT_PRICE = 21;
    /**
     * 构造方法
     *
     * ActivityBaseController constructor.
     */
    function __construct()
    {
        parent::__construct();
        $this->activityApi             = new Activity;
        $this->goodsActivityApi        = new GoodsActivity;
        $this->params                  = $this->request->input();
        $this->params['enterprise_id'] = $this->enterpriseId;
        $this->sourceId                = $this->enterpriseId;
        if (is_numeric($this->params['type'])) {  //兼容加库存类型为数字
            $this->type = $this->typeToCode($this->params['type'], true);
        } else {
            $this->type = $this->camelToSnake($this->params['type']);
        }
        $this->params['type'] = $this->type;
        $this->typeCode       = $this->typeToCode($this->type);

        $this->method = $this->params['method'] ?? '';
    }

    /**
     * 通知商品增加占用关系
     *
     * @param $goodsData
     */
    function noticeGoodsStart($goodsData)
    {
        try {
            (new GoodsActivity)->activityStart($goodsData);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 通知商品增加占用关系
     *
     * @param $goodsData
     */
    function noticeGoodsAddChannel($goodsData)
    {
        try {
            (new MaterialChannel())->edit($goodsData);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 通知商品解除占用关系
     *
     * @param $goodsData
     */
    function noticeGoodsEnd($goodsData)
    {
        try {
            (new GoodsActivity)->activityStop($goodsData);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 获取活动信息
     *
     * @param $activityId
     * @return mixed
     */
    function getActivityInfo($activityId)
    {

        $search['id']            = $activityId;
        $search['type']          = $this->type;
        $search['method']        = 'detail';
        $search['enterprise_id'] = $this->enterpriseId;

        try {
            $activityInfo = $this->activityApi->getActivity($search);
        } catch (\Exception $e) {
            $this->notice($e);
        }
        return $activityInfo;
    }

    /**
     * 组装通知商品增加占用关系的数据格式
     *
     * @param $params
     * @param $type
     * @param $ruleId
     * @return array
     */
    public function formatRelationData($params, $type, $ruleId)
    {
        if (!isset($params['activity_sku']) || empty($params['activity_sku'])) {
            return [];
        }

        $skuList = [];
        array_walk($params['activity_sku'], function ($v) use (&$skuList) {
            if ($v['activity_status'] == 1) {
                $skuList[$v['sku_id']][] = [
                    'id'          => $v['sku_id'],
                    'material_id' => $v['material_id'],
                    'price'       => $v['activity_price'],
                    'stock'       => $v['activity_stock_set'],
                ];
            }
        });

        array_walk($params['activity_sku'], function ($v) use ($skuList, &$result) {
            if ($v['activity_status'] == 1) {
                $result[] = [
                    'spu_id'   => $v['spu_id'],
                    'sku_list' => $skuList[$v['sku_id']]
                ];
            }
        });


        return [
            'enterprise_id'       => $params['enterprise_id'],
            'source_id'           => $params['enterprise_id'],
            'activity_id'         => $ruleId,
            'activity_type'       => $type,
            'activity_name'       => $params['name'],
            'activity_start_time' => (int)$params['start_at'],
            'activity_stop_time'  => (int)$params['end_at'],
            'goods'               => $result ?? []
        ];
    }

    /**
     * 通知商品增加渠道
     * @param $params
     * @param $type
     * @param $ruleId
     * @return array
     */
    protected function formatAddNewChannel($params, $type, $ruleId, $optionType = 'add')
    {
        if (!isset($params['activity_sku']) || empty($params['activity_sku'])) {
            return [];
        }

        $skuList = [];
        array_walk($params['activity_sku'], function ($v) use (&$skuList, $optionType) {
            // 如果活动商品状态开启有效
            //if ($v['activity_status'] == 1) {
                $skuList[] = [
                    'sku_id'                => $v['sku_id'],
                    'sale_type'             => 1,
                    'material_id'           => $v['material_id'],
                    'on_line_price'         => $v['activity_price'],
                    'on_line_current_stock' => $v['activity_stock_set'],
                    'type'                  => $optionType
                ];
            // }
        });

        return [
            'enterprise_id'       => $params['enterprise_id'],
            'source_id'           => $params['enterprise_id'],
            'object_id'           => $ruleId,
            'channel_type'        => $type,
            'object_name'         => $params['name'],
            'activity_start_time' => (int)$params['start_at'],
            'activity_stop_time'  => (int)$params['end_at'],
            'data'                => $skuList
        ];
    }

    /**
     * 自定义抛出错误 抛出服务的错误，方便排查
     *
     * @param \Exception $e
     */
    function notice(\Exception $e, $context = [], $msg = null, $data = [])
    {
        if ($context) {
            // 记录日志
            $msg = $msg ?: "request " . $this->type . " " . $this->method . " error";
            Log::error($msg, $context);
        }
        throw new Error($e->getCode(), $e->getMessage(), $data);
    }

    /**
     * 此处注意与服务的路由对应，newcomer可能有连字符
     *
     * @param $type
     * @return mixed
     */
    function typeToCode($type, $codeToStr = false)
    {

        $map = [
            'help'          => self::TYPE_HELP,
            'newcomer'      => self::TYPE_NEWCOMER,
            'flash-sale'    => self::TYPE_FLASH_SALE,
            'luck'          => self::TYPE_LUCK,
            'group'         => self::TYPE_GROUP,
            'haggle'        => self::TYPE_HAGGLE,
            'discount-card' => self::TYPE_DISCOUNT_CARD,
            'buyout-price'  => self::TYPE_BUYOUT_PRICE,
            'direct-buy'    => self::TYPE_DIRECT_BUY,
        ];

        if ($codeToStr) {
            $map = array_flip($map);
        }

        $type = strtolower($type);

        if (array_key_exists($type, $map)) {
            return $map[$type];
        }

    }

    /**
     * 把驼峰式的type类型转为中划线连接（路由使用）
     *
     * @param $camelType
     * @param string $separator
     * @return string
     */
    function camelToSnake($camelType, $separator = '-')
    {
        return strtolower(preg_replace('/([a-z])([A-Z])/', "$1" . $separator . "$2", $camelType));
    }

}
