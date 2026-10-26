<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/14
 * Time: 10:24 AM
 */

namespace App\Cache;

use App\Models\Enterprise;

class EnterpriseCache extends Base
{
    /**
     * 获取企业详情
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/07 11:16
     * @param string $key 查询字段key
     * @param string $value 查询字段value
     * @return mixed
     */
    public function detail($key, $value)
    {
        return $this->call(function () use ($key, $value) {
            if ($enterprise = Enterprise::query()->where($key, $value)->first()) {
                return $enterprise;
            } else {
                return false;
            }
        });
    }
}
