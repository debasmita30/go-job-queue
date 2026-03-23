from rest_framework import viewsets, status
from rest_framework.response import Response
from .serializers import JobSerializer

class JobViewSet(viewsets.ViewSet):
    
    def list(self, request):
        return Response({'results': []})
    
    def create(self, request, *args, **kwargs):
        return Response(
            {'job': request.data, 'message': 'Job received'},
            status=status.HTTP_201_CREATED
        )

class StatsViewSet(viewsets.ViewSet):
    
    def list(self, request):
        return Response({
            'total_jobs': 0,
            'pending_jobs': 0,
            'processing_jobs': 0,
            'completed_jobs': 0,
            'failed_jobs': 0,
            'dead_jobs': 0,
        })
