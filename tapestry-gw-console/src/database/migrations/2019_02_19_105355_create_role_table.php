<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateRoleTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('role')) {
            Schema::create('role', function (Blueprint $table) {
                $table->increments('id')->comment('角色id-自增');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id');
                $table->string('name', 128)->default('')->comment('角色名称');
                $table->string('describe', 128)->default('')->comment('角色描述');
                $table->string('shop_id', 128)->default('')->comment('门店id 多门店，间隔');
                $table->unsignedTinyInteger('available')->default(1)->comment('是否可用,1：可用，0不可用');
                $table->text('powers')->default('')->comment('前端传入');
                $table->unsignedTinyInteger('type')->default(1)->comment('类型,1:自然流量,2:门店,3:自然流量+门店');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `role` COMMENT '角色表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('role');
    }
}
