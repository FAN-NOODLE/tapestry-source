<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateLogTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('log')) {
            Schema::create('log', function (Blueprint $table) {
                $table->increments('id')->comment('自增ID');
                $table->string('request_id', 50)->default('')->comment('');
                $table->integer('code')->default(0)->comment('');
                $table->string('message', 500)->default('')->comment('');
                $table->text('params')->default('')->comment('');
                $table->text('data')->default('')->comment('');
                $table->string('url', 500)->default('')->comment('');
                $table->string('ip', 50)->default('')->comment('');
                $table->integer('port')->default(0)->comment('');
                $table->string('method', 50)->default('')->comment('');
                $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
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
    if (Schema::hasTable('log')) {
        Schema::dropIfExists('log');
    }
}
}
