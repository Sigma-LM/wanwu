from minio import Minio
from configs.config import config

minio_client: Minio = None

def init_minio():
    global minio_client

    minio_client = Minio(
        endpoint=config.callback_cfg["MINIO"]["ENDPOINT"],
        access_key=config.callback_cfg["MINIO"]["USER"],
        secret_key=config.callback_cfg["MINIO"]["PASSWORD"],
        secure=config.callback_cfg["MINIO"]["SECURE"]
    )

    print("MinIO client initialized.")
    return minio_client