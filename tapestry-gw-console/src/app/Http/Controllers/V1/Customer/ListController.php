<?php

namespace App\Http\Controllers\V1\Customer;

use CrsSdk\Shop\User;
use App\Common\Constant;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 客户列表
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', Constant::DEFAULT_SIZE);
        // 查询条件
        $search = $request->get('search', []);

        if (empty($this->qyDomain)) {
            return success([
                'data' => [],
                'meta' => [
                    'total' => 0,
                    'page'  => intval($page),
                    'size'  => intval($size)
                ]
            ]);
        }

        // 调用sdk
        $params = [
            'enterprise_id'  => $this->enterpriseId,
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'page'           => $page,
            'size'           => $size,
            'search'         => $search
        ];
        $data = ( new Customer() )->customers($params);

        // 数据格式化
        $data['data'] = $this->format($data['data']);

        return success($data);
    }

    /**
     * 新客户列表
     * @param Request $request
     * @return JsonResponse
     */
    public function customerList(Request $request)
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new \CrsSdk\Shop\Customer())->customerList($params);

        return success($data);
    }

    /**
     * 客户购买记录列表
     * @param Request $request
     * @return JsonResponse
     */
    public function buyRecord(Request $request)
    {
        // 分页参数
        $id = $request->get('id');
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', Constant::DEFAULT_SIZE);
        // 调用sdk
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'user_id' => $id,
            'page' => $page,
            'size' => $size,
        ];
        $result = (new User())->userOrderActionRecord($params);
        if (empty($result['list'])){
            return success([]);
        }
        foreach ($result['list'] as &$item) {
            $item['type'] = $item['action_type'];
            $item['create_time'] = isset($item['created_at']) ? date('Y-m-d H:i:s', $item['created_at']) : '';
            $item['order_no'] = $item['main_order_no'];
            unset($item['main_order_no'],$item['order_id'],$item['user_id'],$item['action_type'],$item['created_at']);
        }
        return success($result);
    }

    /**
     * 格式化客户列表
     * @param $list
     * @return array
     */
    protected function format($list)
    {
        // 初始化待返数据
        $data = [];
        if (!$list) {
            return $data;
        }

        foreach ($list as $customer) {
            $data[] = [
                'ec_id'             => $customer['ec_id'],
                'user_id'           => $customer['user_id'],
                'nick_name'         => $customer['nick_name'],
                'real_name'         => $customer['real_name'],
                'avatar'            => $customer['avatar'],
                'level_name'        => $customer['level_name'],
                'birthday'          => $customer['birthday'],
                'mobile'            => $customer['mobile'],
                'sales_username'    => $customer['sales_username'],
                'sales_dp_name'     => $customer['sales_dp_name']
            ];
        }
        return $data;
    }
}