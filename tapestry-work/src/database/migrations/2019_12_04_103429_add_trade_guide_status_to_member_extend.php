<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddTradeGuideStatusToMemberExtend extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('member_extend')) {
            Schema::table('member_extend', function (Blueprint $table) {
                $table->unsignedTinyInteger('trade_guide_status')->default(2)->after('allocation_status')->comment('支付即导购分配状态：1、开启，2、关闭');
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
        if (Schema::hasTable('member_extend')) {
            Schema::table('member_extend', function (Blueprint $table) {
                $table->dropColumn('trade_guide_status');
            });
        }
    }
}
