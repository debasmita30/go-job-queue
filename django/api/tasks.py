from celery import shared_task
from .models import Job
import json

@shared_task(bind=True, max_retries=3)
def process_job(self, job_id):
    job = Job.objects.get(id=job_id)
    try:
        job.status = 'processing'
        job.save()
        
        # Simulate job processing
        result = {
            "processed": True,
            "data": job.payload,
            "message": f"Job {job.id} processed successfully"
        }
        
        job.result = result
        job.status = 'completed'
        job.save()
        
    except Exception as exc:
        job.retries += 1
        if job.retries < job.max_retries:
            job.status = 'pending'
            job.save()
            # Retry after 60 seconds
            raise self.retry(exc=exc, countdown=60)
        else:
            job.status = 'failed'
            job.result = {"error": str(exc)}
            job.save()