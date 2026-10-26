<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddCanCancelToEnterpriseGrayscaleReleaseLogTable extends Migration
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
                $table->tinyInteger('can_cancel')->default(0)->after('status')->comment('是否可以撤销(1：可以，0：不可以)');
                $table->unsignedInteger('cancel_at')->default(0)->after('status')->comment('撤销时间');
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
