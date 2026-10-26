<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseConfigTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_config')) {
            Schema::create('enterprise_config', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('type')->default(0)->comment('类别(1：好物圈设置，2：商品详情页设置，3：下单小程序授权，4：购物袋设置 注意:新增设置时开发人员需要维护该值)');
                $table->string('key',100)->default('')->comment('字段标示');
                $table->unsignedInteger('value')->default(0)->comment('1=开启 2=关闭');
                $table->string('desc',100)->default(0)->comment('字段说明');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_config` comment '企业设置表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_config');
    }
}
