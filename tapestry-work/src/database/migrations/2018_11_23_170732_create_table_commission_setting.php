<?php

use \Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableCommissionSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('commission_setting')) {
            Schema::create('commission_setting', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->unsignedTinyInteger('type')->default(1)->comment('导购类型（1：专属导购，2：普通导购，3: 服务导购）');
                $table->unsignedInteger('value')->default(0)->comment('佣金比例，百分比乘100存储（4%存储为4）');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `commission_setting` COMMENT '佣金设置表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('commission_setting');
    }
}
