<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeSourceTypeInChannel extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('channel')) {
            Schema::table('channel', function ($table) {
                $table->unsignedInteger('source_type')->comment('建立方式，0=默认正常建立，1=内购自建,2=门店商城,3支付即导购')->change();
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
        Schema::table('channel', function (Blueprint $table) {
            //
        });
    }
}
