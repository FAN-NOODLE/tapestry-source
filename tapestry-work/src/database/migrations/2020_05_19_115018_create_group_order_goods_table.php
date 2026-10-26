<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateGroupOrderGoodsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('group_order_goods')) {
            Schema::create('group_order_goods', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('qw_group_id', 32)->default('')->comment('企业微信群ID');
                $table->string('group_user_id', 50)->default('')->comment('用户ID');
                $table->string('union_id', 50)->default('')->comment('用户唯一ID');
                $table->string('order_no', 191)->default('')->comment('订单号');
                $table->unsignedBigInteger('spu_id')->default(0)->comment('SPU ID');
                $table->unsignedBigInteger('sku_id')->default(0)->comment('SKU ID');
                $table->string('spu_name', 191)->default('')->comment('SPU商品名称');
                $table->string('category_name', 191)->default('')->comment('分类名称');
                $table->string('sku_picture', 191)->default('')->comment('SKU主图URL');
                $table->string('spu_picture', 191)->default('')->comment('SPU主图URL');
                $table->unsignedInteger('quantity')->default(0)->comment('购买数量');
                $table->unsignedInteger('refund_quantity')->default(0)->comment('退货数量');
                $table->unsignedInteger('sell_price')->default(0)->comment('销售单价，单位：分');
                $table->unsignedInteger('sub_price')->default(0)->comment('总价，单位：分');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `group_order_goods` comment '企业微信群订单商品表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('group_order_goods');
    }
}
