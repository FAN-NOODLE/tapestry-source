<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateJobsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('jobs')) {
            Schema::create('jobs', function (Blueprint $table) {
                $table->unsignedBigInteger('id',true)->comment("主键");
                $table->string("queue")->default('');
                $table->longText("payload")->default('');
                $table->unsignedTinyInteger("attempts");
                $table->unsignedInteger('reserved_at')->default(0);
                $table->unsignedInteger('available_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->index('queue', 'jobs_queue_index');
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
        Schema::dropIfExists('jobs');
    }
}
