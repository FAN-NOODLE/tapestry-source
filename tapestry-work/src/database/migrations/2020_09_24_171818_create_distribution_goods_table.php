<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateDistributionGoodsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('distribution_goods')) {
            Schema::create('distribution_goods', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('sku_id')->default(0)->comment('skuid');
                $table->unsignedBigInteger('spu_id')->default(0)->comment('spuid');
                $table->unsignedTinyInteger('promotion_commission_type')->default(1)->comment('推广佣金类型：1=百分比2=金额');
                $table->unsignedInteger('promotion_commission_amount')->default(0)->comment('推广佣金设置/百分比整数/元两位小数');
                $table->unsignedTinyInteger('service_commission_type')->default(1)->comment('服务佣金类型：1=百分比2=金额');
                $table->unsignedInteger('service_commission_amount')->default(0)->comment('服务佣金设置/百分比整数/元两位小数');
                $table->unsignedInteger('distribution_start')->default(0)->comment('分销开始时间');
                $table->unsignedInteger('distribution_end')->default(0)->comment('分销结束时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `distribution_goods` comment '分销数据表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
        Schema::dropIfExists('distribution_goods');
    }
}
