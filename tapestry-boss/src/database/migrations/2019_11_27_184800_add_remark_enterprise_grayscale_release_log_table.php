<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddRemarkEnterpriseGrayscaleReleaseLogTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_grayscale_release_log')) {
            Schema::table('enterprise_grayscale_release_log', function (Blueprint $table) {
                $table->string('remark', 90)->default('')->after('can_cancel')->comment('备注');
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
        Schema::table('enterprise_grayscale_release_log', function (Blueprint $table) {
            //
        });
    }
}
