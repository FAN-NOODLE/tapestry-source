<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldToEnterpriseActivityMessage extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_activity_message')) {
            Schema::table('enterprise_activity_message', function (Blueprint $table) {
                $table->unsignedInteger('real_num')->default(0)->comment('真实发送数量');
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
        Schema::table('enterprise_activity_message', function (Blueprint $table) {
            $table->dropColumn('real_num');
        });
    }
}
