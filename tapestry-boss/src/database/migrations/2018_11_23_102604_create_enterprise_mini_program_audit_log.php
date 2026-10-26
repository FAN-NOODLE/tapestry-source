<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseMiniProgramAuditLog extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_mini_program_audit_log')) {
            Schema::create('enterprise_mini_program_audit_log', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('code_library_id')->default(0)->comment('小程序代码库id');
                $table->string('audit_id', 191)->default('')->comment('微信审核id');
                $table->text('commit')->comment('提交审核信息');
                $table->string('status', 191)->default('PENDING')->comment('提交审核状态（PENDING，AUDITING，AUDITED，REFUSED');
                $table->string('reason', 191)->default('')->comment('原因');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_mini_program_audit_log` comment '小程序审核日志表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('enterprise_mini_program_audit_log');
    }
}
