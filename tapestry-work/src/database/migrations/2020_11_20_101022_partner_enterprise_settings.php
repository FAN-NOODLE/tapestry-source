<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class PartnerEnterpriseSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('partner_enterprise_settings')) {

            Schema::create('partner_enterprise_settings', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->unsignedBigInteger('opt_admin_id')->default(0)->comment('操作人ID');
                $table->string('opt_admin_name')->default('')->comment('操作人名称');
                $table->tinyInteger('partner_status')->default(0)->comment('1:开启二级分销 0:关闭二级分销');
                $table->tinyInteger('commission_type')->default(0)->comment('1:基于商品分销 2:基于推广佣金总额的分账');
                $table->unsignedInteger('commission_rate')->default(0)->comment('防止出现小数,设置1 = 10000; 千分之十：100');
                $table->tinyInteger('withdraw_status')->default(0)->comment('0:不可提现 1:可以提现');
                
                $table->index('enterprise_id', 'enterprise_id_index');
            });
            DB::statement("ALTER TABLE `partner_enterprise_settings` COMMENT '二级分销针对企业全局设置表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('partner_enterprise_settings');
    }
}
