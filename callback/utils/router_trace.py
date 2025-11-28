import time
from flask import Flask, request, g 

def add_request_tracing(app: Flask):
    """封装路由追踪"""
    @app.before_request
    def start_trace():
        g.start_time = time.time()
        g.request_body = request.get_json(silent=True)

    @app.after_request
    def end_trace(response):
        if request.path == '/favicon.ico':
            return response
        # 如果已经在异常处理中打印过日志，就跳过 TRACE 日志
        if getattr(g, 'has_exception', False):
            return response
        cost = round((time.time() - g.start_time) * 1000, 2)
        app.logger.info(
            f"[TRACE] {request.method} {request.path} "
            f"Body={g.request_body} Status={response.status_code} Cost={cost}ms"
        )
        return response