<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateUserCouponTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('user_coupon')) {
            Schema::create('user_coupon', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedBigInteger('shop_user_coupon_id')->default(0)->comment('商城用户优惠券ID')->index();
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('user_id')->default(0)->comment('用户id')->index();
                $table->unsignedTinyInteger('effective')->default(1)->comment('有效时间类型（1:无限制，2:固定时间段）');
                $table->unsignedInteger('start_time')->default(0)->comment('有效期开始时间(时间戳-秒)');
                $table->unsignedInteger('end_time')->default(0)->comment('有效期结束时间(时间戳-秒)');
                $table->unsignedTinyInteger('status')->default(1)->comment('使用状态（1:待使用，2:使用中，3:已使用）');
                $table->unsignedBigInteger('wx_coupon_id')->default(0)->comment('微信发放代金券id');
                $table->unsignedInteger('receive_time')->default(0)->comment('领取时间');
                $table->unsignedInteger('use_time')->default(0)->comment('使用时间');
                $table->string('order_no')->default('')->comment('订单编号');
                $table->unsignedInteger('order_amount')->default(0)->comment('订单总额');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `user_coupon` comment '用户优惠券表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('user_coupon');
    }
}
