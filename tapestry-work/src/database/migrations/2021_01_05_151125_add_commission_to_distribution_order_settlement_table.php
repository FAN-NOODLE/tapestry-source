<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddCommissionToDistributionOrderSettlementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('distribution_order_settlement')) {
            if (!Schema::hasColumn('distribution_order_settlement', 'partner_commission_snapshot')) {
                Schema::table('distribution_order_settlement', function (Blueprint $table) {
                    $table->unsignedInteger('partner_commission_snapshot')->default('0')->comment('合伙人佣金快照已取消统计使用');
                    $table->unsignedInteger('promotion_commission_snapshot')->default('0')->comment('导购推广佣金快照已取消统计使用');
                    $table->unsignedInteger('service_commission_snapshot')->default('0')->comment('导购服务佣金快照已取消统计使用');
                });
            }
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('distribution_order_settlement', function (Blueprint $table) {
            $table->dropColumn('partner_commission_snapshot');
            $table->dropColumn('promotion_commission_snapshot');
            $table->dropColumn('service_commission_snapshot');
        });
    }
}
