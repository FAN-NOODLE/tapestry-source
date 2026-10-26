<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateReleaseEnterpriseWork extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('release_enterprise_work')) {
            Schema::create('release_enterprise_work', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('identifier', 255)->default('release_crs_store_assistant')->comment('应用标识');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });

            DB::statement("ALTER TABLE `release_enterprise_work` COMMENT '灰度企业微信'");
        }

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
        Schema::dropIfExists('release_enterprise_work');
    }
}
