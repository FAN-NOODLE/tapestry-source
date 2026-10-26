<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseDataRangeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_data_range')) {
            Schema::create('enterprise_data_range', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->default(0)->comment('管理员关联数据范围表主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('enterprise_admin_id')->default(0)->comment('企业账号表id')->index();
                $table->tinyInteger('type')->default(0)->comment('状态 0-全部门店 1-区域 2-门店  3-商城');
                $table->unsignedBigInteger('object_id')->default(0)->comment('门店|区域|品牌 ID');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("alter table `enterprise_data_range` comment '企业关联数据范围表'");
        }

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_data_range');
    }
}
