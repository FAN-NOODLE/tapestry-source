<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseBaseCodeLibrary extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_base_code_library')) {
            Schema::create('enterprise_base_code_library', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->tinyInteger('type')->default(1)->comment('小程序类型名称(1:商城小程序，2:企业微信小程序)');
                $table->string('version', 191)->default('')->comment('微信版本号');
                $table->string('description', 191)->default('')->comment('微信版本描述');
                $table->string('template_id', 191)->default('')->comment('微信版本对应的模板id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_base_code_library` comment '小程序模板基础库信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('enterprise_base_code_library');
    }
}
