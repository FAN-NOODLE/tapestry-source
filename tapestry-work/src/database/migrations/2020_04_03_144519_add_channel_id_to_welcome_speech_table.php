<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddChannelIdToWelcomeSpeechTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('welcome_speech')) {
            Schema::table('welcome_speech', function (Blueprint $table) {
                $table->unsignedInteger('channel_id')->default('0')->after('enterprise_id')->comment('渠道id');
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
        if (Schema::hasTable('welcome_speech')) {
            Schema::table('welcome_speech', function (Blueprint $table) {
                $table->dropColumn('channel_id');
            });
        }
    }
}
