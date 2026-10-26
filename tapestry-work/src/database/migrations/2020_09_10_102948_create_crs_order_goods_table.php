<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCrsOrderGoodsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('crs_order_goods')) {
            Schema::create('crs_order_goods', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
                $table->string('order_no', 30)->default('')->comment('订单号')->index();
                $table->string('sku_no', 50)->default('')->comment('sku编号')->index();
                $table->unsignedBigInteger('spu_id')->default(0)->comment('SPU ID');
                $table->unsignedBigInteger('sku_id')->default(0)->comment('SKU ID');
                $table->string('spu_name', 191)->default('')->comment('SPU商品名称');
                $table->unsignedInteger('type')->default(0)->comment('订单商品类型：1普通商品，2拼团，3秒杀，5超低购，6助力，7砍价，11新人购商品，12限时购商品');
                $table->unsignedInteger('goods_type')->default(0)->comment('商品类型：1普通，2赠品，3奖品');
                $table->string('category_name', 191)->default('')->comment('分类名称');
                $table->string('sku_picture', 191)->default('')->comment('SKU主图URL');
                $table->text('sku_attribute')->nullable()->comment('SKU商品属性值集合');
                $table->unsignedInteger('quantity')->default(0)->comment('购买数量');
                $table->unsignedInteger('original_price')->default(0)->comment('商品原价，单位：分');
                $table->unsignedInteger('sell_price')->default(0)->comment('销售单价，单位：分');
                $table->unsignedInteger('single_discount')->default(0)->comment('单品优惠，单位：分');
                $table->unsignedInteger('sub_price')->default(0)->comment('总价，单位：分');
                $table->unsignedInteger('goods_pay_price')->default(0)->comment('订单商品分摊实付金额');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `crs_order_goods` comment '订单商品表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('crs_order_goods');
    }
}
