from flask import Blueprint,request 
from callback.services.hello_demo import get_message
from utils.response import response_ok
from utils.response import BizError
hello_demo = Blueprint("hello", __name__,url_prefix="/")

@hello_demo.route("/")
def index():
    """
    首页接口  
    ---
    tags:
      - Hello
    responses:
      200:
        description: 返回 Hello 信息
        schema:
          type: object
          properties:
            msg:
              type: string
              example: Hello from service
    """
    message = get_message()
    return response_ok(f"Hello from service:{message}")


@hello_demo.route("/api/hello")
def hello():
    """
    Hello API 示例  
    ---
    tags:
      - Hello
    responses:
      200:
        description: 返回静态 hello
        schema:
          type: object
          properties:
            msg:
              type: string
              example: Hello from API
    """
    username = request.args.get("user")
    if not username:
      raise BizError("username is None", code=2002)
    
    return response_ok("Hello from API")
