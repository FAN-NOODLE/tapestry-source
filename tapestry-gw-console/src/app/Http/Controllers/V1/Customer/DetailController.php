<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 11:01
 */

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use CrsSdk\Shop\UserLabel;
use CrsSdk\Shop\User;

class DetailController extends BaseController
{
    /**
     *  服务类型： 商城
     */
    const SERVE_TYPE_SHOP = 1;
    /**
     *  服务类型： crs
     */
    const SERVE_TYPE_CRS = 2;

    /**
     * 客户基础信息
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 客户id
        $userId = $request->get('user_id', 0);
        // 外部联系人记录id
        $ecId = $request->get('ec_id', '');

        // 调用sdk
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'qy_domain' => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'user_id' => $userId,
            'ec_id' => $ecId
        ];

        if ($this->serveType == self::SERVE_TYPE_SHOP) {
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'id' => $userId
            ];
            $data = (new User())->detail($params);
            $data['data']['user_id'] = $data['data']['id'];
            unset($data['data']['id']);
        } else {
            $data = (new Customer())->detail($params);
        }

        // 标签
        $customer['user_id'] = array_column($data, 'user_id');
        $userDetail = (new \CrsSdk\Shop\Customer())->detailCustomer(['user_id' => $customer['user_id'], 'enterprise_id' => $this->enterpriseId])['data'];

        // 判断如果不是商城会员
        if ($userId == 0) {
            // 找到该客户的ec_userid
            foreach ($userDetail as $value) {
                if ($value['ec_id'] == $ecId) {
                    $customer['ec_userid'][] = $value['ec_userid'];
                }
            }
        } else {
            $customer['ec_userid'] = array_column($userDetail, 'ec_userid');
        }

        $customer['enterprise_id'] = $this->enterpriseId;
        $customer['serve_type'] = $this->serveType;
        $label = (new UserLabel())->search($customer)['data'];

        // 数据格式化
        $data['data'] = $this->format($data['data'], $label);

        return success($data);
    }

    /**
     * 格式化客户列表
     *
     * @param array $detail 待格式化客户信息
     * @param array $label 标签数组
     *
     * @return array
     */
    protected function format($detail, $label)
    {
        if (!$detail) {
            return [];
        }

        return [
            'ec_id' => $detail['ec_id'] ?? '',
            'user_id' => $detail['user_id'],
            'nick_name' => $detail['nick_name'],
            'real_name' => $detail['real_name'],
            'avatar' => $detail['avatar'],
            'sex' => $detail['sex'],
            'impression_label' => array_values(array_unique(array_column($label['ec_userid'], 'label_name'))),
            'system_label' => array_values(array_unique(array_column($label['user_id'], 'label_name'))),
            'birthday' => $detail['birthday'] ?? '',
            'mobile' => $detail['mobile'],
            'sales_username' => $detail['sales_username'] ?? '',
            'sales_dp_name' => $detail['sales_dp_name'] ?? '',
            'register_at' => $detail['register_at'],
            'level_num' => $detail['level_num'],
            'level_name' => $detail['level_name'],
            'level_is_highest' => $detail['level_is_highest'],
            'expiration_at' => $detail['expiration_at'],
            'amount' => $detail['amount']
        ];
    }
}
