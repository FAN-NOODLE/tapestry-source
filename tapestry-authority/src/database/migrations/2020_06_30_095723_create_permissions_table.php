<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreatePermissionsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('permissions')) {
            Schema::create('permissions', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary();
                $table->unsignedBigInteger('pid')->default(0)->comment('父ID');
                $table->string('name', 64)->default('')->comment('权限名称');
                $table->string('path', 128)->default('')->comment('前端path');
                $table->string('key', 128)->default('')->comment('前端key');
                $table->string('icon', 32)->default('')->comment('前端图标');
                $table->unsignedTinyInteger('type')->default(0)->comment('类型 1:菜单,2:视图,3:按钮,4:菜单+视图');
                $table->unsignedInteger('weight')->default(0)->comment('权重');
                $table->unsignedInteger('system_type')->default(0)->comment('系统类型,1:scrm,2:isg,3:store,4:cdp,5:shop');
                $table->string('version_tag', 16)->default('')->comment('版本号 迭代号');
                $table->unsignedTinyInteger('is_full_release')->default(0)->comment('是否全量发布 0:否,1:是');
                $table->unsignedTinyInteger('center_shop_store')->default(0)->comment('中心商城(有门店)权限 1:标配,2:选配,3:无此功能');
                $table->unsignedTinyInteger('center_shop')->default(0)->comment('中心商城权限 1:标配,2:选配,3:无此功能');
                $table->unsignedTinyInteger('direct_chain_store')->default(0)->comment('直营连锁门店权限 1:标配,2:选配,3:无此功能');

                $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');

                $table->unique('key', 'key_unique-index');
            });

            \Illuminate\Support\Facades\DB::statement("alter table `permissions` comment '权限表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('permissions');
    }
}
