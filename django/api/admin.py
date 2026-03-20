from django.contrib import admin
from .models import Job

@admin.register(Job)
class JobAdmin(admin.ModelAdmin):
    list_display = ['id', 'type', 'status', 'priority', 'created_at']
    list_filter = ['status', 'priority']
    search_fields = ['type', 'id']

