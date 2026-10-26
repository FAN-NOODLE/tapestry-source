<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateStaffSettingTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('staff_setting')) {
            Schema::create('staff_setting', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id');
                $table->unsignedTinyInteger('binding_type')->default(0)->comment('绑定时效：1:永久有效 2：有效期');
                $table->unsignedInteger('duration_time')->default(0)->comment('有效持续时长');
                $table->unsignedTinyInteger('change_owner')->default(0)->comment('是否允许抢客 1:不允许 2：允许');
                $table->unsignedInteger('protected_time')->default(0)->comment('保护时长');
                $table->unsignedTinyInteger('commission_type')->default(0)->comment('内购员工购买商品是否有佣金1:有 2：无');
                $table->unsignedTinyInteger('settle_type')->default(0)->comment('佣金结算方式：1:交易完成结算 2：售后维权期结束结算');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('enterprise_id', 'enterprise_id');
            });

            // 表注释
            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `staff_setting` comment '内购员工设置表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('staff_setting');
    }
}
