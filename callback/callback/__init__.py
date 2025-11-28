from flask import Flask
from flasgger import Swagger
from configs.config import load_config
from extensions.redis import init_redis
from extensions.minio import init_minio
from utils.response import register_error_handlers
from utils.router_trace import add_request_tracing
import logging


def create_app():
    app_v1 = Flask(__name__)
    app_v1.logger.setLevel(logging.INFO)

    # init config
    load_config()

    # init redis
    init_redis()

    # init minio
    init_minio()

    # 初始化 swagger
    Swagger(app_v1)

     # 添加路由追踪
    add_request_tracing(app_v1)

    # 注册异常处理
    register_error_handlers(app_v1)

    # 注册蓝图
    from callback.routes import callback_bp
    app_v1.register_blueprint(callback_bp,url_prefix="/v1")

    return app_v1
