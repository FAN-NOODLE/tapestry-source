<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateApiPowerTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('api_power')) {
            Schema::create('api_power', function (Blueprint $table) {
                $table->increments('id')->comment('自增ID');
                $table->unsignedInteger('pid')->default(0)->comment('父ID');
                $table->string('page', 32)->default('')->comment('前端页面');
                $table->string('route', 32)->default('')->comment('路由地址');
                $table->string('describe', 128)->default('')->comment('描述');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('api_power')) {
            Schema::dropIfExists('api_power');
        }
    }
}
