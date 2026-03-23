from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.decorators import action
from .models import Job
from .serializers import JobSerializer
from django.db import connection
from django.db.utils import ProgrammingError

class JobViewSet(viewsets.ModelViewSet):
    queryset = Job.objects.all()
    serializer_class = JobSerializer
    
    def get_queryset(self):
        queryset = Job.objects.all()
        job_status = self.request.query_params.get('status')
        if job_status:
            queryset = queryset.filter(status=job_status)
        return queryset.order_by('-created_at')
    
    def create(self, request, *args, **kwargs):
        try:
            serializer = self.get_serializer(data=request.data)
            serializer.is_valid(raise_exception=True)
            self.perform_create(serializer)
            return Response(
                {
                    'job': serializer.data,
                    'message': 'Job added to queue'
                },
                status=status.HTTP_201_CREATED
            )
        except ProgrammingError:
            return Response(
                {'error': 'Database table not initialized. Please wait for migrations to complete.'},
                status=status.HTTP_503_SERVICE_UNAVAILABLE
            )

class StatsViewSet(viewsets.ViewSet):
    
    def list(self, request):
        try:
            # Check if table exists first
            with connection.cursor() as cursor:
                cursor.execute("SELECT 1 FROM api_job LIMIT 1")
            
            # If we get here, table exists
            return Response({
                'total_jobs': Job.objects.count(),
                'pending_jobs': Job.objects.filter(status='pending').count(),
                'processing_jobs': Job.objects.filter(status='processing').count(),
                'completed_jobs': Job.objects.filter(status='completed').count(),
                'failed_jobs': Job.objects.filter(status='failed').count(),
                'dead_jobs': Job.objects.filter(status='dead').count(),
            })
        except (ProgrammingError, Exception):
            # Table doesn't exist yet - return empty stats
            return Response({
                'total_jobs': 0,
                'pending_jobs': 0,
                'processing_jobs': 0,
                'completed_jobs': 0,
                'failed_jobs': 0,
                'dead_jobs': 0,
            })
