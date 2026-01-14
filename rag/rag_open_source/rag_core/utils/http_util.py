
from functools import wraps
from flask import request

def validate_request(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        # 统一进行 JSON 解析和校验
        request_json = request.get_json(silent=True)
        if request_json is None:
            raise ValueError("Invalid JSON: Please ensure Content-Type is 'application/json' and body is valid JSON.")

        # 将解析后的数据注入 kwargs
        # 为了通用性，建议使用 'json_data' 作为键名，或者根据函数签名动态注入
        kwargs['request_json'] = request_json
        return f(*args, **kwargs)

    return decorated_function
