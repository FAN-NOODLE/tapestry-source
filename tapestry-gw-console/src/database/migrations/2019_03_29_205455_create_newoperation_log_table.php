<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateNewoperationLogTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('operation_log')) {
            Schema::create('operation_log', function (Blueprint $table) {
                $table->increments('id')->comment('id-自增');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->unsignedInteger('admin_id')->default(0)->index()->comment('管理员id');
                $table->unsignedTinyInteger('serve_type')->default(0)->comment('服务权限类型:0=全部,1=shop,2=crs');
                $table->string('admin_name', 128)->default('')->comment('管理员名称');
                $table->string('module', 32)->default('')->comment('操作模块');
                $table->string('page', 32)->default('')->comment('操作页面');
                $table->string('content', 255)->default('')->comment('操作内容');
                $table->string('ip', 16)->default('')->comment('IP');
                $table->string('http_method', 16)->default('GET')->comment('请求方式');
                $table->text('http_params')->default('')->comment('请求参数');
                $table->text('user_agent')->default('')->comment('User-Agent');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `operation_log` COMMENT '操作日志记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('operation_log')) {
            Schema::dropIfExists('operation_log');
        }
    }
}
