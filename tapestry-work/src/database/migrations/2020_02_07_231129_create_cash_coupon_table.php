<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCashCouponTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('cash_coupon')) {
            Schema::create('cash_coupon', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('cash_coupon_name')->default(0)->comment('活动名称');
                $table->string('number', '50')->default('')->comment('活动编号');
                $table->unsignedInteger('amount')->default(0)->comment('红包金额')->index();
                $table->unsignedInteger('quantity')->default(0)->comment('红包个数')->index();
                $table->unsignedTinyInteger('is_sale')->default(1)->comment('上架状态：1=未上架, 2=上架, 3=已下架');
                $table->unsignedTinyInteger('is_deleted')->default(1)->comment('删除状态 1 未删除 2 已删除');
                $table->string('send_name')->default(0)->comment('商户名称');
                $table->string('wishing')->default(0)->comment('红包祝福语');
                $table->string('remark')->default(0)->comment('红包备注');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            // 表注释
            DB::statement("ALTER TABLE `cash_coupon` comment '现金红包活动表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('cash_coupon');
    }
}
