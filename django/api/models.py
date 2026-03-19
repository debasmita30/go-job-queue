from django.db import models
from django.utils import timezone
import uuid

class Job(models.Model):
    STATUS_CHOICES = [
        ('pending', 'Pending'),
        ('processing', 'Processing'),
        ('completed', 'Completed'),
        ('failed', 'Failed'),
        ('dead', 'Dead'),
    ]
    
    PRIORITY_CHOICES = [
        (1, 'Low'),
        (2, 'Medium'),
        (3, 'High'),
    ]
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    type = models.CharField(max_length=100)
    payload = models.JSONField(default=dict)
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    priority = models.IntegerField(choices=PRIORITY_CHOICES, default=2)
    attempts = models.IntegerField(default=0)
    max_attempts = models.IntegerField(default=3)
    result = models.JSONField(null=True, blank=True)
    error = models.TextField(blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    started_at = models.DateTimeField(null=True, blank=True)
    completed_at = models.DateTimeField(null=True, blank=True)
    
    class Meta:
        ordering = ['-created_at']
        indexes = [
            models.Index(fields=['status', '-created_at']),
            models.Index(fields=['priority', 'status']),
        ]
    
    def __str__(self):
        return f"{self.type} - {self.status}"
