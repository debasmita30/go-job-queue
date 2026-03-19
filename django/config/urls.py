from django.contrib import admin
from django.urls import path, include
from rest_framework.routers import DefaultRouter
from api.views import JobViewSet, StatsViewSet

router = DefaultRouter()
router.register(r'jobs', JobViewSet, basename='job')
router.register(r'stats', StatsViewSet, basename='stats')  # ← ADD THIS

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/', include(router.urls)),
]
