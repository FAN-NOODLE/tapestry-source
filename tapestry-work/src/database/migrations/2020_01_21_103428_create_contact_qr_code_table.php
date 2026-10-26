<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateContactQrCodeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('contact_qr_code')) {
            Schema::create('contact_qr_code', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('channel')->default('')->comment('渠道名称');
                $table->unsignedInteger('type')->default(1)->comment('1 单人 2 多人');
                $table->unsignedInteger('is_verify')->default(1)->comment('0 需要验证 1 无需验证');
                $table->string('official_url')->default('')->comment('官方url');
                $table->string('transform_url')->default('')->comment('转换存储url');
                $table->string('uc_id')->default('')->comment('uc返回二维码id');
                $table->string('config_id')->default('')->comment('二维码id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `contact_qr_code` comment '联系我二维码表'");
            DB::statement("ALTER TABLE `contact_qr_code` AUTO_INCREMENT=10");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('contact_qr_code');
    }
}
