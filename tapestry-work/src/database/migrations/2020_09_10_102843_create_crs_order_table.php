<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCrsOrderTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('crs_order')) {
            Schema::create('crs_order', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedTinyInteger('sync_source')->default(1)->comment('订单同步来源：1-crs, 2-hwj');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('user_id')->default(0)->comment('会员ID');
                $table->string('union_id')->default('')->comment('用户union_id')->index();
                $table->string('open_id')->default('')->comment('用户open_id');
                $table->string('order_no', 30)->default('')->comment('订单号');
                $table->unsignedTinyInteger('status')->default(1)->comment('订单状态：1 待付款，2 已完成，3 已取消, 4 订单确认中');
                $table->unsignedTinyInteger('order_type')->default(1)->comment('订单类型: 1 线上订单，2 门店订单');
                $table->unsignedTinyInteger('payment_channel')->default(1)->comment('支付渠道：1 微信小程序支付，2 美团支付');
                $table->unsignedInteger('total_amount')->default(0)->comment('原价 单位：分');
                $table->unsignedInteger('preferential_price')->default(0)->comment('优惠金额，单位：分');
                $table->unsignedInteger('final_amount')->default(0)->comment('实际付款价格，单位：分');
                $table->unsignedInteger('pay_at')->default(0)->comment('支付时间');
                $table->unsignedInteger('refund_at')->default(0)->comment('退款时间');
                $table->unsignedInteger('order_at')->default(0)->comment('下单时间');
                $table->string('delivery_name', 50)->default('')->comment('收货人姓名');
                $table->string('delivery_mobile', 20)->default('')->comment('收货人手机号码');
                $table->string('delivery_address', 200)->default('')->comment('收货人地址');
                $table->string('remarks', 200)->default('')->comment('卖家备注');
                $table->string('province', 30)->default('')->comment('省');
                $table->string('city', 30)->default('')->comment('市');
                $table->string('area', 30)->default('')->comment('区');
                $table->string('logistic_identity', 200)->default('')->comment('物流公司');
                $table->string('logistic_no', 200)->default('')->comment('物流编号');
                $table->unsignedInteger('send_at')->default(0)->comment('发货时间');
                $table->unsignedInteger('sign_at')->default(0)->comment('签收时间');
                $table->unsignedInteger('cash_fee')->default(0)->comment('实际支付金额，单位：分');
                $table->unsignedInteger('coupon_fee')->default(0)->comment('微信优惠，单位：分');
                $table->unsignedInteger('coupon_count')->default(0)->comment('微信优惠券数量');
                $table->unsignedInteger('other_fee')->default(0)->comment('其他优惠，单位：分');
                $table->unsignedInteger('read_status')->default(1)->comment('已读状态  1 未读 2 已读');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->unique('order_no', 'order_no_index');
            });

            DB::statement("alter table `crs_order` comment '订单表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('crs_order');
    }
}
