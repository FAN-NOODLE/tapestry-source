<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddLatestBuyTimeInMemberOrderStatistics extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('member_order_statistics')) {
            Schema::table('member_order_statistics', function (Blueprint $table) {
                if (!Schema::hasColumn('member_order_statistics', 'latest_buy_time')) {
                    $table->unsignedInteger('latest_buy_time')->default(0)->comment('最近下单时间')->after('total_order_amount');
                }
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
        Schema::table('member_order_statistics', function ($table) {

        });
    }
}
