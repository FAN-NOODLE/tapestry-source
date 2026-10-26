<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangecontactStatusInEnterprise extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise')) {
            Schema::table('enterprise', function ($table) {
                $table->unsignedInteger('contact_status')->comment('企业微信通讯录应用状态（1:正常，2:禁用, 3 自建应用）')->change();
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
        Schema::table('enterprise', function (Blueprint $table) {
            //
        });
    }
}
