from flask import jsonify, make_response,request,g
import time
class BizError(Exception):
    """业务异常，统一返回 code=1"""
    
    def __init__(self, msg, code=1):
        self.code = code
        self.msg = msg
        super().__init__(msg)


def response_ok(data=None, msg="success", code=0):
    return make_response(jsonify({
        "code": code,
        "msg": msg,
        "data": data or {}
    }), 200)


def response_error(code=1, msg="", http_status=400, data=None):
    return make_response(jsonify({
        "code": code,
        "msg": msg,
        "data": data or {}
    }), http_status)


def register_error_handlers(app):
    @app.errorhandler(BizError)
    def handle_biz_error(e):
        # 记录业务异常日志
        g.has_exception = True
        cost = round((time.time() - getattr(g, 'start_time', time.time())) * 1000, 2)
        app.logger.warning(
            f"[BIZ_ERROR] {request.method} {request.path} "
            f"Body={getattr(g, 'request_body', None)} "
            f"Code={e.code} Msg={e.msg} Cost={cost}ms"
        )
        return response_error(code=e.code, msg=e.msg, http_status=400)

    @app.errorhandler(Exception)
    def handle_exception(e):
        # 记录系统异常日志
        g.has_exception = True
        if request.path != '/favicon.ico':
            cost = round((time.time() - getattr(g, 'start_time', time.time())) * 1000, 2)
            app.logger.error(
                f"[EXCEPTION] {request.method} {request.path} "
                f"Body={getattr(g, 'request_body', None)} "
                f"Error={str(e)} Cost={cost}ms"
            )
        return response_error(code=500, msg=str(e), http_status=500)
