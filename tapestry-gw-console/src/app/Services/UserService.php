<?php
/**
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/9/12
 * Time: 16:52
 */

namespace App\Services;

use App\Common\Constant;
use CrsSdk\Shop\User;
use CrsSdk\Work\Customer;

class UserService extends BaseService
{
    /**
     * 用户信息
     * @param array $userIds
     * @param int $size
     * @return array
     */
    public function byList($userIds = [], $size = Constant::DEFAULT_SIZE)
    {
        $resultData = [];

        // 用户列表
        $userData = (new User())->users(
            ['user_ids' => $userIds, 'size' => $size]
        );

        foreach ($userData['data'] as $key => $val) {
            // 昵称
            $name = $val['nick_name'];
            if (!empty($val['real_name'])) {
                $name .= ' (' . $val['real_name'] . ')';
            } else {
                // 备注
                if (!empty($val['union_id'])) {
                    $remarkName = (new Customer())->remark(['union_id' => $val['union_id']])['data']['user_remark_name'] ?? '';
                    $name .= empty($remarkName) ? '' : ' (' . $remarkName . ')';
                }
            }

            $resultData[$val['id']] = [
                'user_nick' => $name,
                'mobile' => $val['phone'] ?? '',
                'avatar' => $val['avatar'],
            ];
        }

        return $resultData;
    }
}
