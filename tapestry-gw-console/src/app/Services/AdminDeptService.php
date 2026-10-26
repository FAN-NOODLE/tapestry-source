<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/28 0028
 * Time: 10:13
 */

namespace App\Services;

use App\Micros\Work\Work;
use CrsSdk\MicroAuthority\Account;

class AdminDeptService
{


    /**
     * @return array
     */
    public function list($enterpriseId, $mobile): array
    {
        $adminDpIds = [];

        $input = [
            'mobile'        => str_replace("****","%", $mobile),
            'enterprise_id' => $enterpriseId,
            'only'          => 2
        ];
        $params = [
            'admin_id' => (new Account())->detailAccount($input)['id']
        ];

        $res = (new Account())->dataRange($params);

        foreach ($res['data_range_list'] as $object) {
            if ($object['type'] == 4) {
                $adminDpIds[] = $object['object_id'];
            }
        }
        if ($adminDpIds) {
            $departmentLists = Work::lists(['enterprise_id' => $enterpriseId, 'params' => ['third_ids' => $adminDpIds]]);
        }

        return [
            'range_type' => 4,
            'list'       => $departmentLists ?? []
        ];
    }


}