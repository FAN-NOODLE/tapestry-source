<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/1/24
 * Time: 21:31
 */

namespace App\Services;

use CrsSdk\Shop\UserLabel;

class CustomerService extends BaseService
{
    /**
     * 获取下单会员
     *
     * @param array $orderCustomer 下单会员id
     * @param int $enterpriseId 企业id
     * @param int $serveType 企业服务类型:默认crs+shop
     *
     * @return array
     */
    public function getLabel($orderCustomer, $enterpriseId, $serveType = 0)
    {
        // 标签
        $customer['ec_userid'] = array_filter(array_column($orderCustomer, 'ec_userid'));
        $customer['user_id'] = array_filter(array_column($orderCustomer, 'user_id'));
        $customer['enterprise_id'] = $enterpriseId;
        $customer['serve_type'] = $serveType;
        $label = (new UserLabel())->search($customer);

        // 获取主值
        $ecUserId = array_column($label['data']['ec_userid'], 'ec_userid');
        $userId = array_column($label['data']['user_id'], 'user_id');
        // 处理数组
        foreach ($label['data'] as $key => &$value) {
            if ($key == 'ec_userid') {
                $ecUserIdLabel = [];
                foreach ($value as $ecUserIdLabelValue) {
                    $ecUserIdLabel[$ecUserIdLabelValue['ec_userid']][] = $ecUserIdLabelValue;
                }
                $value = $ecUserIdLabel;
            }
            if ($key == 'user_id') {
                $userIdLabel = [];
                foreach ($value as $userIdLabelValue) {
                    $userIdLabel[$userIdLabelValue['user_id']][] = $userIdLabelValue;
                }
                $value = $userIdLabel;
            }
        }

        return [$ecUserId, $userId, $label];
    }

}
