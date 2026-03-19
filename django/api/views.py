from rest_framework.viewsets import ViewSet
from rest_framework.response import Response
from rest_framework.decorators import action
from .models import Job

class StatsViewSet(ViewSet):
    """
    API endpoint for queue statistics
    """
    
    @action(detail=False, methods=['get'])
    def list(self, request):
        """
        Get overall queue statistics
        GET /api/stats/
        """
        total_jobs = Job.objects.count()
        pending_jobs = Job.objects.filter(status='pending').count()
        processing_jobs = Job.objects.filter(status='processing').count()
        completed_jobs = Job.objects.filter(status='completed').count()
        failed_jobs = Job.objects.filter(status='failed').count()
        dead_jobs = Job.objects.filter(status='dead').count()
        
        return Response({
            'total_jobs': total_jobs,
            'pending_jobs': pending_jobs,
            'processing_jobs': processing_jobs,
            'completed_jobs': completed_jobs,
            'failed_jobs': failed_jobs,
            'dead_jobs': dead_jobs,
        })
        
