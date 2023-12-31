import boto3
import pytest

from moto import mock_s3
from io import BytesIO
from PIL import Image

from app.libs.aws.aws_client import AWSClient
from app.libs.aws.s3_service import S3Service
from app.libs.face_detection.image_service import ImageService
from app.libs.face_detection.positions_service import PositionsService
from app.helpers.create_test_image import create_test_image


@mock_s3
def test_generate_recognited_image(mocker, data):
    bucket_name = "test_bucket"
    image_path = "test_image_path"
    input_image = create_test_image()
    
    s3_client = AWSClient(
        "s3",
        {"aws_access_key_id": "test_key", "aws_secret_access_key": "test_secret",},
    )
    s3_service = S3Service(s3_client)
    positions_service = PositionsService()

    image_service = ImageService(positions_service, s3_service)
    image_service._bucket_name = bucket_name
    image_service._image_path = image_path
    
    s3 = s3_client.get_instance()
    s3.create_bucket(Bucket=bucket_name)
    s3.put_object(
        Bucket=bucket_name, Key=image_path, Body=input_image,
    )

    recognited_image = image_service.generate_recognited_image(data)

    assert isinstance(recognited_image, BytesIO)
    assert recognited_image.getvalue() != input_image.getvalue()


def test_generate_recognited_image_fail(mocker, data):
    s3_client = None

    s3_client = mocker.Mock(
        get_instance=lambda: mocker.Mock(get_object=lambda Bucket, Key: {"Body": {}})
    )

    s3_service = S3Service(s3_client)
    positions_service = PositionsService()

    image_service = ImageService(positions_service, s3_service)

    with pytest.raises(Exception, match="Can't generate the recognited image"):
        image_service.generate_recognited_image(data)
