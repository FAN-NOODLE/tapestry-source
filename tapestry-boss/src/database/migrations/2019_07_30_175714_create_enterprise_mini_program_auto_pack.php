<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseMiniProgramAutoPack extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_mini_program_auto_pack')) {
            Schema::create('enterprise_mini_program_auto_pack', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('type')->default(1)->comment('小程序类型(1:商城小程序，2:企业微信小程序)');
                $table->string('task_id', 64)->default('')->comment('build平台打包任务id')->index();
                $table->string('version', 64)->default('')->comment('版本号');
                $table->unsignedBigInteger('version_id')->default(0)->comment('版本id')->index();
                $table->string('wx_version', 64)->default('')->comment('提交到开放平台版本号');
                $table->string('wx_desc', 64)->default('')->comment('提交到开放平台版本描述');
                $table->unsignedInteger('status')->default(0)->comment('0=排队中,1=进行中,2=成功,3=失败,4=需要登录');
                $table->text('desc')->comment('状态描述');
                $table->text('task_config')->comment('build平台打包任务请求参数');
                $table->text('login_url')->comment('需要登录开发者工具二维码base64码');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_mini_program_auto_pack` comment '企业小程序自动打包'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('enterprise_mini_program_auto_pack')) {
            Schema::dropIfExists('enterprise_mini_program_auto_pack');
        }
    }
}
