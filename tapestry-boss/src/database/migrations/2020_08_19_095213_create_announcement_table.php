<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateAnnouncementTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('announcement', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary();
            $table->string('title', 128)->default('')->comment('标题');
            $table->unsignedInteger('type')->default(0)->comment('1:商城公告,2:功能上线,3:注意事项');
            $table->text('content')->nullable()->comment('内容');
            $table->unsignedInteger('released_at')->default(0)->comment('发布时间');
            $table->unsignedInteger('weight')->default(0)->comment('权重');
            $table->unsignedInteger('updated_at')->default(0);
            $table->unsignedInteger('created_at')->default(0);
            $table->unsignedInteger('deleted_at')->default(0);
        });
        DB::statement("alter table `announcement` comment '公告表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('announcement');
    }
}
