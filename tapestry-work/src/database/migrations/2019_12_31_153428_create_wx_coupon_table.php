<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateWxCouponTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('wx_coupon')) {
            Schema::create('wx_coupon', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedBigInteger('shop_wx_coupon_id')->default(0)->comment('商城微信优惠券ID')->index();
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('coupon_name')->default('')->comment('优惠券名称');
                $table->unsignedInteger('coupon_value')->default(0)->comment('优惠券面额（单位分）');
                $table->unsignedInteger('coupon_minimum')->default(0)->comment('优惠券使用最低额（单位分）');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `wx_coupon` comment '微信优惠券表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('wx_coupon');
    }
}
