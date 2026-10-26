<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddPartnerCommissionToDistributionOrderGoodsSettlementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('distribution_order_goods_settlement')) {
            if (!Schema::hasColumn('distribution_order_goods_settlement', 'partner_commission')) {
                Schema::table('distribution_order_goods_settlement',function (Blueprint $table) {
                    $table->unsignedInteger('partner_commission')->default('0')->comment('合伙人佣金');
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
        Schema::table('distribution_order_goods_settlement', function (Blueprint $table) {
            $table->dropColumn('partner_commission');
        });
    }
}
