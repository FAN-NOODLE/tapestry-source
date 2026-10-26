<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateWhiteListTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('white_list'))
        Schema::create('white_list', function (Blueprint $table) {
            $table->increments('id')->comment('自增ID');
            $table->unsignedInteger('enterprise_id')->default(0)->comment('白名单企业id');
            $table->unsignedInteger('type')->default(0)->index()->comment('类型 1：先享卡');
            $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('white_list');
    }
}
