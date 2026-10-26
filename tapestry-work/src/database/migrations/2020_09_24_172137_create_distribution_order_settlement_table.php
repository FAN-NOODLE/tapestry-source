<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateDistributionOrderSettlementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('distribution_order_settlement')) {
            Schema::create('distribution_order_settlement', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('order_no',60)->default('')->comment('订单编号');
                $table->unsignedInteger('order_status')->default(0)->comment('订单状态订单状态：20:未支付;21:已取消;23:待支付;25:支付中;30:待确认(已支付);40:待发货;41:售后中(待发货,仅退款);60:待收货;61:售后中;80:已签收;81:售后中(已签收,退货退款);83:部分退款;100:已完成;200:退款中;201:退款失败;205:全部退款');
                $table->unsignedTinyInteger('settlement_status')->default(0)->comment('结算状态：0待结算1已结算2已取消');
                $table->unsignedInteger('settlement_standard')->default(100)->comment('结算标准：100已完成');
                $table->unsignedInteger('order_pay_time')->default(0)->comment('订单支付时间');
                $table->unsignedBigInteger('promotion_worker_id')->default(0)->comment('推广导购id');
                $table->unsignedInteger('promotion_commission')->default(0)->comment('推广佣金');
                $table->unsignedBigInteger('service_shopping_guid_id')->default(0)->comment('服务导购id');
                $table->unsignedInteger('service_commission')->default(0)->comment('服务佣金');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `distribution_order_settlement` comment '订单结算表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('distribution_order_settlement');
    }
}
