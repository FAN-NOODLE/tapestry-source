<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddRefundStatusAndQuantityToDistributionOrderGoodsSettlementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('distribution_order_goods_settlement')) {
            if (!Schema::hasColumn('distribution_order_goods_settlement', 'refund_status')) {
                Schema::table('distribution_order_goods_settlement',function (Blueprint $table) {
                    $table->tinyInteger('refund_status')->default('0')->comment('refund_status 0:默认 1:部分退款 2:全部退款');
                    $table->bigInteger('refund_quantity')->default('0')->comment('refund_quantity 退款件数');
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
            $table->dropColumn('refund_status');
            $table->dropColumn('refund_quantity');
        });
    }
}
