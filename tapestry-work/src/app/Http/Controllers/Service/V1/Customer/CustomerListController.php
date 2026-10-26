<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 11:00
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerService;
use Illuminate\Http\Request;

class CustomerListController extends BaseController {
    /**
     * index
     * @desc 顾客列表筛选接口
     * @param Request $request
     *
     * @return array
     * @date 2020-07-17
     * @author XuHuitao
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index() {
        $param = $this->getParam();
        switchAppHost(['enterprise_id' => $param['enterprise_id']]);
        $result = CustomerService::getCustomerList($param);
        $result['page'] = $param['page'];
        $result['size'] = $param['size'];
        return success($result);
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'enterprise_id'     => 'integer|required',
            'store_id'          => 'integer',
            'mobile'            => 'required',
            'qy_domain'         => 'required',
            'page'              => 'integer',
            'size'              => 'integer',
            'gender'            => 'integer',
            'order_num_min'     => 'integer',
            'order_num_max'     => 'integer',
            'order_amount_min'  => 'integer',
            'order_amount_max'  => 'integer'
        ], [
            'enterprise_id.*'             => '企业id参数不合法',
            'store_id.*'                  => '门店id参数不合法',
            'mobile.*'                    => '用户手机号参数不合法',
            'page.integer'                => '分页不合法',
            'size.integer'                => '分页数量不合法',
            'gender.integer'              => '性别参数不合法',
            'order_num_min.integer'       => '订单数量参数不合法',
            'order_num_max.integer'       => '订单数量参数不合法',
            'order_amount_min.integer'    => '订单金额参数不合法',
            'order_amount_max.integer'    => '订单金额参数不合法',
        ]);
        $param = $this->request->only(['page', 'size', 'gender', 'order_num_min', 'order_num_max',
                                       'order_amount_min', 'order_amount_max', 'tag_ids', 'enterprise_id',
                                       'join_time_start', 'join_time_end', 'store_id', 'mobile', 'near_browse_start', 'near_browse_end', 'achievement_ascription', 'guide_tag_ids', 'buy_time_start', 'buy_time_end', 'add_type']);
        $param['achievement_ascription'] = $param['achievement_ascription'] ?? 0;
        $param['add_type'] = $param['add_type'] ?? 0;

        if (isset($param['order_num_min']) || isset($param['order_num_max'])) {
            $param['order_num_min'] = $param['order_num_min'] ?? Constant::CUSTOMER_ORDER_NUM_MIN;
            $param['order_num_max'] = $param['order_num_max'] ?? Constant::CUSTOMER_ORDER_NUM_MAX;

            if ($param['order_num_min'] > $param['order_num_max']) {
                list($param['order_num_max'], $param['order_num_min']) = [$param['order_num_min'], $param['order_num_max']];
            }
        }

        if (isset($param['order_amount_min']) || isset($param['order_amount_max'])) {
            $param['order_amount_min'] = $param['order_amount_min'] ?? Constant::CUSTOMER_ORDER_AMOUNT_MIN;
            $param['order_amount_max'] = $param['order_amount_max'] ?? Constant::CUSTOMER_ORDER_AMOUNT_MAX;

            if ($param['order_amount_min'] > $param['order_amount_max']) {
                list($param['order_amount_max'], $param['order_amount_min']) = [$param['order_amount_min'], $param['order_amount_max']];
            }
        }

        if (!isset($param['page'])) {
            $param['page'] = Constant::DEFAULT_PAGE;
        }

        if (!isset($param['size'])) {
            $param['size'] = Constant::DEFAULT_SIZE;
        }
        if (isset($param['tag_ids'])) {
            $param['tag_ids'] = explode(',', $param['tag_ids']);
        }
        if (isset($param['guide_tag_ids'])) {
            $param['guide_tag_ids'] = explode(',', $param['guide_tag_ids']);
        }
        if (isset($param['buy_time_start'])) {
            if (isset($param['buy_time_end'])) {
                if ($param['buy_time_start'] > $param['buy_time_end']) {
                    list($param['buy_time_start'], $param['buy_time_end']) = [$param['buy_time_end'], $param['buy_time_start']];
                }
            } else {
                // 未设置结束时间则查询开始时间当天数据
                $param['buy_time_end'] = $param['buy_time_start'] + 86399;
            }
        } else if (isset($param['buy_time_end'])) {
            // 未设置开始时间则查询结束时间当天数据
            $param['buy_time_start'] = $param['buy_time_end'] - 86399;
        }
        if (isset($param['join_time_start'])) {
            if (isset($param['join_time_end'])) {
                if ($param['join_time_start'] > $param['join_time_end']) {
                    list($param['join_time_start'], $param['join_time_end']) = [$param['join_time_end'], $param['join_time_start']];
                }
            } else {
                // 未设置结束时间则查询开始时间当天数据
                $param['join_time_end'] = $param['join_time_start'] + 86399;
            }
        } else if (isset($param['join_time_end'])) {
            // 未设置开始时间则查询结束时间当天数据
            $param['join_time_start'] = $param['join_time_end'] - 86399;
        }

        if (isset($param['near_browse_start'])) {
            if (isset($param['near_browse_end'])) {
                if ($param['near_browse_start'] > $param['near_browse_end']) {
                    list($param['near_browse_start'], $param['near_browse_end']) = [$param['near_browse_end'], $param['near_browse_start']];
                }
            } else {
                // 未设置结束时间则查询开始时间当天数据
                $param['near_browse_end'] = $param['near_browse_start'] + 86399;
            }
        } else if (isset($param['near_browse_end'])) {
            // 未设置开始时间则查询结束时间当天数据
            $param['near_browse_start'] = $param['near_browse_end'] - 86399;
        }

        return $param;
    }

    /**
     * search
     * @desc 搜索顾客
     * @return array
     * @date 2020-07-21
     * @throws \Illuminate\Validation\ValidationException
     */
    public function search()
    {
        $this->validate($this->request, [
            'enterprise_id'     => 'integer|required',
            'store_id'          => 'integer',
            'mobile'            => 'required',
            'qy_domain'         => 'required',
            'search_key'        => 'required',
            'page'              => 'integer',
            'size'              => 'integer',
        ], [
            'enterprise_id.*'               => '企业id参数不合法',
            'store_id.*'                    => '门店id参数不合法',
            'mobile.*'                      => '用户手机号参数不合法',
            'qy_domain.required'            => 'qy_domain参数缺失',
            'search_key.*'                  => '搜索关键字不能为空',
            'page.integer'                  => '分页不合法',
            'size.integer'                  => '分页数量不合法',
        ]);
        $param = $this->request->only(['page', 'size', 'enterprise_id', 'store_id', 'mobile', 'search_key', 'qy_domain']);
        $param['page'] = $param['page'] ?? Constant::DEFAULT_PAGE;
        $param['size'] = $param['size'] ?? Constant::DEFAULT_SIZE;
       //switchAppHost(['enterprise_id' => $param['enterprise_id']]);
        $customerList = CustomerService::searchCustomer($param['enterprise_id'], $param['search_key']);
        $param['search_union_id'] = array_keys($customerList['shop_customer']);
        $param['search_relation_id'] = array_keys($customerList['relations']);
        $param['add_type'] = 0;

        $result = CustomerService::getCustomerList($param);
        $result['page'] = $param['page'];
        $result['size'] = $param['size'];
        return success($result);
    }


}