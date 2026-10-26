<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/14
 * Time: 8:01 PM
 */

namespace App\Console\Commands;

use CrsSdk\Sdk;
use Illuminate\Console\Command;

class Test extends Command
{


    /**
     * @var string
     */
    protected $signature = 'doc';

    /**
     * @var string
     */

    public function handle()
    {
        $params['enterprise_id'] = 3363603667487424;

        $res = Sdk::work()->setApi("service/v1/data-compass/customer-statistic/fiscal-month")->post($params);
        dd($res);

    }
}
