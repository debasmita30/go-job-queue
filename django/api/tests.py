import pytest
from django.test import TestCase
from rest_framework.test import APIClient
from .models import Job

@pytest.mark.django_db
class TestJobAPI(TestCase):
    def setUp(self):
        self.client = APIClient()
    
    def test_create_job(self):
        data = {
            'name': 'Test Job',
            'priority': 3,
            'payload': {'data': 'test'}
        }
        response = self.client.post('/api/jobs/', data, format='json')
        assert response.status_code == 201
        assert Job.objects.count() == 1
    
    def test_job_priority_ordering(self):
        Job.objects.create(name='Low', priority=1, payload={})
        Job.objects.create(name='High', priority=3, payload={})
        
        jobs = Job.objects.all()
        assert jobs[0].priority == 3