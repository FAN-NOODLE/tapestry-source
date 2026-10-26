<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIndexToDailyMemberCustomerStatisticTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('daily_member_customer_statistic')) {
            DB::statement("ALTER TABLE daily_member_customer_statistic ADD INDEX `idx_enterprise_id_stat_time_mem_uid` (`enterprise_id`,`stat_time`,`mem_uid`,`cumulative_customer_num`,`today_customer_num`,`chat_cnt`,`reply_percentage`)");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('daily_member_customer_statistic')) {
            DB::statement("ALTER TABLE `daily_member_customer_statistic` DROP INDEX `idx_enterprise_id_stat_time_mem_uid`");
        }
    }
}
