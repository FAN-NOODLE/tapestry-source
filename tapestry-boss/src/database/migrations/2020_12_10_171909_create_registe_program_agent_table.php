<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateRegisteProgramAgentTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('register_program_agent', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary()->comment('自增id');
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
            $table->string('legal_persona_name',191)->default('')->comment('法人姓名');
            $table->string('legal_persona_wechat',191)->default('')->comment('法人微信');
            $table->string('name',191)->default('')->comment('企业名称');
            $table->tinyInteger('code_type')->default(0)->comment('企业类型');
            $table->string('code',191)->default('')->comment('企业代码');
            $table->string('app_id',191)->default('')->comment('第三方appid')->index();
            $table->string('auth_code',191)->default('')->comment('第三方授权码');
            $table->tinyInteger('status')->default(2)->comment('审核状态（1:审核成功，2:审核中，3：审核失败）');
            $table->text('fail_reason')->default('')->comment('审核失败原因');
            $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
        });
        DB::statement("alter table `register_program_agent` comment '小程序代注册信息表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('register_program_agent');
    }
}
