<?php
/**
 * Created by PhpStorm.
 * User: shizhixiong
 * Date: 2019-01-24
 * Time: 21:30
 */

namespace App\Services;

use CrsSdk\Work\Member;

class MemberService extends BaseService
{
    /**
     * 获取员工详情
     *
     * @param       $domain
     * @param       $appIdentifier
     * @param array $memUids
     * @param int   $pageSize 默认1000 超过需遍历
     *
     * @return array
     */
    public function getFirstMemberData($domain, $appIdentifier, $memUids= [], $pageSize = 1000)
    {
        if (empty($domain) || empty($appIdentifier) || empty($memUids)) {
            return [];
        }

        $data = (new Member())->memberList([
            'params' => ['memUids' => $memUids, 'pageSize' => $pageSize],
            'qy_domain' => $domain,
            'app_identifier' => $appIdentifier,
        ])['list'];

        $returnData = [];

        array_walk($data, function ($value) use (&$returnData) {
            $dpName = '';
            if (isset($value['dpName']) && current($value['dpName'])) {
                $dpName = current($value['dpName'])['dpName'];
            }
            $returnData[$value['memUid']] = [
                'memUid' => $value['memUid'],
                'memUsername' => $value['memUsername'],
                'memRole' => $value['memRole'],
                'dpName' => $dpName
            ];
        });

        return $returnData;
    }
}
