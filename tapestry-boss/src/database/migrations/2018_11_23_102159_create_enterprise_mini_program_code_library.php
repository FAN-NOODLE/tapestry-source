<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseMiniProgramCodeLibrary extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_mini_program_code_library')) {
            Schema::create('enterprise_mini_program_code_library', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('account_id')->default(0)->comment('关联小程序id')->index();
                $table->string('version', 191)->default('')->comment('版本号');
                $table->string('description', 191)->default('')->comment('描述');
                $table->text('extra_info')->comment('附加信息');
                $table->string('template_id', 191)->default('')->comment('第三方模版id');
                $table->string('audit_status', 191)->default('PENDING')->comment('审核状态（PENDING，AUDITING，AUDITED，REFUSED');
                $table->string('online_status', 191)->default('PENDING')->comment('上线状态（PENDING，ONLINE，OFFLINE');
                $table->unsignedTinyInteger('using')->default(0)->comment('是否在使用该版本');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_mini_program_code_library` comment '小程序版本管理'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('enterprise_mini_program_code_library');
    }
}
