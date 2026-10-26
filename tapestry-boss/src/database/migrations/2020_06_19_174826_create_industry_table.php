<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateIndustryTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('industry', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary()->comment('自增id');
            $table->unsignedBigInteger('pid')->default(0)->comment('父行业');
            $table->string('name')->default('')->comment('行业名称');
            $table->unsignedInteger('updated_at')->default(0);
            $table->unsignedInteger('created_at')->default(0);
            $table->unsignedInteger('deleted_at')->default(0);
        });

        DB::statement("alter table `industry` comment '行业表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('industry');
    }
}
