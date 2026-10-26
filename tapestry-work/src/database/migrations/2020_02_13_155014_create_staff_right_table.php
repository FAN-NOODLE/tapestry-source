<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateStaffRightTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('staff_right')) {
            Schema::create('staff_right', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('obj_id', 32)->default('')->comment('员工/部门ID');
                $table->unsignedTinyInteger('obj_type')->default(0)->comment('obj_id类型（1=部门、2=员工）');
                $table->string('label')->default('')->comment('员工名称');
                $table->string('value')->default('')->comment('type为1时:dp_third_id,type为2时:user_id (UC)');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('enterprise_id', 'enterprise_id');
                $table->index('obj_id', 'obj_id');
            });

            // 表注释
            DB::statement("ALTER TABLE `staff_right` comment '内购员工权限表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('staff_right');
    }
}
