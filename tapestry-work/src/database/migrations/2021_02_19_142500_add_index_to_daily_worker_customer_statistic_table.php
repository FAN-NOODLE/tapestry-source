<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIndexToDailyWorkerCustomerStatisticTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('daily_worker_customer_statistic')) {
            DB::statement("ALTER TABLE daily_worker_customer_statistic ADD INDEX `idx_enterprise_id_stat_time_worker_id` (`enterprise_id`,`stat_time`,`worker_id`,`cumulative_contact_card_way_customer_num`,`cumulative_other_way_customer_num`,`cumulative_binding_customer_num`,`today_contact_card_way_customer_num`,`today_other_way_customer_num`,`today_binding_customer_num`)");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('daily_worker_customer_statistic')) {
            DB::statement("ALTER TABLE `daily_worker_customer_statistic` DROP INDEX `idx_enterprise_id_stat_time_worker_id`");
        }

    }
}
