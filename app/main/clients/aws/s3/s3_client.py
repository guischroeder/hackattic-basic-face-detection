import boto3
from app.main.clients.aws.aws_client import AWSClient


class S3Client(AWSClient):
    _client = None

    def __init__(self, credentials, region_name):
        self._credentials = credentials
        self._region_name = region_name

    def get_instance(self):
        if not self._client:
            self._client = boto3.client(
                "s3", region_name=self._region_name, **self._credentials
            )

        return self._client