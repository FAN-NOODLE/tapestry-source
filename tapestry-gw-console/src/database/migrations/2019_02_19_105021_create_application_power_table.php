<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateApplicationPowerTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {

        if (!Schema::hasTable('application_power')) {
            Schema::create('application_power', function (Blueprint $table) {
                $table->increments('id')->comment('应用权限id-自增');
                $table->unsignedTinyInteger('serve_type')->default(0)->comment('服务权限类型:0=全部,1=shop,2=crs');
                $table->string('name', 128)->default('')->comment('权限名称');
                $table->string('alias', 128)->default('')->comment('路由别名');
                $table->string('component', 128)->default('')->comment('前端component');
                $table->string('path', 128)->default('')->comment('权限标示');
                $table->string('parent', 128)->default('')->comment('父级权限标识');
                $table->unsignedTinyInteger('available')->default(1)->comment('是否可用,1：可用，0不可用');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `application_power` COMMENT '应用权限表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('application_power');
    }
}
