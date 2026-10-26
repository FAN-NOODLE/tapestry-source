<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterprisePaymentTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_payment')) {
            Schema::create('enterprise_payment', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('name')->default('')->comment('活动名称');
                $table->string('number', '50')->default('')->comment('活动编号');
                $table->unsignedTinyInteger('type')->default(1)->comment('红包类型：1=固定金额, 2=2随机金额');
                $table->unsignedInteger('min_amount')->default(0)->comment('最小红包区间金额');
                $table->unsignedInteger('max_amount')->default(0)->comment('最大红包区间金额');
                $table->unsignedInteger('fixed_amount')->default(0)->comment('固定红包金额');
                $table->unsignedInteger('enterprise_payment_amount')->default(0)->comment('红包总金额');
                $table->unsignedInteger('quantity')->default(0)->comment('红包个数');
                $table->unsignedTinyInteger('is_old_user_receive')->default(1)->comment('企业现有客户是否可以领取：1=否, 2=是');
                $table->string('qr_code_url')->default('')->comment('员工二维码图');
                $table->unsignedTinyInteger('is_sale')->default(1)->comment('上架状态：1=未上架, 2=上架, 3=已下架,4=已结束');
                $table->unsignedTinyInteger('is_deleted')->default(1)->comment('删除状态 1 未删除 2 已删除');
                $table->string('send_name')->default('')->comment('商户名称');
                $table->string('wishing')->default('')->comment('红包祝福语');
                $table->string('remark')->default('')->comment('红包备注');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `enterprise_payment` comment '企业微信红包表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_payment');
    }
}
