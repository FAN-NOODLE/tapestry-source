<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-18
 * Time: 10:06
 */

namespace App\Console\Commands;

use App\Common\Constant;
use CrsSdk\Work\Customer;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class ExportCustomer extends Command
{
    /**
     * @var string
     */
    protected $signature = 'export:customer-list';

    /**
     * @var string
     */
    protected $description = '导出顾客数据';

    protected $headers
        = [
            'mem_name' => '',
            'name'     => '',
            'phone'    => '',
            'gender'   => '',
            'dep_name' => ''
        ];

    public function handle()
    {
        // 写入头
        $customer = new Customer();
        $params   = [
            'enterprise_id'  => '100987608609664',
            'qy_domain'      => 'C6B0D5057F00000155E244CC883F01C3',
            'app_identifier' => 'crs-shop',
            'start_time'     => '1597420800',
            'end_time'       => time(),
            'page'           => 1,
            'size'           => Constant::MAX_DEFAULT_SIZE
        ];

        while (true) {
            $data = [];
            $list = $customer->customerList($params);
            if (empty($list['list']['data'])) {
                break;
            }
            foreach ($list['list']['data'] as $item) {
                $tmp = [];
                foreach ($this->headers as $key => $value) {
                    $tmp[$key] = $item[$key];
                }
                $data[] = $tmp;
            }
            DB::table('customer_cache_hxm')->insert($data);
            $params['page'] += 1;
        }
    }
}
