<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class WorkSunCode extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('work_sun_code')) {

            Schema::create('work_sun_code', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('导购ID');
                $table->string('sun_code')->default('')->comment('太阳码');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("ALTER TABLE `work_sun_code` COMMENT '导购太阳码'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('work_sun_code');
    }
}
