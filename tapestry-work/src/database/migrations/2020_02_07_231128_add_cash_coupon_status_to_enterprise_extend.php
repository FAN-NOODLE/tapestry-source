<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddCashCouponStatusToEnterpriseExtend extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_extend')) {
            Schema::table('enterprise_extend', function (Blueprint $table) {
                $table->unsignedTinyInteger('cash_coupon_status')->default(2)->after('allocation_status')->comment('开启红包管理：1、开启，2、关闭');
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
        if (Schema::hasTable('enterprise_extend')) {
            Schema::table('enterprise_extend', function (Blueprint $table) {
                $table->dropColumn('cash_coupon_status');
            });
        }
    }
}
