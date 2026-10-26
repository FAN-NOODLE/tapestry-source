<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateContactQrCodeParamsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('contact_qr_code_params')) {
            Schema::create('contact_qr_code_params', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('contact_qr_code_id')->default(0)->comment('联系我二维码id');
                $table->string('object')->default('')->comment('员工id/部门id');
                $table->text('object_extra')->default('')->comment('type为部门时员工id列表,逗号分隔');
                $table->string('type')->default(1)->comment('1 员工 2 部门');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `contact_qr_code_params` comment '联系我二维码关联员工/部门表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('contact_qr_code_params');
    }
}
