# utils/decorators.py
import http
import sys
from functools import wraps

from flask import g, jsonify, request

from callback.services.tavily_news import TavilyNewsAgency
from utils.response import BizError


def require_api_key(f: str) -> str:
    """
    鉴权装饰器：
    1. 校验 Header 中的 API Key。
    2. 实例化 TavilyNewsAgency 并挂载到 g.agency。
    """

    @wraps(f)
    def decorated_function(*args, **kwargs):
        api_key = None
        # 兼容 Bearer Token 格式 (Authorization: Bearer <key>)
        auth_header = request.headers.get("Authorization")
        if auth_header and auth_header.startswith("Bearer "):
            api_key = auth_header.split(" ")[1]

        # 校验 Key 是否存在
        if not api_key:
            raise BizError(
                "Authentication required:Please provide tavily api key.",
                code=http.HTTPStatus.UNAUTHORIZED,
            )

        # 实例化对象并挂载到 Flask 全局上下文 g 中
        # 这样视图函数就可以直接用 g.agency 了
        g.agency = TavilyNewsAgency(api_key=api_key)

        return f(*args, **kwargs)

    return decorated_function
