<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateChannelTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('channel')) {
            Schema::create('channel', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('channel', 100)->default('')->comment('渠道名称')->index();
                $table->unsignedInteger('admin_id')->default(0)->comment('管理员id');
                $table->string('admin_name')->default('')->comment('管理员名称');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `channel` comment '渠道表'");
            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `channel` AUTO_INCREMENT=4");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('channel');
    }
}
