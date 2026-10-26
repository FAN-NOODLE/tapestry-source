<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateWelcomeSpeechTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('welcome_speech')) {
            Schema::create('welcome_speech', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('contact_qr_code_id')->default(0)->comment('联系我二维码id');
                $table->unsignedInteger('type')->default(1)->comment('1 无链接 2 图片 3 网页 4 小程序');
                $table->unsignedInteger('is_default')->default(0)->comment('是否默认 0 否 1 是');
                $table->text('object')->default('')->comment('链接内容json');
                $table->text('content')->default('')->comment('消息内容');
                $table->text('original_data')->default('')->comment('原始数据');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `welcome_speech` comment '联系我二维码表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('welcome_speech');
    }
}
