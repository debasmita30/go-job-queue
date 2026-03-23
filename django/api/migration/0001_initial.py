# Generated migration file

from django.db import migrations, models
import uuid


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Job',
            fields=[
                ('id', models.UUIDField(default=uuid.uuid4, primary_key=True, serialize=False)),
                ('type', models.CharField(max_length=100)),
                ('payload', models.JSONField(default=dict)),
                ('status', models.CharField(choices=[('pending', 'Pending'), ('processing', 'Processing'), ('completed', 'Completed'), ('failed', 'Failed'), ('dead', 'Dead')], default='pending', max_length=20)),
                ('priority', models.IntegerField(default=2)),
                ('attempts', models.IntegerField(default=0)),
                ('max_attempts', models.IntegerField(default=3)),
                ('result', models.JSONField(blank=True, null=True)),
                ('error', models.TextField(blank=True)),
                ('created_at', models.DateTimeField(auto_now_add=True)),
                ('updated_at', models.DateTimeField(auto_now=True)),
            ],
            options={
                'ordering': ['-created_at'],
            },
        ),
    ]
