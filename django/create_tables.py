import os
import django
from django.conf import settings

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
django.setup()

from django.db import connection
from django.db.migrations.executor import MigrationExecutor

def create_tables():
    """Run migrations to create tables"""
    executor = MigrationExecutor(connection)
    plan = executor.migration_plan(executor.loader.graph.leaf_nodes())
    
    if plan:
        executor.migrate(executor.loader.graph.leaf_nodes())
        print("✅ Tables created successfully!")
    else:
        print("✅ Tables already exist!")

if __name__ == '__main__':
    create_tables()
