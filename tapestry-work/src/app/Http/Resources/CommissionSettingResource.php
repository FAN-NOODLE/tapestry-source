<?php

namespace App\Http\Resources;

use Illuminate\Http\Resources\Json\Resource;

class CommissionSettingResource extends Resource
{
    /**
     * Transform the resource collection into an array.
     *
     * @param  \Illuminate\Http\Request $request
     *
     * @return array
     */
    public function toArray($request)
    {
        return [
            'id'            => $this->id,
            'enterprise_id' => $this->enterprise_id,
            'type'          => $this->type,
            'value'         => $this->value
        ];
    }
}
