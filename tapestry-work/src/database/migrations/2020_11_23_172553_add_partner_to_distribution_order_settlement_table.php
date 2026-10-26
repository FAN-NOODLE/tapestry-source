<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddPartnerToDistributionOrderSettlementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('distribution_order_settlement')) {
            if (!Schema::hasColumn('distribution_order_settlement', 'partner_id')) {
                Schema::table('distribution_order_settlement',function (Blueprint $table) {
                    $table->unsignedBigInteger('partner_id')->default('0')->comment('合伙人ID');
                    $table->unsignedBigInteger('partner_enterprise_id')->default('0')->comment('合伙人在具体企业的ID');
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
        Schema::table('distribution_order_settlement', function (Blueprint $table) {
            $table->dropColumn('partner_id');
            $table->dropColumn('partner_enterprise_id');
            $table->dropColumn('partner_commission');
        });
    }
}
