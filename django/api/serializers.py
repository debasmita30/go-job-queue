from rest_framework import serializers
from .models import Job

class JobSerializer(serializers.ModelSerializer):
    class Meta:
        model = Job
        fields = [
            'id',
            'type',
            'payload',
            'status',
            'priority',
            'attempts',
            'max_attempts',
            'result',
            'error',
            'created_at',
            'updated_at',
            'started_at',
            'completed_at',
        ]
        read_only_fields = [
            'id',
            'attempts',
            'result',
            'error',
            'created_at',
            'updated_at',
            'started_at',
            'completed_at',
        ]
