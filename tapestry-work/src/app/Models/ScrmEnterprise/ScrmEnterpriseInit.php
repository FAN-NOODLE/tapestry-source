<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models\ScrmEnterprise;

use App\Models\Base;
use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\DB;

class ScrmEnterpriseInit extends Base
{
    /**
     * 连接
     */
    protected $connection = 'micro_scrm_enterprise';

    /**
     * 连接这一部分数据库时要动态切换链接，因为数据库是根据企业id分库的
     *
     * @param $query
     * @param $enterpriseId
     * @return mixed
     * User: zhaoyang@vchangyi.com
     * Date: 2021/4/25 15:48
     */
    public function scopeSetConn($query, $enterpriseId){
        if (Config::get('database.connections.micro_scrm_enterprise.database') != 'micro_scrm_enterprise_' . $enterpriseId) {
            Config::set('database.connections.micro_scrm_enterprise.database', 'micro_scrm_enterprise_' . $enterpriseId);
            DB::connection('micro_scrm_enterprise')->reconnect();
        }
        return $query;
    }
}

