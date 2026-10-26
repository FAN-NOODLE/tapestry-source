<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateDistributionOrderGoodsSettlementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('distribution_order_goods_settlement')) {
            Schema::create('distribution_order_goods_settlement', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('order_no',60)->default('')->comment('订单编号');
                $table->unsignedBigInteger('sku_id')->default(0)->comment('商品sku_id');
                $table->unsignedInteger('purchase_quantity')->comment('商品数量');
                $table->unsignedInteger('order_goods_price')->default(0)->comment('订单商品售卖价');
                $table->unsignedInteger('promotion_commission')->default(0)->comment('推广佣金');
                $table->unsignedInteger('service_commission')->default(0)->comment('服务佣金');
                $table->text('set_type_ratio_json')->comment('推广佣金及服务佣金设置方式及比例');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `distribution_order_goods_settlement` comment '订单商品结算表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('distribution_order_goods_settlement');
    }
}
